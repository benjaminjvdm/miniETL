package main

import (
	"os"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	content := []byte(`
input:
  path: testdata/input.csv
  type: csv
output:
  path: testdata/output.json
  type: json
transformations:
  - new_col: old_col
`)
	tmpfile, err := os.CreateTemp("", "config.*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Call LoadConfig with the path to the temporary file
	config, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Check if the returned Config struct has the expected values
	expected := Config{
		Input: struct {
			Path string `mapstructure:"path"`
			Type string `mapstructure:"type"`
		}{Path: "testdata/input.csv", Type: "csv"},
		Output: struct {
			Path string `mapstructure:"path"`
			Type string `mapstructure:"type"`
		}{Path: "testdata/output.json", Type: "json"},
		Transformations: []map[string]interface{}{
			{"new_col": "old_col"},
		},
	}

	if !reflect.DeepEqual(config, expected) {
		t.Errorf("LoadConfig() = %v, want %v", config, expected)
	}
}
