package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"mindfulme-backend/database"
	"mindfulme-backend/handlers"
	"mindfulme-backend/middleware"
	"mindfulme-backend/models"
)

func main() {
	// 環境変数の読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// MongoDB接続
	ctx := context.Background()
	db, err := database.Connect(ctx)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Disconnect(ctx)

	// サービスの初期化
	userService := models.NewUserService(db)
	sessionService := models.NewSessionService(db)

	// インデックスの作成
	if err := userService.CreateIndexes(ctx); err != nil {
		log.Println("Warning: Failed to create user indexes:", err)
	}
	if err := sessionService.CreateIndexes(ctx); err != nil {
		log.Println("Warning: Failed to create session indexes:", err)
	}

	// Ginルーターの設定
	router := gin.Default()

	// CORS設定
	router.Use(middleware.CORSMiddleware())

	// ヘルスチェック
	router.GET("/health", func(c *gin.Context) {
		dbStatus := "ok"
		if err := database.Ping(ctx); err != nil {
			dbStatus = "error"
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"app":    "MindfulMe",
			"time":   time.Now().Format(time.RFC3339),
			"database": gin.H{
				"status": dbStatus,
			},
		})
	})

	// APIルートグループ
	api := router.Group("/api/v1")
	{
		// 認証エンドポイント
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register(userService))
			auth.POST("/login", handlers.Login(userService))
			auth.POST("/refresh", handlers.RefreshToken())
			auth.POST("/logout", middleware.AuthMiddleware(), handlers.Logout())
		}

		// 瞑想セッションエンドポイント（認証必須）
		sessions := api.Group("/sessions")
		sessions.Use(middleware.AuthMiddleware())
		{
			sessions.GET("", handlers.GetSessions(sessionService))
			sessions.POST("", handlers.CreateSession(sessionService))
			sessions.PUT("/:id", handlers.UpdateSession(sessionService))
			sessions.GET("/stats", handlers.GetStats(sessionService))
			sessions.GET("/streak", handlers.GetStreak(sessionService))
		}

		// ユーザープロファイル（認証必須）
		profile := api.Group("/profile")
		profile.Use(middleware.AuthMiddleware())
		{
			profile.GET("", handlers.GetProfile(userService))
			profile.PUT("", handlers.UpdateProfile(userService))
		}
	}

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// グレースフルシャットダウン
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Printf("MindfulMe server started on port %s", port)

	// シグナル待機
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}