package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"g-seeker-backend/internal/router"
)

func main() {
	r := gin.Default()

	router.RegisterRoutes(r)

	addr := ":6657"
	log.Printf("server is running at %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
