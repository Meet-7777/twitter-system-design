package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgres() *sql.DB {
	conn := "host=localhost port=5432 user=dishabohra dbname=twitter sslmode=disable"

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Postgres connected 🚀")
	return db
}
