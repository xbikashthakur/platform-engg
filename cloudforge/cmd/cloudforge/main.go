package main

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

func main() {
	// Ensure the user provides a config file path as an argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: cloudforge <config-file>")
	}

	configPath := os.Args[1]
	// Read the entire file into memory
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}

	// Unmarshal (parse) the YAML data into our Config struct
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Printf("Error parsing config file: %v\n", err)
		os.Exit(1)
	}

	// Print the loaded configuration to verify it works
	fmt.Printf("CloudForge CLI\n")
	fmt.Printf("Project: %s\n", config.ProjectName)
	fmt.Printf("Environment: %s\n", config.Environment)
	fmt.Printf("Region: %s\n", config.Region)
	fmt.Printf("Tags: %v\n", config.Tags)
}
