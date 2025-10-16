package main

import (
	"fmt"
	"log"
	"os"
	"gopkg.in/yaml.v3"

)

type Config struct {
	ProjectName	string				`yaml:"project_name"`
	Environment	string				`yaml:"environment"`
	Region		string				`yaml:"region"`
	Tags		map[string]string	`yaml:"tags"`

}

func main() {
	logger := log.New(os.Stdout, "[Cloudforge] ", log.LstdFlags)

	if len(os.Args) < 2 {
		logger.Fatal("Usage: cloudforge <config-file>")
	}

	configPath := os.Args[1]
	logger.Printf("Loading config from: %s", configPath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		logger.Fatalf("Error reading config: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		logger.Fatalf("Error parsing config: %v", err)
	}

	logger.Printf("Project: %s, Environment: %s, Region: %s",
        config.ProjectName, config.Environment, config.Region)

	fmt.Printf("\n✓ Configuration loaded successfully\n")
    fmt.Printf("  Project: %s\n", config.ProjectName)
    fmt.Printf("  Environment: %s\n", config.Environment)
    fmt.Printf("  Region: %s\n", config.Region)
}