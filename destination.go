package main

import "fmt"

// Destination represents an output destination
type Destination interface {
	Write(<-chan map[string]interface{}) <-chan error
}

// NewDestination creates a new destination based on the configuration
func NewDestination(config DestinationConfig) (Destination, error) {
	switch config.Type {
	case "csv":
		return NewCSVDestination(config.Path), nil
	default:
		return nil, fmt.Errorf("unsupported destination type: %s", config.Type)
	}
}
