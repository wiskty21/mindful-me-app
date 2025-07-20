package models

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User ユーザーモデル
type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email          string             `bson:"email" json:"email"`
	Name           string             `bson:"name" json:"name"`
	PasswordHash   string             `bson:"password_hash" json:"-"`
	TotalSessions  int                `bson:"total_sessions" json:"total_sessions"`
	TotalMinutes   int                `bson:"total_minutes" json:"total_minutes"`
	CurrentStreak  int                `bson:"current_streak" json:"current_streak"`
	LongestStreak  int                `bson:"longest_streak" json:"longest_streak"`
	LastSessionAt  *time.Time         `bson:"last_session_at" json:"last_session_at"`
	PreferredTime  int                `bson:"preferred_time" json:"preferred_time"` // デフォルト瞑想時間（分）
	ReminderTime   string             `bson:"reminder_time" json:"reminder_time"`   // リマインダー時刻 "HH:MM"
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

// UserService ユーザーサービス
type UserService struct {
	collection *mongo.Collection
}

// NewUserService ユーザーサービスを作成
func NewUserService(db *mongo.Database) *UserService {
	return &UserService{
		collection: db.Collection("users"),
	}
}

// CreateUser ユーザーを作成
func (s *UserService) CreateUser(ctx context.Context, user *User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.TotalSessions = 0
	user.TotalMinutes = 0
	user.CurrentStreak = 0
	user.LongestStreak = 0
	user.PreferredTime = 10 // デフォルト10分

	result, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("email already exists")
		}
		return err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetUserByEmail メールアドレスでユーザーを取得
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByID IDでユーザーを取得
func (s *UserService) GetUserByID(ctx context.Context, id primitive.ObjectID) (*User, error) {
	var user User
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUserStats ユーザーの統計情報を更新
func (s *UserService) UpdateUserStats(ctx context.Context, userID primitive.ObjectID, duration int) error {
	now := time.Now()
	
	// 現在のユーザー情報を取得
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// ストリークの計算
	newStreak := user.CurrentStreak
	if user.LastSessionAt != nil {
		// 前回のセッションから24時間以内なら継続
		if now.Sub(*user.LastSessionAt) <= 24*time.Hour {
			newStreak++
		} else if now.Sub(*user.LastSessionAt) > 48*time.Hour {
			// 48時間以上空いたらリセット
			newStreak = 1
		}
	} else {
		newStreak = 1
	}

	// 最長ストリークの更新
	longestStreak := user.LongestStreak
	if newStreak > longestStreak {
		longestStreak = newStreak
	}

	update := bson.M{
		"$inc": bson.M{
			"total_sessions": 1,
			"total_minutes":  duration,
		},
		"$set": bson.M{
			"current_streak":  newStreak,
			"longest_streak":  longestStreak,
			"last_session_at": now,
			"updated_at":      now,
		},
	}

	_, err = s.collection.UpdateByID(ctx, userID, update)
	return err
}

// UpdateProfile プロフィール情報を更新
func (s *UserService) UpdateProfile(ctx context.Context, userID primitive.ObjectID, name string, preferredTime int, reminderTime string) error {
	update := bson.M{
		"$set": bson.M{
			"name":           name,
			"preferred_time": preferredTime,
			"reminder_time":  reminderTime,
			"updated_at":     time.Now(),
		},
	}

	_, err := s.collection.UpdateByID(ctx, userID, update)
	return err
}

// CreateIndexes インデックスを作成
func (s *UserService) CreateIndexes(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := s.collection.Indexes().CreateOne(ctx, indexModel)
	return err
}