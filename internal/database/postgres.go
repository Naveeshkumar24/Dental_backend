package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func Connect(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}
