package main

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestLoadToSQLite(t *testing.T) {
	// Sample data
	data := []map[string]interface{}{
		{"name": "Alice", "age": 30},
		{"name": "Bob", "age": 25},
	}

	// Create a temporary file for the output
	tmpfile, err := os.CreateTemp("", "output.*.db")
	if err != nil {
		t.Fatal(err)
	}
	dbPath := tmpfile.Name()
	tmpfile.Close() // Close the file so the db driver can use it
	defer os.Remove(dbPath) // clean up

	// Call loadToSQLite
	if err := loadToSQLite(dbPath, data); err != nil {
		t.Fatalf("loadToSQLite() error = %v", err)
	}

	// Connect to the database and verify the data
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	tableName := strings.TrimSuffix(filepath.Base(dbPath), filepath.Ext(dbPath))
	query := "SELECT name, age FROM \"" + tableName + "\" ORDER BY name"
	rows, err := db.Query(query)
	if err != nil {
		t.Fatalf("Failed to query test database with query '%s': %v", query, err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var name string
		var age int
		if err := rows.Scan(&name, &age); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		results = append(results, map[string]interface{}{"name": name, "age": age})
	}

	expectedCount := 2
	if len(results) != expectedCount {
		t.Errorf("Expected %d rows, got %d", expectedCount, len(results))
	}

	// Note: reflect.DeepEqual is tricky with database types. Manual check is safer.
	if results[0]["name"] != "Alice" || results[0]["age"] != 30 {
		t.Errorf("Unexpected row 1: got %v", results[0])
	}
	if results[1]["name"] != "Bob" || results[1]["age"] != 25 {
		t.Errorf("Unexpected row 2: got %v", results[1])
	}
}
