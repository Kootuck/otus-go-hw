package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger  LoggerConf  `yaml:"logger"`
	Storage StorageConf `yaml:"storage"`
	HTTP    HTTPConf    `yaml:"http"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type StorageConf struct {
	Method string       `yaml:"method"`
	PG     PGConnection `yaml:"pg"`
}

type HTTPConf struct {
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	ShutdownTimeout int    `yaml:"shutdownTimeout"`
}

type PGConnection struct {
	DBUser     string `yaml:"dbUser"`
	DBPassword string `yaml:"dbPassword"`
	DBName     string `yaml:"dbName"`
	DBHost     string `yaml:"dbHost"`
	DBPort     string `yaml:"dbPort"`
	Timeout    int    `yaml:"execDefaultTimeoutSeconds"`
}

func NewConfig(path string) (Config, error) {
	var config Config

	// Read the config file
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal the YAML into our struct
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return config, fmt.Errorf("error parsing config file: %w", err)
	}

	return config, nil
}
