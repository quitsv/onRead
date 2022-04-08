package controllers

import (
	"database/sql"
	"log"
)

func Connect() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/onread?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
