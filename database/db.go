package database

import (
	"database/sql"
	"fmt"
)

type Database struct {
	DB *sql.DB
}

var DbInstance *Database

func InitDB(user, password, host, port, dbName string) error {
	dbURI := fmt.Sprintf("host=%s user=%s password=%s port=%s database=%s",
		host, user, password, port, dbName)
	db, err := sql.Open("pgx", dbURI)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	DbInstance = &Database{DB: db}
	return nil
}
