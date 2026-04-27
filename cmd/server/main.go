package main

import (
	"log"
	"ride-hailing/internal/config"
	"ride-hailing/internal/router"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
	config.ConnectDB()
	config.ConnectRedis()
	config.InitNewRelic()

	r := router.SetupRouter()
	r.Run(":8080")
}