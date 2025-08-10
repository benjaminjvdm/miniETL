package main

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	// Sample data
	data := []map[string]interface{}{
		{"field1": "value1", "field2": "value2"},
	}

	// Create a temporary file for the output
	tmpfile, err := os.CreateTemp("", "output.*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Create a config struct with the output path
	config := Config{
		Output: struct {
			Path string `mapstructure:"path"`
			Type string `mapstructure:"type"`
		}{Path: tmpfile.Name(), Type: "json"},
	}

	// Call Load
	if err := Load(config, data); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Read the temporary file
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Unmarshal the content and check if it's correct
	var loadedData []map[string]interface{}
	if err := json.Unmarshal(content, &loadedData); err != nil {
		t.Fatalf("Error unmarshalling loaded data: %v", err)
	}

	if !reflect.DeepEqual(data, loadedData) {
		t.Errorf("Load() produced %v, want %v", loadedData, data)
	}
}
