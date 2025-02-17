package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Extractor interface {
	Extract(path string) ([]map[string]interface{}, error)
}

func Extract(config Config) ([]map[string]interface{}, error) {
	extractors := map[string]Extractor{
		"csv":  csvExtractor{},
		"json": jsonExtractor{},
		"txt":  txtExtractor{},
	}

	extractor, ok := extractors[config.Input.Type]
	if !ok {
		return nil, fmt.Errorf("unsupported input type: %s", config.Input.Type)
	}

	data, err := extractor.Extract(config.Input.Path)
	if err != nil {
		ErrorLogger.Println("Error extracting data:", err)
		DebugLogger.Println("Error extracting data:", err)
		return nil, err
	}

	return data, nil
}

type csvExtractor struct{}

func (e csvExtractor) Extract(path string) ([]map[string]interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	headers, err := reader.Read()
	if err != nil {
		fmt.Println("Error reading CSV headers:", err)
		return nil, err
	}

	var data []map[string]interface{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading CSV record:", err)
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

type jsonExtractor struct{}

func (e jsonExtractor) Extract(path string) ([]map[string]interface{}, error) {
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

type txtExtractor struct{}

func (e txtExtractor) Extract(path string) ([]map[string]interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read the config to get the delimiter
	config, err := LoadConfig(".")
	if err != nil {
		return nil, err
	}

	delimiter := config.Input.Delimiter
	if delimiter == "" {
		delimiter = "," // Default delimiter
	}

	reader := csv.NewReader(f)
	reader.Comma = rune(delimiter[0]) // Set the delimiter

	var data []map[string]interface{}
	headers := []string{}
	firstLine := true

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if firstLine {
			headers = record
			firstLine = false
			continue
		}

		row := make(map[string]interface{})
		for i, value := range record {
			row[headers[i]] = value
		}
		data = append(data, row)
	}

	return data, nil
}
