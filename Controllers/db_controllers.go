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

	db, err := sql.Open("mysql", os.Getenv("DATABASE"))
	if err != nil {
		log.Fatal(err)
	}
	return db
}
