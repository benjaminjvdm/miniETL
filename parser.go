package main

import (
	"fmt"
	"strings"
)

// EvaluateExpression evaluates a simple expression against a data row.
func EvaluateExpression(row map[string]interface{}, expression string) (interface{}, error) {
	parts := strings.Split(expression, "+")
	var result interface{}
	var isNumeric bool

	for i, part := range parts {
		part = strings.TrimSpace(part)
		value, err := evaluateTerm(row, part)
		if err != nil {
			return nil, err
		}

		if i == 0 {
			result = value
			_, isNumeric = result.(int)
		} else {
			if isNumeric {
				// All parts must be numeric to add
				currentVal, ok := value.(int)
				if !ok {
					return nil, fmt.Errorf("type mismatch: cannot add non-numeric value to a number")
				}
				result = result.(int) + currentVal
			} else {
				// All parts must be strings to concatenate
				result = fmt.Sprintf("%v%v", result, value)
			}
		}
	}
	return result, nil
}

// evaluateTerm evaluates a single part of an expression.
func evaluateTerm(row map[string]interface{}, term string) (interface{}, error) {
	// 1. String literal
	if strings.HasPrefix(term, "\"") && strings.HasSuffix(term, "\"") {
		return term[1 : len(term)-1], nil
	}

	// 2. upper(field)
	if strings.HasPrefix(term, "upper(") && strings.HasSuffix(term, ")") {
		field := strings.TrimSpace(term[6 : len(term)-1])
		if val, ok := row[field]; ok {
			if strVal, ok := val.(string); ok {
				return strings.ToUpper(strVal), nil
			}
		}
		return nil, fmt.Errorf("invalid field for upper(): %s", field)
	}

	// 3. len(field)
	if strings.HasPrefix(term, "len(") && strings.HasSuffix(term, ")") {
		field := strings.TrimSpace(term[4 : len(term)-1])
		if val, ok := row[field]; ok {
			if strVal, ok := val.(string); ok {
				return len(strVal), nil
			}
		}
		return nil, fmt.Errorf("invalid field for len(): %s", field)
	}

	// 4. Field name
	if val, ok := row[term]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("unsupported or invalid term: %s", term)
}
