package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var UsersCollection *mongo.Collection
var MedicalRecordsCollection *mongo.Collection

func ConnectMongoDB() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MongoDB URI is not set in the .env file")
	}

	// MongoDB client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Set a timeout for connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new MongoDB client and connect to the database
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("MongoDB Connection Error: %v\n", err)
	}

	// Ensure that the connection is established
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v\n", err)
	}

	// Use your actual database name here
	DB = client.Database("khidmaat")

	// Initialize collections
	UsersCollection = DB.Collection("users")
	MedicalRecordsCollection = DB.Collection("medical_records")

	log.Println("✅ Connected to MongoDB Atlas")
}
