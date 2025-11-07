package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error

	DB, err = sql.Open("sqlite3", "api.db")

	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	// Verify database connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	err = createTables()
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
}

func createTables() error {
	createUserTable := `
     CREATE TABLE IF NOT EXISTS users(
         id INTEGER PRIMARY KEY AUTOINCREMENT,
         email TEXT NOT NULL UNIQUE,
         password TEXT NOT NULL
     )`

	_, err := DB.Exec(createUserTable)
	if err != nil {
		return fmt.Errorf("could not create users table: %v", err)
	}

	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		dateTime DATETIME NOT NULL,
		user_id INTEGER,
	    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`

	_, err = DB.Exec(createEventsTable)
	if err != nil {
		return fmt.Errorf("could not create events table: %v", err)
	}

	createRegistrationsTable := `
	CREATE TABLE IF NOT EXISTS registration (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    event_id INTEGER,
	    user_id INTEGER,
	    FOREIGN KEY(event_id) REFERENCES events(id) ON DELETE CASCADE,
	    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
	)`

	_, err = DB.Exec(createRegistrationsTable)
	if err != nil {
		return fmt.Errorf("could not create registrations table: %v", err)
	}

	return nil
}
