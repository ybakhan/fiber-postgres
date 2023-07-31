package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents configurable properties of tax-calculator
type Config struct {
	Port    uint   `yaml:"port"`
	Version string `yaml:"version"`

	Database struct {
		Host     string `yaml:"host"`
		Port     uint   `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		SSLMode  string `yaml:"sslMode"`
	} `yaml:"database"`

	Cache struct {
		Host string `yaml:"host"`
		Port uint   `yaml:"port"`
	} `yaml:"cache"`
}

func ReadConfig() *Config {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	return &config
}
