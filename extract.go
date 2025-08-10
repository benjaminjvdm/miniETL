package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
)

func Extract(config Config) ([]map[string]interface{}, error) {
	switch config.Input.Type {
	case "csv":
		return extractFromCSV(config.Input.Path)
	case "json":
		return extractFromJSON(config.Input.Path)
	case "txt":
		return extractFromTXT(config.Input.Path)
	default:
		return nil, ErrInvalidInputType
	}
}

func extractFromCSV(path string) ([]map[string]interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	var data []map[string]interface{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
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
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var data []map[string]interface{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var row map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &row); err != nil {
			return nil, err
		}
		data = append(data, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return data, nil
}
