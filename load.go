package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

func Load(config Config, data []map[string]interface{}) error {
	switch config.Output.Type {
	case "csv":
		err := loadToCSV(config.Output.Path, data)
		if err != nil {
			return err
		}
		return nil
	case "json":
		err := loadToJSON(config.Output.Path, data)
		if err != nil {
			return err
		}
		return nil
	case "txt":
		err := loadToTXT(config.Output.Path, data)
		if err != nil {
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
	// Implement TXT loading logic here
	return nil
}
