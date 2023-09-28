// config.go
package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	DumpPath string   `yaml:"dump_path"`
	Queries  []string `yaml:"queries"`
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

	return &cfg, nil
}
