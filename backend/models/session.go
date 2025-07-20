package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Session 瞑想セッションモデル
type Session struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	StartTime  time.Time          `bson:"start_time" json:"start_time"`
	EndTime    time.Time          `bson:"end_time" json:"end_time"`
	Duration   int                `bson:"duration" json:"duration"` // 分単位
	Type       string             `bson:"type" json:"type"`         // "timer", "guided", "breathing"
	Mood       string             `bson:"mood,omitempty" json:"mood,omitempty"` // "calm", "focused", "relaxed", "energized"
	Note       string             `bson:"note,omitempty" json:"note,omitempty"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}

// SessionStats セッション統計
type SessionStats struct {
	TotalSessions   int     `json:"total_sessions"`
	TotalMinutes    int     `json:"total_minutes"`
	AverageMinutes  float64 `json:"average_minutes"`
	CurrentStreak   int     `json:"current_streak"`
	LongestStreak   int     `json:"longest_streak"`
	ThisWeek        int     `json:"this_week"`
	ThisMonth       int     `json:"this_month"`
	MostCommonMood  string  `json:"most_common_mood"`
}

// SessionService セッションサービス
type SessionService struct {
	collection *mongo.Collection
}

// NewSessionService セッションサービスを作成
func NewSessionService(db *mongo.Database) *SessionService {
	return &SessionService{
		collection: db.Collection("sessions"),
	}
}

// CreateSession セッションを作成
func (s *SessionService) CreateSession(ctx context.Context, session *Session) error {
	session.CreatedAt = time.Now()
	
	// 時間を計算
	if session.Duration == 0 && !session.EndTime.IsZero() {
		session.Duration = int(session.EndTime.Sub(session.StartTime).Minutes())
	}

	result, err := s.collection.InsertOne(ctx, session)
	if err != nil {
		return err
	}

	session.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetSessions ユーザーのセッション履歴を取得
func (s *SessionService) GetSessions(ctx context.Context, userID primitive.ObjectID, limit int) ([]Session, error) {
	var sessions []Session

	opts := options.Find().
		SetSort(bson.D{{Key: "start_time", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := s.collection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

// UpdateSession セッションを更新（気分やメモの追加）
func (s *SessionService) UpdateSession(ctx context.Context, sessionID, userID primitive.ObjectID, mood, note string) error {
	filter := bson.M{
		"_id":     sessionID,
		"user_id": userID,
	}

	update := bson.M{
		"$set": bson.M{
			"mood": mood,
			"note": note,
		},
	}

	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// GetStats ユーザーの統計情報を取得
func (s *SessionService) GetStats(ctx context.Context, userID primitive.ObjectID) (*SessionStats, error) {
	stats := &SessionStats{}

	// 全セッション数と合計時間
	pipeline := []bson.M{
		{"$match": bson.M{"user_id": userID}},
		{"$group": bson.M{
			"_id": nil,
			"total_sessions": bson.M{"$sum": 1},
			"total_minutes": bson.M{"$sum": "$duration"},
		}},
	}

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	if len(result) > 0 {
		if total, ok := result[0]["total_sessions"].(int32); ok {
			stats.TotalSessions = int(total)
		}
		if minutes, ok := result[0]["total_minutes"].(int32); ok {
			stats.TotalMinutes = int(minutes)
		}
		if stats.TotalSessions > 0 {
			stats.AverageMinutes = float64(stats.TotalMinutes) / float64(stats.TotalSessions)
		}
	}

	// 今週のセッション数
	weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
	weekCount, err := s.collection.CountDocuments(ctx, bson.M{
		"user_id": userID,
		"start_time": bson.M{"$gte": weekStart},
	})
	if err == nil {
		stats.ThisWeek = int(weekCount)
	}

	// 今月のセッション数
	monthStart := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)
	monthCount, err := s.collection.CountDocuments(ctx, bson.M{
		"user_id": userID,
		"start_time": bson.M{"$gte": monthStart},
	})
	if err == nil {
		stats.ThisMonth = int(monthCount)
	}

	// 最も多い気分
	moodPipeline := []bson.M{
		{"$match": bson.M{"user_id": userID, "mood": bson.M{"$ne": ""}}},
		{"$group": bson.M{
			"_id": "$mood",
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"count": -1}},
		{"$limit": 1},
	}

	moodCursor, err := s.collection.Aggregate(ctx, moodPipeline)
	if err == nil {
		defer moodCursor.Close(ctx)
		var moodResult []bson.M
		if err := moodCursor.All(ctx, &moodResult); err == nil && len(moodResult) > 0 {
			if mood, ok := moodResult[0]["_id"].(string); ok {
				stats.MostCommonMood = mood
			}
		}
	}

	return stats, nil
}

// GetStreak 連続日数を計算
func (s *SessionService) GetStreak(ctx context.Context, userID primitive.ObjectID) (current int, longest int, err error) {
	// 日付順にセッションを取得
	opts := options.Find().SetSort(bson.D{{Key: "start_time", Value: -1}})
	cursor, err := s.collection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return 0, 0, err
	}
	defer cursor.Close(ctx)

	var sessions []Session
	if err := cursor.All(ctx, &sessions); err != nil {
		return 0, 0, err
	}

	if len(sessions) == 0 {
		return 0, 0, nil
	}

	// 連続日数を計算
	current = 1
	longest = 1
	tempStreak := 1
	
	today := time.Now().Truncate(24 * time.Hour)
	lastDate := sessions[0].StartTime.Truncate(24 * time.Hour)
	
	// 今日セッションがない場合はカレントを0に
	if !lastDate.Equal(today) && !lastDate.Equal(today.AddDate(0, 0, -1)) {
		current = 0
	}

	for i := 1; i < len(sessions); i++ {
		currentDate := sessions[i].StartTime.Truncate(24 * time.Hour)
		dayDiff := lastDate.Sub(currentDate).Hours() / 24

		if dayDiff == 1 {
			tempStreak++
			if tempStreak > longest {
				longest = tempStreak
			}
		} else if dayDiff > 1 {
			tempStreak = 1
		}

		lastDate = currentDate
	}

	if current > 0 {
		current = tempStreak
	}

	return current, longest, nil
}

// CreateIndexes インデックスを作成
func (s *SessionService) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "start_time", Value: -1},
			},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	}

	_, err := s.collection.Indexes().CreateMany(ctx, indexes)
	return err
}