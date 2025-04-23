package main

import (
	"log"

	"github.com/zain-0/khidmaat-backend/config"
	"github.com/zain-0/khidmaat-backend/routers"
)

func main() {
	// Initialize the MongoDB connection
	config.ConnectMongoDB()

	// Set up the router
	router := routers.SetupRouter()

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
