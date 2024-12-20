package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewPostgresConnection() (*sql.DB, error) {
	connStr := "postgres://postgres:@localhost:5432/recordings?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
