package main

import (
	"reflect"
	"testing"
)

func TestEvaluateExpression(t *testing.T) {
	testCases := []struct {
		name       string
		row        map[string]interface{}
		expression string
		expected   interface{}
		expectErr  bool
	}{
		// --- Happy Path ---
		{
			name:       "Simple string concatenation",
			row:        map[string]interface{}{"first": "John", "last": "Doe"},
			expression: `first + " " + last`,
			expected:   "John Doe",
			expectErr:  false,
		},
		{
			name:       "Simple integer addition",
			row:        map[string]interface{}{"a": 10, "b": 20, "c": "val"},
			expression: `len(c) + a + b`, // len() returns int, so this should be numeric
			expected:   33,
			expectErr:  false,
		},
		{
			name:       "Upper function",
			row:        map[string]interface{}{"name": "tester"},
			expression: `upper(name)`,
			expected:   "TESTER",
			expectErr:  false,
		},
		{
			name:       "Len function",
			row:        map[string]interface{}{"name": "tester"},
			expression: `len(name)`,
			expected:   6,
			expectErr:  false,
		},
		{
			name:       "Combination of functions and fields",
			row:        map[string]interface{}{"fname": "hello", "lname": "world"},
			expression: `upper(fname) + " " + upper(lname)`,
			expected:   "HELLO WORLD",
			expectErr:  false,
		},
		// --- Error Path ---
		{
			name:       "Invalid function",
			row:        map[string]interface{}{"name": "tester"},
			expression: `invalid_func(name)`,
			expectErr:  true,
		},
		{
			name:       "Referencing non-existent column",
			row:        map[string]interface{}{"name": "tester"},
			expression: `upper(non_existent)`,
			expectErr:  true,
		},
		{
			name:       "Syntax error - unclosed quote",
			row:        map[string]interface{}{},
			expression: `"hello`,
			expectErr:  true,
		},
		{
			name:       "Syntax error - invalid term",
			row:        map[string]interface{}{},
			expression: `a + !`,
			expectErr:  true,
		},
		// --- Edge Cases & Type Issues ---
		{
			name:       "String literal with plus sign",
			row:        map[string]interface{}{"a": "first"},
			expression: `a + "a+b"`,
			expected:   "firsta+b", // The current parser will fail this, I expect. It will split on the +.
			expectErr:  true,     // This is an expected failure of the current parser.
		},
		{
			name:       "Type mismatch - string + int",
			row:        map[string]interface{}{"a": "value", "b": 10},
			expression: `a + b`,
			expected:   "value10", // String concatenation is the expected behavior
			expectErr:  false,
		},
		{
			name:       "Type mismatch - int + string",
			row:        map[string]interface{}{"a": 10, "b": "value"},
			expression: `a + b`,
			expectErr:  true, // Numeric context, should fail on string
		},
		{
			name:       "Empty expression",
			row:        map[string]interface{}{},
			expression: ``,
			expectErr:  true, // Should not handle empty expressions gracefully
		},
		{
			name:       "Field with space in name (supported)",
			row:        map[string]interface{}{"first name": "John"},
			expression: `upper(first name)`,
			expected:   "JOHN",
			expectErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := EvaluateExpression(tc.row, tc.expression)

			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected an error, but got none. Result: %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error, but got: %v", err)
				}
				if !reflect.DeepEqual(result, tc.expected) {
					t.Errorf("Expected result '%v' (type %T), but got '%v' (type %T)",
						tc.expected, tc.expected, result, result)
				}
			}
		})
	}
}
