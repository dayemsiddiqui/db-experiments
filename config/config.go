package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"os"
)

type QueryConfig struct {
	Name           string `yaml:"name"`
	Query          string `yaml:"query"`
	TrafficPercent int    `yaml:"traffic_percent"`
}

type InputParameters struct {
	Name  string   `yaml:"name"`
	Value []string `yaml:"value"`
}

type Config struct {
	DumpPath   string            `yaml:"dump_path"`
	Traffic    int               `yaml:"traffic"`
	Queries    []QueryConfig     `yaml:"queries"`
	Parameters []InputParameters `yaml:"parameters"`
}

func ReadConfig() (*Config, error) {
	configFile, err := os.ReadFile("experiments.yaml")
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		return nil, err
	}

	trafficValidationError := ValidateTrafficConfig(cfg)
	if trafficValidationError != nil {
		return nil, trafficValidationError
	}
	return &cfg, nil
}

func ValidateTrafficConfig(cfg Config) error {
	totalPercent := 0
	for _, q := range cfg.Queries {
		totalPercent += q.TrafficPercent
	}
	if totalPercent != 100 {
		return errors.New("total traffic_percent does not sum up to 100%")
	}
	return nil
}
