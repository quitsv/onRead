package Controllers

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Connect() *sql.DB {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	db_url := os.Getenv("DATABASE")
	db, err := sql.Open("mysql", db_url)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
