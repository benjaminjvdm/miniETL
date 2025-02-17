package main

import (
	"testing"
)

func TestExtract(t *testing.T) {
	// Test cases for Extract function
	// Test case 1: CSV input
	t.Run("CSV Input", func(t *testing.T) {
		config := Config{
			Input: struct {
				Path      string `mapstructure:"path"`
				Type      string `mapstructure:"type"`
				Delimiter string `mapstructure:"delimiter"`
			}{
				Path:      "input.csv",
				Type:      "csv",
				Delimiter: ",",
			},
		}
		_, err := Extract(config)
		if err != nil {
			t.Errorf("Error extracting data: %v", err)
		}
	})

	// Test case 2: JSON input
	t.Run("JSON Input", func(t *testing.T) {
		// Implement test case for JSON input
		config := Config{
			Input: struct {
				Path      string `mapstructure:"path"`
				Type      string `mapstructure:"type"`
				Delimiter string `mapstructure:"delimiter"`
			}{
				Path:      "output.json",
				Type:      "json",
				Delimiter: ",",
			},
		}
		_, err := Extract(config)
		if err != nil {
			t.Errorf("Error extracting data: %v", err)
		}
	})

	// Test case 3: TXT input
	t.Run("TXT Input", func(t *testing.T) {
		// Implement test case for TXT input
		config := Config{
			Input: struct {
				Path      string `mapstructure:"path"`
				Type      string `mapstructure:"type"`
				Delimiter string `mapstructure:"delimiter"`
			}{
				Path:      "input.txt",
				Type:      "txt",
				Delimiter: ",",
			},
		}
		_, err := Extract(config)
		if err != nil {
			t.Errorf("Error extracting data: %v", err)
		}
	})
}
