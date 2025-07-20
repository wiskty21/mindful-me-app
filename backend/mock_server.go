package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type User struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	TotalSessions  int    `json:"total_sessions"`
	TotalMinutes   int    `json:"total_minutes"`
	CurrentStreak  int    `json:"current_streak"`
	LongestStreak  int    `json:"longest_streak"`
	LastSessionAt  string `json:"last_session_at,omitempty"`
	PreferredTime  int    `json:"preferred_time"`
	ReminderTime   string `json:"reminder_time"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type AuthResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func main() {
	r := gin.Default()

	// CORS設定
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3003", "http://192.168.0.7:3003"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ヘルスチェック
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now()})
	})

	// API v1 グループ
	v1 := r.Group("/api/v1")

	// 認証 API
	auth := v1.Group("/auth")
	{
		auth.POST("/register", func(c *gin.Context) {
			mockUser := User{
				ID:            "mock-user-id-123",
				Email:         "test@mindfulme.app",
				Name:          "テストユーザー",
				TotalSessions: 0,
				TotalMinutes:  0,
				CurrentStreak: 0,
				LongestStreak: 0,
				PreferredTime: 10,
				ReminderTime:  "08:00",
				CreatedAt:     time.Now().Format(time.RFC3339),
				UpdatedAt:     time.Now().Format(time.RFC3339),
			}

			response := AuthResponse{
				User:         mockUser,
				AccessToken:  "mock-access-token-" + time.Now().Format("20060102150405"),
				RefreshToken: "mock-refresh-token-" + time.Now().Format("20060102150405"),
			}

			c.JSON(http.StatusOK, response)
		})

		auth.POST("/login", func(c *gin.Context) {
			mockUser := User{
				ID:            "mock-user-id-123",
				Email:         "test@mindfulme.app",
				Name:          "テストユーザー",
				TotalSessions: 5,
				TotalMinutes:  47,
				CurrentStreak: 3,
				LongestStreak: 7,
				LastSessionAt: time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
				PreferredTime: 10,
				ReminderTime:  "08:00",
				CreatedAt:     time.Now().Add(-30 * 24 * time.Hour).Format(time.RFC3339),
				UpdatedAt:     time.Now().Format(time.RFC3339),
			}

			response := AuthResponse{
				User:         mockUser,
				AccessToken:  "mock-access-token-" + time.Now().Format("20060102150405"),
				RefreshToken: "mock-refresh-token-" + time.Now().Format("20060102150405"),
			}

			c.JSON(http.StatusOK, response)
		})

		auth.POST("/refresh", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"access_token":  "mock-new-access-token-" + time.Now().Format("20060102150405"),
				"refresh_token": "mock-new-refresh-token-" + time.Now().Format("20060102150405"),
			})
		})

		auth.POST("/logout", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ログアウトしました"})
		})
	}

	// プロフィール API
	v1.GET("/profile", func(c *gin.Context) {
		mockUser := User{
			ID:            "mock-user-id-123",
			Email:         "test@mindfulme.app",
			Name:          "テストユーザー",
			TotalSessions: 5,
			TotalMinutes:  47,
			CurrentStreak: 3,
			LongestStreak: 7,
			LastSessionAt: time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
			PreferredTime: 10,
			ReminderTime:  "08:00",
			CreatedAt:     time.Now().Add(-30 * 24 * time.Hour).Format(time.RFC3339),
			UpdatedAt:     time.Now().Format(time.RFC3339),
		}

		c.JSON(http.StatusOK, mockUser)
	})

	v1.PUT("/profile", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "プロフィールを更新しました"})
	})

	// セッション API
	sessions := v1.Group("/sessions")
	{
		sessions.GET("/", func(c *gin.Context) {
			mockSessions := []gin.H{
				{
					"id":         "session-1",
					"user_id":    "mock-user-id-123",
					"start_time": time.Now().Add(-2 * 24 * time.Hour).Format(time.RFC3339),
					"end_time":   time.Now().Add(-2*24*time.Hour + 10*time.Minute).Format(time.RFC3339),
					"duration":   10,
					"type":       "timer",
					"mood":       "calm",
					"note":       "朝の瞑想",
					"created_at": time.Now().Add(-2 * 24 * time.Hour).Format(time.RFC3339),
				},
				{
					"id":         "session-2",
					"user_id":    "mock-user-id-123",
					"start_time": time.Now().Add(-1 * 24 * time.Hour).Format(time.RFC3339),
					"end_time":   time.Now().Add(-1*24*time.Hour + 15*time.Minute).Format(time.RFC3339),
					"duration":   15,
					"type":       "timer",
					"mood":       "focused",
					"note":       "昼休みの瞑想",
					"created_at": time.Now().Add(-1 * 24 * time.Hour).Format(time.RFC3339),
				},
			}

			c.JSON(http.StatusOK, mockSessions)
		})

		sessions.POST("/", func(c *gin.Context) {
			newSession := gin.H{
				"id":         "session-new-" + time.Now().Format("20060102150405"),
				"user_id":    "mock-user-id-123",
				"start_time": time.Now().Format(time.RFC3339),
				"end_time":   time.Now().Add(10 * time.Minute).Format(time.RFC3339),
				"duration":   10,
				"type":       "timer",
				"created_at": time.Now().Format(time.RFC3339),
			}

			c.JSON(http.StatusCreated, newSession)
		})

		sessions.GET("/stats", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"total_sessions":    5,
				"total_minutes":     47,
				"current_streak":    3,
				"longest_streak":    7,
				"weekly_sessions":   3,
				"monthly_sessions":  5,
				"average_duration":  9.4,
				"most_common_mood":  "calm",
			})
		})

		sessions.GET("/streak", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"current_streak": 3,
				"longest_streak": 7,
				"last_session":   time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
			})
		})
	}

	log.Println("Mock server starting on :8081")
	r.Run(":8081")
}