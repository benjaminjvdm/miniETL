package main

import (
	"fmt"
)

type TransformationFunc func(map[string]interface{}) (map[string]interface{}, error)

func Transform(data []map[string]interface{}, transformations []map[string]interface{}) ([]map[string]interface{}, error) {
	for _, row := range data {
		fmt.Println(row)
	}
	return data, nil
}
