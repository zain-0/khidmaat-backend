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

	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("MongoDB Connection Error:", err)
	}

	DB = client.Database("your_db_name") // Replace with your actual DB name

	// Initialize collections
	UsersCollection = DB.Collection("users")
	MedicalRecordsCollection = DB.Collection("medical_records")

	log.Println("✅ Connected to MongoDB Atlas")
}
