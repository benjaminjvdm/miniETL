package main

import (
	"fmt"
	"strings"
	"time"
)

type TransformationFunc func(map[string]interface{}) (map[string]interface{}, error)

func Transform(data []map[string]interface{}, transformations []map[string]interface{}) ([]map[string]interface{}, error) {
	for _, transformation := range transformations {
		field, ok := transformation["field"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid transformation: missing field")
		}
		action, ok := transformation["action"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid transformation: missing action")
		}

		for _, row := range data {
			switch action {
			case "uppercase":
				if str, ok := row[field].(string); ok {
					row[field] = toUppercase(str)
				}
			case "lowercase":
				if str, ok := row[field].(string); ok {
					row[field] = strings.ToLower(str)
				}
			case "trim":
				if str, ok := row[field].(string); ok {
					row[field] = strings.TrimSpace(str)
				}
			case "date_format":
				if str, ok := row[field].(string); ok {
					params := transformation["params"].(map[string]interface{})
					format, ok := params["format"].(string)
					if !ok {
						return nil, fmt.Errorf("invalid transformation: missing format parameter")
					}
					row[field] = formatDate(str, format)
				} else {
					ErrorLogger.Println("Error: Field is not a string")
				}
			case "custom":
				if functionName, ok := transformation["params"].(map[string]interface{})["function"].(string); ok {
					// Call the custom function
					// Assuming custom functions are defined elsewhere and accessible
					// This is a placeholder, actual implementation will depend on how custom functions are registered and called
					fmt.Println("Calling custom function:", functionName)
					// row[field] = callCustomFunction(functionName, row[field])
				}
			}
		}
	}
	return data, nil
}

func formatDate(dateStr string, format string) string {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return "INVALID DATE"
	}
	return t.Format(format)
}

func toUppercase(s string) string {
	return strings.ToUpper(s)
}
