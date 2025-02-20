package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

// CSVDestination represents a CSV output destination
type CSVDestination struct {
	path string
}

// NewCSVDestination creates a new CSV output destination
func NewCSVDestination(path string) *CSVDestination {
	return &CSVDestination{path: path}
}

// Write writes data to the CSV file
func (d *CSVDestination) Write(recordChan <-chan map[string]interface{}) <-chan error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)

		file, err := os.Create(d.path)
		if err != nil {
			errChan <- fmt.Errorf("failed to create CSV file: %w", err)
			return
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		var records []map[string]interface{}
		for record := range recordChan {
			records = append(records, record)
		}

		if len(records) == 0 {
			return
		}

		var header []string
		for key := range records[0] {
			header = append(header, key)
		}

		err = writer.Write(header)
		if err != nil {
			errChan <- fmt.Errorf("failed to write CSV header: %w", err)
			return
		}

		for _, recordMap := range records {
			record := make([]string, len(header))
			for i, key := range header {
				record[i] = fmt.Sprintf("%v", recordMap[key])
			}
			err = writer.Write(record)
			if err != nil {
				errChan <- fmt.Errorf("failed to write CSV record: %w", err)
				return
			}
		}
	}()

	return errChan
}
