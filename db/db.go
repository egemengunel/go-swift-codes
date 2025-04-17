package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

//Init DB initilaizes the SQLite DB connection

func InitDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("error openind db connection %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging db %w", err)
	}

	//create the table
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS swift_codes (
		country_iso2 TEXT,
		swift_code   TEXT PRIMARY KEY,
		code_type    TEXT,
		name         TEXT,
		address      TEXT,
		town_name    TEXT,
		country_name TEXT,
		time_zone    TEXT
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating swift_codes table %w", err)
	}
	return db, nil
}
