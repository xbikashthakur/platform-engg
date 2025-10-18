// platform-engg/cloudforge/cmd/cloudforge/main.go
package main

import (
	"fmt"
	"log"
	"os"

	// Full module path required for the internal package
	"github.com/xbikashthakur/platform-engg/cloudforge/internal/config"
)

func main() {
	logger := log.New(os.Stdout, "[Cloudforge] ", log.LstdFlags)

	if len(os.Args) < 2 {
		logger.Fatal("Usage: cloudforge <config-file>")
	}

	configPath := os.Args[1]
	logger.Printf("Loading config from: %s", configPath)

	cfg, err := config.Load(configPath)
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	logger.Printf("Project: %s, Environment: %s, Region: %s",
		cfg.ProjectName, cfg.Environment, cfg.Region)

	fmt.Printf("\n✓ Configuration loaded successfully\n")
	fmt.Printf("  Project: %s\n", cfg.ProjectName)
	fmt.Printf("  Environment: %s\n", cfg.Environment)
	fmt.Printf("  Region: %s\n", cfg.Region)
}
