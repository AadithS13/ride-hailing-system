package main

import (
	"log"
	"ride-hailing/internal/config"
	"ride-hailing/internal/router"

	"github.com/joho/godotenv"
)

func main() {

	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	} else {
		log.Println("Loaded .env file")
	}

	// DB
	log.Println("Connecting to DB...")
	config.ConnectDB()
	log.Println("DB connected")

	// Migrations
	config.RunMigrations()
	log.Println("Migrations done")

	// Redis
	log.Println("Connecting to Redis...")
	config.ConnectRedis()
	log.Println("Redis connected")

	// New Relic
	config.InitNewRelic()
	log.Println("New Relic initialized")

	// Server
	r := router.SetupRouter()

	log.Println("Server running on :8080")
	r.Run(":8080")
}