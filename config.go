package main

import (
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

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
