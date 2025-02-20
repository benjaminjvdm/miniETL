package main

import (
	"fmt"
	"strconv"
)

// ConvertTransform represents a type conversion transformation
type ConvertTransform struct {
	fields map[string]string
}

// NewConvertTransform creates a new type conversion transformation
func NewConvertTransform(fields map[string]interface{}) (*ConvertTransform, error) {
	convertFields := make(map[string]string)
	for field, typeStr := range fields {
		typeStrVal, ok := typeStr.(string)
		if !ok {
			return nil, fmt.Errorf("type for field %s is not a string", field)
		}
		convertFields[field] = typeStrVal
	}
	return &ConvertTransform{fields: convertFields}, nil
}

// Apply applies the type conversion transformation to the data
func (t *ConvertTransform) Apply(recordChan <-chan map[string]interface{}) (<-chan map[string]interface{}, <-chan error) {
	outChan := make(chan map[string]interface{})
	errChan := make(chan error)

	go func() {
		defer close(outChan)
		defer close(errChan)

		for record := range recordChan {
			newRecord := make(map[string]interface{})
			for key, value := range record {
				newRecord[key] = value
			}

			for field, typeStr := range t.fields {
				if value, ok := newRecord[field]; ok {
					switch typeStr {
					case "int":
						strVal, ok := value.(string)
						if !ok {
							errChan <- fmt.Errorf("field %s is not a string", field)
							continue
						}
						intVal, err := strconv.Atoi(strVal)
						if err != nil {
							errChan <- fmt.Errorf("failed to convert field %s to int: %w", field, err)
							continue
						}
						newRecord[field] = intVal
					case "float":
						strVal, ok := value.(string)
						if !ok {
							errChan <- fmt.Errorf("field %s is not a string", field)
							continue
						}
						floatVal, err := strconv.ParseFloat(strVal, 64)
						if err != nil {
							errChan <- fmt.Errorf("failed to convert field %s to float: %w", field, err)
							continue
						}
						newRecord[field] = floatVal
					case "string":
						newRecord[field] = fmt.Sprintf("%v", value)
					default:
						errChan <- fmt.Errorf("unsupported type: %s", typeStr)
					}
				}
			}
			outChan <- newRecord
		}
	}()

	return outChan, errChan
}
