package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	dsn := "postgres://postgres:student@localhost:5432/bookstore?sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to open db:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("failed to connect to db:", err)
	}

	migrate()
}

func migrate() {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS favorite_books (
			user_id    INTEGER NOT NULL,
			book_id    INTEGER NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			PRIMARY KEY (user_id, book_id)
		)
	`)
	if err != nil {
		log.Fatal("migration failed:", err)
	}
}
