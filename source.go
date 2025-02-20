package main

import "fmt"

// Source represents a data source
type Source interface {
	Read() (<-chan map[string]interface{}, <-chan error)
}

// NewSource creates a new data source based on the configuration
func NewSource(config SourceConfig) (Source, error) {
	switch config.Type {
	case "csv":
		return NewCSVSource(config.Path), nil
	default:
		return nil, fmt.Errorf("unsupported source type: %s", config.Type)
	}
}
