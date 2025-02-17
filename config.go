package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Input struct {
		Path string `mapstructure:"path"`
		Type string `mapstructure:"type"`
	} `mapstructure:"input"`
	Output struct {
		Path string `mapstructure:"path"`
		Type string `mapstructure:"type"`
	} `mapstructure:"output"`
	Transformations []map[string]interface{} `mapstructure:"transformations"`
}

func LoadConfig(path string) (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	var config Config

	// Get the absolute path to the config file
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("Error getting absolute path: %s\n", err)
		return config, err
	}

	// Add the directory containing the config file to the search path
	viper.AddConfigPath(filepath.Dir(absPath))

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("Unable to decode into struct: %v\n", err)
		return config, err
	}

	fmt.Printf("Config values: %+v\n", config) // Print config values

	return config, nil
}
