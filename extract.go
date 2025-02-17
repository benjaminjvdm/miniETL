package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func Extract(config Config) ([]map[string]interface{}, error) {
	// var data []map[string]interface{} // Remove this line

	switch config.Input.Type {
	case "csv":
		data, err := extractFromCSV(config.Input.Path)
		if err != nil {
			fmt.Println("Error extracting from CSV:", err) // Add this line
			return nil, err
		}
		fmt.Println("CSV data:", data) // Print the data here
		return data, nil
	case "json":
		data, err := extractFromJSON(config.Input.Path)
		if err != nil {
			return nil, err
		}
		return data, nil
	case "txt":
		data, err := extractFromTXT(config.Input.Path)
		if err != nil {
			return nil, err
		}
		return data, nil
	default:
		return nil, fmt.Errorf("unsupported input type: %s", config.Input.Type)
	}
}

func extractFromCSV(path string) ([]map[string]interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening CSV file:", err) // Add this line
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	headers, err := reader.Read()
	if err != nil {
		fmt.Println("Error reading CSV headers:", err) // Add this line
		return nil, err
	}

	var data []map[string]interface{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading CSV record:", err) // Add this line
			return nil, err
		}

		row := make(map[string]interface{})
		for i, value := range record {
			row[headers[i]] = value
		}
		data = append(data, row)
	}

	return data, nil
}

func extractFromJSON(path string) ([]map[string]interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var data []map[string]interface{}
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func extractFromTXT(path string) ([]map[string]interface{}, error) {
	// Implement TXT extraction logic here
	return nil, nil
}
