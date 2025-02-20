package main

import (
	"fmt"
	"strconv"
)

// FilterTransform represents a filter transformation
type FilterTransform struct {
	field     string
	condition string
	value     interface{}
}

// NewFilterTransform creates a new filter transformation
func NewFilterTransform(fields map[string]interface{}) (*FilterTransform, error) {
	field, ok := fields["field"].(string)
	if !ok {
		return nil, fmt.Errorf("field is not a string")
	}

	condition, ok := fields["condition"].(string)
	if !ok {
		return nil, fmt.Errorf("condition is not a string")
	}

	value := fields["value"]

	return &FilterTransform{
		field:     field,
		condition: condition,
		value:     value,
	}, nil
}

// Apply applies the filter transformation to the data
func (t *FilterTransform) Apply(recordChan <-chan map[string]interface{}) (<-chan map[string]interface{}, <-chan error) {
	outChan := make(chan map[string]interface{})
	errChan := make(chan error)

	go func() {
		defer close(outChan)
		defer close(errChan)

		for record := range recordChan {
			value, ok := record[t.field]
			if !ok {
				continue // Skip records where the field is missing
			}

			match := false
			switch t.condition {
			case "equals":
				match = value == t.value
			case "greater_than":
				numValue, ok := value.(int)
				if !ok {
					errChan <- fmt.Errorf("field %s is not an integer", t.field)
					continue
				}
				var numFilterValue int
				numFilterValueIface, ok := t.value.(int)
				if ok {
					numFilterValue = numFilterValueIface
				} else {
					strFilterValue, ok := t.value.(string)
					if !ok {
						errChan <- fmt.Errorf("filter value is not an integer or string")
						continue
					}
					var err error
					numFilterValue, err = strconv.Atoi(strFilterValue)
					if err != nil {
						errChan <- fmt.Errorf("filter value is not a valid integer: %w", err)
						continue
					}

				}
				match = numValue > numFilterValue
			default:
				errChan <- fmt.Errorf("unsupported condition: %s", t.condition)
				continue
			}

			if match {
				outChan <- record
			}
		}
	}()

	return outChan, errChan
}
