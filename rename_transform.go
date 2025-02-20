package main

import "fmt"

// RenameTransform represents a rename transformation
type RenameTransform struct {
	fields map[string]string
}

// NewRenameTransform creates a new rename transformation
func NewRenameTransform(fields map[string]interface{}) (*RenameTransform, error) {
	renameFields := make(map[string]string)
	for oldName, newName := range fields {
		newNameStr, ok := newName.(string)
		if !ok {
			return nil, fmt.Errorf("new name for field %s is not a string", oldName)
		}
		renameFields[oldName] = newNameStr
	}
	return &RenameTransform{fields: renameFields}, nil
}

// Apply applies the rename transformation to the data
func (t *RenameTransform) Apply(recordChan <-chan map[string]interface{}) (<-chan map[string]interface{}, <-chan error) {
	outChan := make(chan map[string]interface{})
	errChan := make(chan error)

	go func() {
		defer close(outChan)
		defer close(errChan)

		for record := range recordChan {
			newRecord := make(map[string]interface{})
			for key, value := range record {
				newKey := key
				if newName, ok := t.fields[key]; ok {
					newKey = newName
				}
				newRecord[newKey] = value
			}
			outChan <- newRecord
		}
	}()

	return outChan, errChan
}
