package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Keywords []string `yaml:"keywords"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile("sus.yaml")
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
