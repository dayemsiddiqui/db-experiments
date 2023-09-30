package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"math"
	"os"
)

type QueryConfig struct {
	Name           string  `yaml:"name"`
	Query          string  `yaml:"query"`
	TrafficPercent float64 `yaml:"traffic_percent"`
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

func roundToPlaces(value float64, places int) float64 {
	multiplier := math.Pow(10, float64(places))
	return math.Round(value*multiplier) / multiplier
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

	totalPercent := 0.0
	for _, q := range cfg.Queries {
		totalPercent += q.TrafficPercent
	}
	if roundToPlaces(totalPercent, 2) != 1.0 {
		return nil, errors.New("total traffic_percent does not sum up to 100%")
	}
	return &cfg, nil
}
