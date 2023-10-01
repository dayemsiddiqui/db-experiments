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

	trafficValidationError := ValidateAndNormalizeTrafficConfig(&cfg)
	if trafficValidationError != nil {
		return nil, trafficValidationError
	}
	return &cfg, nil
}

func ValidateAndNormalizeTrafficConfig(cfg *Config) error {
	totalPercent := 0
	queriesWithoutTraffic := 0

	// Calculate total traffic percent and count queries without specified traffic
	for _, q := range cfg.Queries {
		totalPercent += q.TrafficPercent
		if q.TrafficPercent == 0 {
			queriesWithoutTraffic++
		}
	}

	// If totalPercent is 100, no need to normalize
	if totalPercent == 100 {
		return nil
	}

	// If totalPercent is not 100 and there are queries without specified traffic, distribute remaining traffic among them
	if queriesWithoutTraffic > 0 {
		remainingTraffic := 100 - totalPercent
		trafficPerQuery := remainingTraffic / queriesWithoutTraffic

		for i := range cfg.Queries {
			if cfg.Queries[i].TrafficPercent == 0 {
				cfg.Queries[i].TrafficPercent = trafficPerQuery
			}
		}

		return nil
	}

	return errors.New("total traffic_percent does not sum up to 100%")
}
