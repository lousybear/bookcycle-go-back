package db

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/lousybear/bookcycle-go-back/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	once        sync.Once
)

func Init() {
	once.Do(func() {
		mongoURI := utils.GetEnv("MONGO_URI", "")
		if mongoURI == "" {
			log.Fatal("MONGO_URI is not set in the environment")
		}

		clientOpts := options.Client().
			ApplyURI(mongoURI).
			SetConnectTimeout(10 * time.Second).
			SetServerSelectionTimeout(10 * time.Second).
			SetSocketTimeout(10 * time.Second).
			SetRetryWrites(true)

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, clientOpts)
		if err != nil {
			log.Fatalf("Failed to create MongoDB client: %v", err)
		}

		if err := client.Ping(ctx, nil); err != nil {
			log.Fatalf("Failed to ping MongoDB: %v", err)
		}

		log.Println("✅ MongoDB connected successfully")
		mongoClient = client
	})
}

func Client() *mongo.Client {
	if mongoClient == nil {
		log.Fatal("MongoDB client not initialized. Call db.Init() in main.go")
	}
	return mongoClient
}

func GetCollection(name string) *mongo.Collection {
	dbName := utils.GetEnv("MONGO_DB", "test")
	return Client().Database(dbName).Collection(name)
}

func Disconnect() {
	if mongoClient == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := mongoClient.Disconnect(ctx); err != nil {
		log.Printf("⚠️ Error disconnecting MongoDB: %v", err)
	} else {
		log.Println("🔌 MongoDB disconnected cleanly")
	}
}
