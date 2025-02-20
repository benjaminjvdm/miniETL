package main

import "fmt"

// Transform represents a data transformation
type Transform interface {
	Apply(<-chan map[string]interface{}) (<-chan map[string]interface{}, <-chan error)
}

// NewTransform creates a new transformation based on the configuration
func NewTransform(config TransformConfig) (Transform, error) {
	switch config.Type {
	case "rename":
		transform, err := NewRenameTransform(config.Fields)
		if err != nil {
			return nil, fmt.Errorf("failed to create rename transform: %w", err)
		}
		return transform, nil
	case "convert":
		transform, err := NewConvertTransform(config.Fields)
		if err != nil {
			return nil, fmt.Errorf("failed to create convert transform: %w", err)
		}
		return transform, nil
	case "filter":
		transform, err := NewFilterTransform(config.Fields)
		if err != nil {
			return nil, fmt.Errorf("failed to create filter transform: %w", err)
		}
		return transform, nil
	case "calculate":
		transform, err := NewCalculateTransform(config.Fields)
		if err != nil {
			return nil, fmt.Errorf("failed to create calculate transform: %w", err)
		}
		return transform, nil
	default:
		return nil, fmt.Errorf("unsupported transform type: %s", config.Type)
	}
}
