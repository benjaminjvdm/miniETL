package main

import (
	"fmt"
	"os"

	"time"

	"github.com/schollz/progressbar/v3"
)

func main() {
	// Load configuration
	config, err := LoadConfig("./config.yaml")
	if err != nil {
		fmt.Printf("Error loading config: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Config:", config) // Print the config

	// Extract
	data, err := Extract(config)
	if err != nil {
		fmt.Printf("Error extracting data: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Data:", data) // Print the data

	// Transform
	data, err = Transform(data, config.Transformations)
	if err != nil {
		fmt.Printf("Error transforming data: %s\n", err)
		os.Exit(1)
	}

	// Load
	err = Load(config, data)
	if err != nil {
		fmt.Printf("Error loading data: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("ETL process completed.")

	bar := progressbar.Default(100, "Processing")
	for i := 0; i < 100; i++ {
		time.Sleep(10 * time.Millisecond)
		bar.Add(1)
	}
}
