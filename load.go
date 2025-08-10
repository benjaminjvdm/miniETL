package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func Load(config Config, data []map[string]interface{}) error {
	switch config.Output.Type {
	case "csv":
		return loadToCSV(config.Output.Path, data)
	case "json":
		return loadToJSON(config.Output.Path, data)
	case "txt":
		return loadToTXT(config.Output.Path, data)
	case "sqlite":
		return loadToSQLite(config.Output.Path, data)
	default:
		return ErrInvalidOutputType
	}
}

func loadToCSV(path string, data []map[string]interface{}) error {
	if len(data) == 0 {
		// Create an empty file
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		f.Close()
		return nil
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	// Get headers from the first row and sort them for consistent order
	var headers []string
	for k := range data[0] {
		headers = append(headers, k)
	}
	sort.Strings(headers)

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

	for _, row := range data {
		jsonRow, err := json.Marshal(row)
		if err != nil {
			return err
		}
		if _, err := f.Write(append(jsonRow, '\n')); err != nil {
			return err
		}
	}

	return nil
}

func loadToSQLite(path string, data []map[string]interface{}) error {
	if len(data) == 0 {
		return nil // Nothing to load
	}

	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory for sqlite db: %w", err)
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	defer db.Close()

	// Infer table name from path, default to "data"
	tableName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	if tableName == "" {
		tableName = "data"
	}

	// Get headers and infer types from the first row
	var headers []string
	firstRow := data[0]
	for k := range firstRow {
		headers = append(headers, k)
	}
	sort.Strings(headers) // Sort headers for deterministic order

	var colTypes []string
	for _, k := range headers {
		colTypes = append(colTypes, getSQLiteType(firstRow[k]))
	}

	// Create table
	var createTableBuilder strings.Builder
	createTableBuilder.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS \"%s\" (", tableName))
	for i, header := range headers {
		createTableBuilder.WriteString(fmt.Sprintf("\"%s\" %s", header, colTypes[i]))
		if i < len(headers)-1 {
			createTableBuilder.WriteString(", ")
		}
	}
	createTableBuilder.WriteString(");")
	_, err = db.Exec(createTableBuilder.String())
	if err != nil {
		return fmt.Errorf("create table failed: %w", err)
	}

	// Prepare insert statement
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	var stmtBuilder strings.Builder
	stmtBuilder.WriteString(fmt.Sprintf("INSERT INTO \"%s\" (", tableName))
	quotedHeaders := make([]string, len(headers))
	for i, h := range headers {
		quotedHeaders[i] = fmt.Sprintf("\"%s\"", h)
	}
	stmtBuilder.WriteString(strings.Join(quotedHeaders, ", "))
	stmtBuilder.WriteString(") VALUES (")
	stmtBuilder.WriteString(strings.Repeat("?,", len(headers)-1) + "?)")
	stmt, err := tx.Prepare(stmtBuilder.String())
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("prepare statement failed: %w", err)
	}
	defer stmt.Close()

	// Insert data
	for _, row := range data {
		values := make([]interface{}, len(headers))
		for i, header := range headers {
			values[i] = row[header]
		}
		_, err := stmt.Exec(values...)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("insert row failed: %w", err)
		}
	}

	return tx.Commit()
}

func getSQLiteType(value interface{}) string {
	switch value.(type) {
	case int, int32, int64, bool:
		return "INTEGER"
	case float32, float64:
		return "REAL"
	case []byte:
		return "BLOB"
	default:
		return "TEXT"
	}
}
