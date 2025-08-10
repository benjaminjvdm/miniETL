package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestExtract(t *testing.T) {
	// Create a temporary directory for test data
	tmpDir, err := os.MkdirTemp("", "testdata")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a temporary CSV file
	csvContent := []byte("header1,header2\nvalue1,value2\nvalue3,value4\n")
	csvPath := filepath.Join(tmpDir, "input.csv")
	if err := os.WriteFile(csvPath, csvContent, 0666); err != nil {
		t.Fatal(err)
	}

	// Create a config struct pointing to the temporary file
	config := Config{
		Input: struct {
			Path string `mapstructure:"path"`
			Type string `mapstructure:"type"`
		}{Path: csvPath, Type: "csv"},
	}

	// Call Extract
	data, err := Extract(config)
	if err != nil {
		t.Fatalf("Extract() error = %v", err)
	}

	// Check if the returned data is correct
	expected := []map[string]interface{}{
		{"header1": "value1", "header2": "value2"},
		{"header1": "value3", "header2": "value4"},
	}

	if !reflect.DeepEqual(data, expected) {
		t.Errorf("Extract() = %v, want %v", data, expected)
	}
}
