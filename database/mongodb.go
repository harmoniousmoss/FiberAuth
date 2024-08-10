package database

import (
	"context"
	"fmt"
	"log"
	"myfibergotemplate/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func ConnectMongoDB() error {
	log.Println("Attempting to connect to MongoDB...")
	mongoURI := config.GetEnv("MONGO_URI", "")
	if mongoURI == "" {
		log.Println("MONGO_URI environment variable is not set.")
		return fmt.Errorf("MONGO_URI environment variable is not set")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	var err error

	for attempt := 1; attempt <= 5; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		MongoClient, err = mongo.Connect(ctx, clientOptions)
		if err == nil {
			err = MongoClient.Ping(ctx, nil)
			if err == nil {
				log.Println("Connected to MongoDB!")
				cancel()
				return nil
			}
		}

		log.Printf("Failed to connect to MongoDB: %v, retrying in %d seconds", err, attempt*2)
		time.Sleep(time.Duration(attempt*2) * time.Second)
		cancel()
	}

	return fmt.Errorf("failed to connect to MongoDB after several attempts: %v", err)
}

func GetMongoClient() *mongo.Client {
	return MongoClient
}
