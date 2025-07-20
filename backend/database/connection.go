package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var database *mongo.Database

// Connect MongoDBに接続
func Connect(ctx context.Context) (*mongo.Database, error) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	// Mock mode for development
	if uri == "mock" {
		log.Println("Using mock database for development")
		// Return a mock database that doesn't actually connect
		database = &mongo.Database{}
		return database, nil
	}

	// 接続オプション
	clientOptions := options.Client().
		ApplyURI(uri).
		SetConnectTimeout(10 * time.Second).
		SetMaxPoolSize(10).
		SetMinPoolSize(2)

	// MongoDB接続
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// 接続確認
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")

	// データベース選択
	database = client.Database("mindfulme")
	return database, nil
}

// Disconnect MongoDB接続を切断
func Disconnect(ctx context.Context) error {
	if client != nil {
		return client.Disconnect(ctx)
	}
	return nil
}

// Ping 接続確認
func Ping(ctx context.Context) error {
	if client == nil {
		return mongo.ErrClientDisconnected
	}
	return client.Ping(ctx, nil)
}

// GetDatabase データベースを取得
func GetDatabase() *mongo.Database {
	return database
}