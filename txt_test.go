package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestExtractFromTXT(t *testing.T) {
	// Create a temporary directory for test data
	tmpDir, err := os.MkdirTemp("", "testdata")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a temporary TXT file with JSON objects per line
	txtContent := []byte(`{"col1": "val1", "col2": 123}` + "\n" + `{"col1": "val2", "col2": 456}` + "\n")
	txtPath := filepath.Join(tmpDir, "input.txt")
	if err := os.WriteFile(txtPath, txtContent, 0666); err != nil {
		t.Fatal(err)
	}

	// Call extractFromTXT
	data, err := extractFromTXT(txtPath)
	if err != nil {
		t.Fatalf("extractFromTXT() error = %v", err)
	}

	// Check if the returned data is correct
	expected := []map[string]interface{}{
		{"col1": "val1", "col2": float64(123)}, // Note: json.Unmarshal reads numbers as float64
		{"col1": "val2", "col2": float64(456)},
	}

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("extractFromTXT() = %v, want %v", data, expected)
	}
}

func TestLoadToTXT(t *testing.T) {
	// Sample data
	data := []map[string]interface{}{
		{"fieldA": "valueA", "fieldB": 111},
		{"fieldA": "valueB", "fieldB": 222},
	}

	// Create a temporary file for the output
	tmpfile, err := os.CreateTemp("", "output.*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Call loadToTXT
	if err := loadToTXT(tmpfile.Name(), data); err != nil {
		t.Fatalf("loadToTXT() error = %v", err)
	}

	// Read the temporary file
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Expected content
	expected := `{"fieldA":"valueA","fieldB":111}` + "\n" + `{"fieldA":"valueB","fieldB":222}` + "\n"

	if string(content) != expected {
		t.Errorf("loadToTXT() produced %q, want %q", string(content), expected)
	}
}
