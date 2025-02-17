package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

func Load(config Config, data []map[string]interface{}) error {
	switch config.Output.Type {
	case "csv":
		err := loadToCSV(config.Output.Path, data)
		if err != nil {
			ErrorLogger.Println("Error loading to CSV:", err)
			return err
		}
		return nil
	case "json":
		err := loadToJSON(config.Output.Path, data)
		if err != nil {
			ErrorLogger.Println("Error loading to JSON:", err)
			DebugLogger.Println("Error loading to JSON:", err)
			return err
		}
		return nil
	case "txt":
		err := loadToTXT(config.Output.Path, data)
		if err != nil {
			ErrorLogger.Println("Error loading to TXT:", err)
			return err
		}
		return nil
	case "sqlite":
		err := loadToSQLite(config, data)
		if err != nil {
			ErrorLogger.Println("Error loading to SQLite:", err)
			return err
		}
		return nil
	default:
		return fmt.Errorf("unsupported output type: %s", config.Output.Type)
	}
}

func loadToCSV(path string, data []map[string]interface{}) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	// Get headers
	var headers []string
	for k := range data[0] {
		headers = append(headers, k)
	}

	err = writer.Write(headers)
	if err != nil {
		return err
	}

	for _, row := range data {
		var record []string
		for _, header := range headers {
			record = append(record, fmt.Sprintf("%v", row[header]))
		}
		err = writer.Write(record)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadToJSON(path string, data []map[string]interface{}) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err) // Add this line
		return err
	}

	return nil
}

func loadToTXT(path string, data []map[string]interface{}) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Load the config to get the delimiter
	config, err := LoadConfig(".")
	if err != nil {
		return err
	}

	delimiter := config.Input.Delimiter
	if delimiter == "" {
		delimiter = "," // Default delimiter
	}

	for _, row := range data {
		var record []string
		for _, value := range row {
			record = append(record, fmt.Sprintf("%v", value))
		}
		_, err = fmt.Fprintln(f, strings.Join(record, delimiter))
		if err != nil {
			return err
		}
	}

	return nil
}

func loadToSQLite(config Config, data []map[string]interface{}) error {
	db, err := sql.Open("sqlite3", config.Output.Path)
	if err != nil {
		return err
	}
	defer db.Close()

	// Get headers
	var headers []string
	for k := range data[0] {
		headers = append(headers, k)
	}

	// Create table
	var createTableSQL string
	createTableSQL += "CREATE TABLE IF NOT EXISTS data ("
	for i, header := range headers {
		createTableSQL += header + " TEXT"
		if i < len(headers)-1 {
			createTableSQL += ", "
		}
	}
	createTableSQL += ");"

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	// Prepare statement for inserting data
	var insertSQL string
	insertSQL += "INSERT INTO data ("
	for i, header := range headers {
		insertSQL += header
		if i < len(headers)-1 {
			insertSQL += ", "
		}
	}
	insertSQL += ") VALUES ("
	for i := 0; i < len(headers); i++ {
		insertSQL += "?"
		if i < len(headers)-1 {
			insertSQL += ", "
		}
	}
	insertSQL += ");"

	stmt, err := db.Prepare(insertSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Insert data
	for _, row := range data {
		var values []interface{}
		for _, header := range headers {
			values = append(values, row[header])
		}
		_, err = stmt.Exec(values...)
		if err != nil {
			return err
		}
	}

	return nil
}
