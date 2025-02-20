package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config represents the ETL configuration
type Config struct {
	Source      SourceConfig      `yaml:"source"`
	Transform   []TransformConfig `yaml:"transform"`
	Destination DestinationConfig `yaml:"destination"`
}

// SourceConfig represents the configuration for the data source
type SourceConfig struct {
	Type string `yaml:"type"`
	Path string `yaml:"path"`
}

// TransformConfig represents the configuration for a transformation
type TransformConfig struct {
	Type   string                 `yaml:"type"`
	Fields map[string]interface{} `yaml:"fields"`
}

// DestinationConfig represents the configuration for the output destination
type DestinationConfig struct {
	Type string `yaml:"type"`
	Path string `yaml:"path"`
}

// LoadConfig loads the configuration from the specified file
func LoadConfig(filename string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}

	err = validateConfig(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func validateConfig(config *Config) error {
	if config.Source.Type == "" {
		return fmt.Errorf("source type is required")
	}
	if config.Source.Path == "" {
		return fmt.Errorf("source path is required")
	}
	if config.Destination.Type == "" {
		return fmt.Errorf("destination type is required")
	}
	if config.Destination.Path == "" {
		return fmt.Errorf("destination path is required")
	}
	return nil
}
