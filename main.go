package main

import (
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/pflag"
)

func main() {
	// Command-line flags
	var configFile string
	pflag.StringVar(&configFile, "config", "./config.yaml", "Path to the config file")
	pflag.Parse()

	// Load configuration
	InfoLogger.Println("Loading configuration from", configFile)
	config, err := LoadConfig(configFile)
	if err != nil {
		ErrorLogger.Fatalf("Error loading config: %s", err)
	}
	InfoLogger.Println("Configuration loaded successfully")

	bar := progressbar.Default(3, "ETL Process")

	// Extract
	InfoLogger.Println("Extracting data...")
	data, err := Extract(config)
	if err != nil {
		ErrorLogger.Fatalf("Error extracting data: %s", err)
	}
	InfoLogger.Printf("Extracted %d records", len(data))
	bar.Add(1)

	// Transform
	InfoLogger.Println("Transforming data...")
	data, err = Transform(data, config.Transformations)
	if err != nil {
		ErrorLogger.Fatalf("Error transforming data: %s", err)
	}
	InfoLogger.Println("Data transformed successfully")
	bar.Add(1)

	// Load
	InfoLogger.Println("Loading data...")
	err = Load(config, data)
	if err != nil {
		ErrorLogger.Fatalf("Error loading data: %s", err)
	}
	InfoLogger.Println("Data loaded successfully")
	bar.Add(1)

	InfoLogger.Println("ETL process completed successfully.")
}
