package main

import (
	"fmt"
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
	viper.SetConfigFile(path)

	var config Config

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	// Basic validation
	if config.Input.Path == "" {
		return config, fmt.Errorf("validation failed: Input.Path is required")
	}
	if config.Input.Type == "" {
		return config, fmt.Errorf("validation failed: Input.Type is required")
	}
	if config.Output.Path == "" {
		return config, fmt.Errorf("validation failed: Output.Path is required")
	}
	if config.Output.Type == "" {
		return config, fmt.Errorf("validation failed: Output.Type is required")
	}

	return config, nil
}
