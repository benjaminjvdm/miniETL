package main

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/pflag"
)

func runETL(configFile string) error {
	// Load configuration
	InfoLogger.Println("Loading configuration from", configFile)
	config, err := LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}
	InfoLogger.Println("Configuration loaded successfully")

	bar := progressbar.Default(3, "ETL Process")

	// Extract
	InfoLogger.Println("Extracting data...")
	data, err := Extract(config)
	if err != nil {
		return fmt.Errorf("error extracting data: %w", err)
	}
	InfoLogger.Printf("Extracted %d records", len(data))
	bar.Add(1)

	// Transform
	InfoLogger.Println("Transforming data...")
	data, err = Transform(data, config.Transformations)
	if err != nil {
		return fmt.Errorf("error transforming data: %w", err)
	}
	InfoLogger.Println("Data transformed successfully")
	bar.Add(1)

	// Load
	InfoLogger.Println("Loading data...")
	err = Load(config, data)
	if err != nil {
		return fmt.Errorf("error loading data: %w", err)
	}
	InfoLogger.Println("Data loaded successfully")
	bar.Add(1)

	InfoLogger.Println("ETL process completed successfully.")
	return nil
}

func main() {
	// Command-line flags
	var configFile string
	pflag.StringVar(&configFile, "config", "./config.yaml", "Path to the config file")
	pflag.Parse()

	if err := runETL(configFile); err != nil {
		ErrorLogger.Fatalf(err.Error())
	}
}
