package main

import (
	"fmt"
	"strconv"
	"strings"
)

// CalculateTransform represents a field calculation transformation
type CalculateTransform struct {
	field      string
	expression string
}

// NewCalculateTransform creates a new field calculation transformation
func NewCalculateTransform(fields map[string]interface{}) (*CalculateTransform, error) {
	field, ok := fields["field"].(string)
	if !ok {
		return nil, fmt.Errorf("field is not a string")
	}

	expression, ok := fields["expression"].(string)
	if !ok {
		return nil, fmt.Errorf("expression is not a string")
	}

	return &CalculateTransform{
		field:      field,
		expression: expression,
	}, nil
}

// Apply applies the field calculation transformation to the data
func (d *CalculateTransform) Apply(recordChan <-chan map[string]interface{}) (<-chan map[string]interface{}, <-chan error) {
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

			result, err := evaluateExpression(d.expression, newRecord)
			if err != nil {
				errChan <- fmt.Errorf("failed to evaluate expression: %w", err)
				continue
			}

			newRecord[d.field] = result
			outChan <- newRecord
		}
	}()

	return outChan, errChan
}

func evaluateExpression(expression string, record map[string]interface{}) (interface{}, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	parts := strings.Split(expression, "*")

	if len(parts) != 2 {
		return nil, fmt.Errorf("unsupported expression format")
	}

	field1 := parts[0]
	field2Str := parts[1]

	value1, ok := record[field1].(int)
	if !ok {
		return nil, fmt.Errorf("field %s is not an integer", field1)
	}

	value2, err := strconv.ParseFloat(field2Str, 64)
	if err != nil {
		return nil, fmt.Errorf("field %s is not a valid number", field2Str)
	}

	return float64(value1) * value2, nil
}
