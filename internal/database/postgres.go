package database

import (
	"database/sql"
	"log"
	"time"
)

func NewPostgres() *sql.DB {
	conn := "host=localhost port=5432 user=postgres dbname=twitter sslmode=disable"
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db

}
