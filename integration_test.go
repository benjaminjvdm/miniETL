package main

import (
	"database/sql"
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Test Cases ---

func TestIntegration_CSV_to_SQLite(t *testing.T) {
	// --- Setup ---
	// 1. Create temp directory
	tmpDir, err := os.MkdirTemp("", "etl_integration_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// 2. Create temp input CSV file
	inputPath := filepath.Join(tmpDir, "input.csv")
	csvContent := "id,name,value\n1,Alice,100\n2,Bob,200\n"
	err = os.WriteFile(inputPath, []byte(csvContent), 0644)
	require.NoError(t, err)

	// 3. Define output path for SQLite DB
	dbPath := filepath.Join(tmpDir, "output.db")

	// 4. Create temp YAML config file
	configContent := `
input:
  path: ` + inputPath + `
  type: csv
output:
  path: ` + dbPath + `
  type: sqlite
transformations:
  - name_upper: upper(name)
`
	configPath := filepath.Join(tmpDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// --- Execute ---
	err = runETL(configPath)
	require.NoError(t, err)

	// --- Verify ---
	// 1. Connect to the created SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	require.NoError(t, err)
	defer db.Close()

	// 2. Query the data
	tableName := strings.TrimSuffix(filepath.Base(dbPath), filepath.Ext(dbPath))
	rows, err := db.Query("SELECT id, name, value, name_upper FROM \"" + tableName + "\" ORDER BY id")
	require.NoError(t, err)
	defer rows.Close()

	// 3. Check the results
	var results []map[string]interface{}
	for rows.Next() {
		var id, value int
		var name, nameUpper string
		err := rows.Scan(&id, &name, &value, &nameUpper)
		require.NoError(t, err)
		results = append(results, map[string]interface{}{
			"id":         id,
			"name":       name,
			"value":      value,
			"name_upper": nameUpper,
		})
	}
	require.NoError(t, rows.Err())

	expected := []map[string]interface{}{
		{"id": 1, "name": "Alice", "value": 100, "name_upper": "ALICE"},
		{"id": 2, "name": "Bob", "value": 200, "name_upper": "BOB"},
	}

	assert.Equal(t, len(expected), len(results))
	// Manual checks because type comparison with DB results is tricky
	assert.Equal(t, expected[0]["name_upper"], results[0]["name_upper"])
	assert.Equal(t, expected[1]["name_upper"], results[1]["name_upper"])
}

func TestIntegration_TXT_to_CSV(t *testing.T) {
	// --- Setup ---
	// 1. Create temp directory
	tmpDir, err := os.MkdirTemp("", "etl_integration_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// 2. Create temp input TXT (JSONL) file
	inputPath := filepath.Join(tmpDir, "input.txt")
	txtContent := `{"user_id": 1, "email": "a@test.com"}` + "\n" + `{"user_id": 2, "email": "b@test.com"}` + "\n"
	err = os.WriteFile(inputPath, []byte(txtContent), 0644)
	require.NoError(t, err)

	// 3. Define output path for CSV
	outputPath := filepath.Join(tmpDir, "output.csv")

	// 4. Create temp YAML config file
	configContent := `
input:
  path: ` + inputPath + `
  type: txt
output:
  path: ` + outputPath + `
  type: csv
`
	configPath := filepath.Join(tmpDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// --- Execute ---
	err = runETL(configPath)
	require.NoError(t, err)

	// --- Verify ---
	// 1. Read the output CSV file
	file, err := os.Open(outputPath)
	require.NoError(t, err)
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// 2. Check the results
	// Note: Order of columns from map is not guaranteed
	assert.Equal(t, 3, len(records), "Expected 3 records (header + 2 data rows)")
	assert.Contains(t, records[0], "user_id")
	assert.Contains(t, records[0], "email")
	// Find which row is which
	var row1, row2 []string
	if records[1][0] == "1" || records[1][1] == "1" {
		row1 = records[1]
		row2 = records[2]
	} else {
		row1 = records[2]
		row2 = records[1]
	}
	assert.Contains(t, row1, "1")
	assert.Contains(t, row1, "a@test.com")
	assert.Contains(t, row2, "2")
	assert.Contains(t, row2, "b@test.com")
}

func TestIntegration_TXT_to_SQLite(t *testing.T) {
	// --- Setup ---
	tmpDir, err := os.MkdirTemp("", "etl_integration_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.txt")
	txtContent := `{"id": 1, "name": "Charlie"}` + "\n" + `{"id": 2, "name": "Diana"}` + "\n"
	err = os.WriteFile(inputPath, []byte(txtContent), 0644)
	require.NoError(t, err)

	dbPath := filepath.Join(tmpDir, "output.db")

	configContent := `
input:
  path: ` + inputPath + `
  type: txt
output:
  path: ` + dbPath + `
  type: sqlite
`
	configPath := filepath.Join(tmpDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// --- Execute ---
	err = runETL(configPath)
	require.NoError(t, err)

	// --- Verify ---
	db, err := sql.Open("sqlite3", dbPath)
	require.NoError(t, err)
	defer db.Close()

	tableName := strings.TrimSuffix(filepath.Base(dbPath), filepath.Ext(dbPath))
	rows, err := db.Query("SELECT id, name FROM \"" + tableName + "\" ORDER BY id")
	require.NoError(t, err)
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		require.NoError(t, err)
		results = append(results, map[string]interface{}{"id": id, "name": name})
	}
	require.NoError(t, rows.Err())

	expected := []map[string]interface{}{
		{"id": 1, "name": "Charlie"},
		{"id": 2, "name": "Diana"},
	}
	assert.Equal(t, expected[0]["name"], results[0]["name"])
}

func TestIntegration_TOML_Config(t *testing.T) {
	// --- Setup ---
	tmpDir, err := os.MkdirTemp("", "etl_integration_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.csv")
	csvContent := "header\nvalue\n"
	err = os.WriteFile(inputPath, []byte(csvContent), 0644)
	require.NoError(t, err)

	outputPath := filepath.Join(tmpDir, "output.json")

	// 4. Create temp TOML config file
	configContent := `
[input]
path = "` + inputPath + `"
type = "csv"

[output]
path = "` + outputPath + `"
type = "json"
`
	configPath := filepath.Join(tmpDir, "config.toml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// --- Execute ---
	err = runETL(configPath)
	require.NoError(t, err)

	// --- Verify ---
	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.JSONEq(t, `[{"header": "value"}]`, string(content))
}

// --- Error Handling Tests ---

func TestIntegration_Error_NonExistentInputFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "etl_integration_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	configContent := `
input:
  path: /path/to/non/existent/file.csv
  type: csv
output:
  path: ` + filepath.Join(tmpDir, "output.json") + `
  type: json
`
	configPath := filepath.Join(tmpDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	err = runETL(configPath)
	require.Error(t, err, "Expected an error for non-existent input file")
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestIntegration_Error_MalformedInputCSV(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "etl_integration_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.csv")
	// Jagged CSV: header has 2 columns, data row has 3
	csvContent := "h1,h2\nv1,v2,v3\n"
	err = os.WriteFile(inputPath, []byte(csvContent), 0644)
	require.NoError(t, err)

	configContent := `
input:
  path: ` + inputPath + `
  type: csv
output:
  path: ` + filepath.Join(tmpDir, "output.json") + `
  type: json
`
	configPath := filepath.Join(tmpDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	err = runETL(configPath)
	require.Error(t, err, "Expected an error for malformed CSV")
	assert.Contains(t, err.Error(), "record on line 2: wrong number of fields")
}

func TestIntegration_Error_InvalidTransformation(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "etl_integration_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.csv")
	csvContent := "header1\nvalue1\n"
	err = os.WriteFile(inputPath, []byte(csvContent), 0644)
	require.NoError(t, err)

	configContent := `
input:
  path: ` + inputPath + `
  type: csv
output:
  path: ` + filepath.Join(tmpDir, "output.json") + `
  type: json
transformations:
  - new_col: upper(non_existent_col)
`
	configPath := filepath.Join(tmpDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	err = runETL(configPath)
	require.Error(t, err, "Expected an error for invalid transformation")
	assert.Contains(t, err.Error(), "invalid field for upper(): non_existent_col")
}

func TestIntegration_Error_OutputToInvalidPath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "etl_integration_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.csv")
	csvContent := "h1\nv1\n"
	err = os.WriteFile(inputPath, []byte(csvContent), 0644)
	require.NoError(t, err)

	configContent := `
input:
  path: ` + inputPath + `
  type: csv
output:
  path: /non/existent/path/output.json
  type: json
`
	configPath := filepath.Join(tmpDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	err = runETL(configPath)
	require.Error(t, err, "Expected an error for invalid output path")
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestIntegration_Error_MalformedInputTXT(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "etl_integration_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.txt")
	// Malformed JSONL: second line is not valid JSON
	txtContent := `{"id":1}` + "\n" + `not-json` + "\n"
	err = os.WriteFile(inputPath, []byte(txtContent), 0644)
	require.NoError(t, err)

	configContent := `
input:
  path: ` + inputPath + `
  type: txt
output:
  path: ` + filepath.Join(tmpDir, "output.json") + `
  type: json
`
	configPath := filepath.Join(tmpDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	err = runETL(configPath)
	require.Error(t, err, "Expected an error for malformed TXT/JSONL")
	assert.Contains(t, err.Error(), "invalid character 'o' in literal null (expecting 'u')")
}

func TestIntegration_Error_InvalidConfigFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "etl_integration_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Config file is missing the 'output' section
	configContent := `
input:
  path: /some/path
  type: csv
`
	configPath := filepath.Join(tmpDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	err = runETL(configPath)
	require.Error(t, err, "Expected an error for invalid config file")
	// This error comes from our validation logic in LoadConfig
	assert.Contains(t, err.Error(), "validation failed")
	assert.Contains(t, err.Error(), "Output.Path is required")
}
