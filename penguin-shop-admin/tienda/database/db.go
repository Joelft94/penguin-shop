package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB(mongoURI string) {
	if mongoURI == "" {
		// Change this to match your Node.js admin panel's MongoDB port
		mongoURI = "mongodb://localhost:27017"
	}

	log.Printf("Connecting to MongoDB at: %s", mongoURI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	DB = client.Database("penguin-shop")
	
	// Test the products collection
	collection := DB.Collection("products")
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("Error counting products: %v", err)
	} else {
		log.Printf("Found %d products in database", count)
	}

	log.Printf("Connected to MongoDB database: %s", DB.Name())
}