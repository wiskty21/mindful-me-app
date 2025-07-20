package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/argon2"
	"mindfulme-backend/models"
)

// RegisterRequest ユーザー登録リクエスト
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

// LoginRequest ログインリクエスト
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse 認証レスポンス
type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

// UserResponse ユーザー情報レスポンス
type UserResponse struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	Name          string     `json:"name"`
	TotalSessions int        `json:"total_sessions"`
	TotalMinutes  int        `json:"total_minutes"`
	CurrentStreak int        `json:"current_streak"`
	LastSessionAt *time.Time `json:"last_session_at"`
}

// generateSalt ランダムソルトを生成
func generateSalt() ([]byte, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	return salt, err
}

// hashPassword パスワードをハッシュ化
func hashPassword(password string, salt []byte) string {
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return base64.StdEncoding.EncodeToString(hash)
}

// verifyPassword パスワードを検証
func verifyPassword(password, hash string, salt []byte) bool {
	passwordHash := hashPassword(password, salt)
	return passwordHash == hash
}

// generateToken JWTトークンを生成
func generateToken(userID, email string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(duration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// Register ユーザー登録ハンドラー
func Register(userService *models.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := context.Background()

		// 既存ユーザーチェック
		existingUser, _ := userService.GetUserByEmail(ctx, req.Email)
		if existingUser != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}

		// パスワードハッシュ化
		salt, err := generateSalt()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate salt"})
			return
		}

		passwordHash := hashPassword(req.Password, salt)

		// ユーザー作成
		user := &models.User{
			Email:        req.Email,
			Name:         req.Name,
			PasswordHash: passwordHash + ":" + base64.StdEncoding.EncodeToString(salt),
		}

		if err := userService.CreateUser(ctx, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		// トークン生成
		accessToken, err := generateToken(user.ID.Hex(), user.Email, 15*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
			return
		}

		refreshToken, err := generateToken(user.ID.Hex(), user.Email, 7*24*time.Hour)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
			return
		}

		c.JSON(http.StatusCreated, AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			User: UserResponse{
				ID:            user.ID.Hex(),
				Email:         user.Email,
				Name:          user.Name,
				TotalSessions: 0,
				TotalMinutes:  0,
				CurrentStreak: 0,
			},
		})
	}
}

// Login ログインハンドラー
func Login(userService *models.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := context.Background()

		// ユーザー検索
		user, err := userService.GetUserByEmail(ctx, req.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// パスワード検証
		hashParts := []byte(user.PasswordHash)
		separatorIndex := len(hashParts) - 45 // base64エンコードされた32バイトのソルト
		hash := string(hashParts[:separatorIndex])
		saltBase64 := string(hashParts[separatorIndex+1:])

		salt, err := base64.StdEncoding.DecodeString(saltBase64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid password format"})
			return
		}

		if !verifyPassword(req.Password, hash, salt) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// トークン生成
		accessToken, err := generateToken(user.ID.Hex(), user.Email, 15*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
			return
		}

		refreshToken, err := generateToken(user.ID.Hex(), user.Email, 7*24*time.Hour)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
			return
		}

		c.JSON(http.StatusOK, AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			User: UserResponse{
				ID:            user.ID.Hex(),
				Email:         user.Email,
				Name:          user.Name,
				TotalSessions: user.TotalSessions,
				TotalMinutes:  user.TotalMinutes,
				CurrentStreak: user.CurrentStreak,
				LastSessionAt: user.LastSessionAt,
			},
		})
	}
}

// RefreshToken リフレッシュトークンハンドラー
func RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// トークン検証
		token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		userID := claims["user_id"].(string)
		email := claims["email"].(string)

		// 新しいトークン生成
		accessToken, err := generateToken(userID, email, 15*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token": accessToken,
		})
	}
}

// Logout ログアウトハンドラー
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		// JWTトークンはステートレスなので、クライアント側でトークンを削除するだけ
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	}
}

// GetProfile プロフィール取得ハンドラー
func GetProfile(userService *models.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetString("userID")
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		ctx := context.Background()
		user, err := userService.GetUserByID(ctx, userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user": UserResponse{
				ID:            user.ID.Hex(),
				Email:         user.Email,
				Name:          user.Name,
				TotalSessions: user.TotalSessions,
				TotalMinutes:  user.TotalMinutes,
				CurrentStreak: user.CurrentStreak,
				LastSessionAt: user.LastSessionAt,
			},
			"preferences": gin.H{
				"preferred_time": user.PreferredTime,
				"reminder_time":  user.ReminderTime,
			},
		})
	}
}

// UpdateProfile プロフィール更新ハンドラー
func UpdateProfile(userService *models.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetString("userID")
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var req struct {
			Name          string `json:"name"`
			PreferredTime int    `json:"preferred_time"`
			ReminderTime  string `json:"reminder_time"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := context.Background()
		if err := userService.UpdateProfile(ctx, userID, req.Name, req.PreferredTime, req.ReminderTime); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
	}
}