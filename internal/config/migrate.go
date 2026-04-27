package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/pressly/goose/v3"
	_ "github.com/lib/pq"
)

func RunMigrations() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to open DB for migrations:", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal("failed to run migrations:", err)
	}

	log.Println("Migrations applied ✅")
}