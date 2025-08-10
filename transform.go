package main

import "fmt"

type TransformationFunc func(map[string]interface{}) (map[string]interface{}, error)

func Transform(data []map[string]interface{}, transformations []map[string]interface{}) ([]map[string]interface{}, error) {
	for _, row := range data {
		for _, transform := range transformations {
			for newField, expression := range transform {
				output, err := EvaluateExpression(row, expression.(string))
				if err != nil {
					return nil, fmt.Errorf("failed to evaluate expression for field %s: %w", newField, err)
				}
				row[newField] = output
			}
		}
	}
	return data, nil
}
