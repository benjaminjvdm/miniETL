package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

// CSVSource represents a CSV data source
type CSVSource struct {
	path string
}

// NewCSVSource creates a new CSV data source
func NewCSVSource(path string) *CSVSource {
	return &CSVSource{path: path}
}

// Read reads data from the CSV file
func (s *CSVSource) Read() (<-chan map[string]interface{}, <-chan error) {
	recordChan := make(chan map[string]interface{})
	errChan := make(chan error)

	go func() {
		defer close(recordChan)
		defer close(errChan)

		file, err := os.Open(s.path)
		if err != nil {
			errChan <- fmt.Errorf("failed to open CSV file: %w", err)
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		header, err := reader.Read()
		if err != nil {
			errChan <- fmt.Errorf("failed to read CSV header: %w", err)
			return
		}

		for {
			record, err := reader.Read()
			if err != nil {
				if err.Error() == "EOF" {
					return
				}
				errChan <- fmt.Errorf("failed to read CSV record: %w", err)
				return
			}

			recordMap := make(map[string]interface{})
			for i, value := range record {
				recordMap[header[i]] = value
			}

			recordChan <- recordMap
		}
	}()

	return recordChan, errChan
}
