package db

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- Test createTables() directly ---
func TestCreateTables(t *testing.T) {
	var err error

	DB, err = sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err, "Failed to open in-memory database")

	defer DB.Close()

	err = createTables()
	assert.NoError(t, err, "Should create tables without error")

	// Verify tables exist
	tables := []string{"users", "events", "registration"}
	for _, table := range tables {
		var name string
		err := DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		assert.NoError(t, err, "table %s should exist", table)
		assert.Equal(t, table, name)
	}
}

// --- Test InitDB() end-to-end ---
func TestInitDB(t *testing.T) {
	// Create a temp file DB name
	tempDBFile := "test_api.db"
	defer os.Remove(tempDBFile) // cleanup after test

	// Temporarily override the file-based DB
	oldDB := DB
	defer func() { DB = oldDB }()

	DB, _ = sql.Open("sqlite3", tempDBFile)
	err := createTables()
	assert.NoError(t, err, "InitDB should create tables successfully")

	// Check if tables exist
	tables := []string{"users", "events", "registration"}
	for _, table := range tables {
		var name string
		err := DB.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		assert.NoError(t, err, "table %s should exist after InitDB", table)
	}
}
