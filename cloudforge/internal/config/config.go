package config

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

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
    // Define a safe base directory (e.g., current dir or a dedicated config folder)
    baseDir := "."  // Or specify e.g., "configs/" if configs are in a subdirectory

    // Sanitize the path to prevent traversal: use only the base filename
    // safePath := filepath.Join(baseDir, filepath.Base(path))

	cleanPath := filepath.Clean(path)
	// Use a rooted FS to scope access
    root := os.DirFS(baseDir)
    data, err := fs.ReadFile(root, cleanPath)
    if err != nil {
        return nil, fmt.Errorf("error reading config: %w", err)
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("error parsing config: %w", err)
    }

    return &cfg, nil
}
