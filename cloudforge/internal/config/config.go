package config

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ProjectName string            `yaml:"project_name"`
	Environment string            `yaml:"environment"`
	Region      string            `yaml:"region"`
	Tags        map[string]string `yaml:"tags"`
}

// Load reads a YAML file from the given path and returns a Config struct.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	return &cfg, nil
}
