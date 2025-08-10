package main

import (
	"reflect"
	"testing"
)

func TestTransformWithSimpleParser(t *testing.T) {
	// Sample data
	data := []map[string]interface{}{
		{"header1": "value1", "header2": "value2"},
	}

	// Transformations config
	transformations := []map[string]interface{}{
		{"combined": `header1 + " - " + header2`},
		{"upper_h1": `upper(header1)`},
		{"total_len": `len(header1) + len(header2)`},
	}

	// Call Transform
	transformedData, err := Transform(data, transformations)
	if err != nil {
		t.Fatalf("Transform() error = %v", err)
	}

	// Expected data
	expected := []map[string]interface{}{
		{
			"header1":    "value1",
			"header2":    "value2",
			"combined":   "value1 - value2",
			"upper_h1":   "VALUE1",
			"total_len":  12, // len("value1") + len("value2") = 6 + 6
		},
	}

	// Check if the data was transformed correctly
	if !reflect.DeepEqual(transformedData, expected) {
		t.Errorf("Transform() = \n%v, want \n%v", transformedData, expected)
	}
}
