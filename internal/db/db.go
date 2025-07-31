package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// NewConnection establishes a new connection to the PostgreSQL database.
// It returns a pointer to the sql.DB object or an error if the connection fails.
func NewConnection() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=postgres password=mypassword dbname=computers sslmode=disable"

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to DB: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to the database.")

	return conn, nil
}
