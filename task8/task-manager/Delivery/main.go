package main

import (
	"fmt"
	"log"
	"os" // For potential future port configuration from env
	"task_manager/Delivery/routers" // Import the router setup function
)

func main() {
	fmt.Println("Starting Task Manager API...")

	// Initialize router with all dependencies wired inside SetupRouter
	router := routers.SetupRouter()

	// --- Server Configuration ---
	port := os.Getenv("PORT") // Example: Get port from environment variable
	if port == "" {
		port = "8080" // Default port if not set
		log.Printf("Defaulting to port %s", port)
	}

	serverAddr := ":" + port
	log.Printf("Server starting on %s", serverAddr)

	// --- Start Server ---
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}