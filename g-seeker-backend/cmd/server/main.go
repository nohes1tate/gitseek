package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"g-seeker-backend/internal/router"
)

func main() {
	// Load .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	r := gin.Default()

	router.RegisterRoutes(r)

	// Get port from environment variable, default to 3000
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "3000"
	}
	addr := ":" + port
	log.Printf("server is running at %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
