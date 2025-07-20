package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mindfulme-backend/models"
)

// CreateSessionRequest セッション作成リクエスト
type CreateSessionRequest struct {
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
	Type      string    `json:"type" binding:"required,oneof=timer guided breathing"`
	Mood      string    `json:"mood,omitempty"`
	Note      string    `json:"note,omitempty"`
}

// UpdateSessionRequest セッション更新リクエスト
type UpdateSessionRequest struct {
	Mood string `json:"mood" binding:"required,oneof=calm focused relaxed energized"`
	Note string `json:"note"`
}

// GetSessions セッション履歴取得ハンドラー
func GetSessions(sessionService *models.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetString("userID")
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// クエリパラメータからlimitを取得（デフォルト30）
		limit := 30
		if limitStr := c.Query("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
				limit = l
			}
		}

		ctx := context.Background()
		sessions, err := sessionService.GetSessions(ctx, userID, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sessions"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"sessions": sessions,
			"count":    len(sessions),
		})
	}
}

// CreateSession セッション作成ハンドラー
func CreateSession(sessionService *models.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetString("userID")
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var req CreateSessionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// セッション時間を計算
		duration := int(req.EndTime.Sub(req.StartTime).Minutes())
		if duration <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session duration"})
			return
		}

		session := &models.Session{
			UserID:    userID,
			StartTime: req.StartTime,
			EndTime:   req.EndTime,
			Duration:  duration,
			Type:      req.Type,
			Mood:      req.Mood,
			Note:      req.Note,
		}

		ctx := context.Background()
		if err := sessionService.CreateSession(ctx, session); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
			return
		}

		// ユーザーの統計情報を更新
		userService := c.MustGet("userService").(*models.UserService)
		if err := userService.UpdateUserStats(ctx, userID, duration); err != nil {
			// エラーログだけ記録して続行
			c.JSON(http.StatusCreated, gin.H{
				"session": session,
				"warning": "Failed to update user stats",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"session": session})
	}
}

// UpdateSession セッション更新ハンドラー
func UpdateSession(sessionService *models.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetString("userID")
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		sessionIDStr := c.Param("id")
		sessionID, err := primitive.ObjectIDFromHex(sessionIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
			return
		}

		var req UpdateSessionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := context.Background()
		if err := sessionService.UpdateSession(ctx, sessionID, userID, req.Mood, req.Note); err != nil {
			if err == primitive.ErrInvalidHex {
				c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Session updated successfully"})
	}
}

// GetStats 統計情報取得ハンドラー
func GetStats(sessionService *models.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetString("userID")
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		ctx := context.Background()
		stats, err := sessionService.GetStats(ctx, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
			return
		}

		c.JSON(http.StatusOK, stats)
	}
}

// GetStreak 連続日数取得ハンドラー
func GetStreak(sessionService *models.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetString("userID")
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		ctx := context.Background()
		current, longest, err := sessionService.GetStreak(ctx, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get streak"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"current_streak": current,
			"longest_streak": longest,
		})
	}
}