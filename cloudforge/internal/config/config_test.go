package config

import (
	"os"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	// --- Test Case 1: Successful Load ---
	t.Run("Successful Load", func(t *testing.T) {
		expectedConfig := &Config{
			ProjectName: "CloudForge",
			Environment: "dev",
			Region:      "us-east-1",
			Tags: map[string]string{
				"cost-center": "12345",
				"owner":       "platform-engg",
			},
		}

		content := `
project_name: CloudForge
environment: dev
region: us-east-1
tags:
  cost-center: "12345"
  owner: "platform-engg"
`
		tmpFile, err := os.CreateTemp("", "config-*.yaml")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		// Use t.Cleanup to handle the error.
		t.Cleanup(func() {
			if err := os.Remove(tmpFile.Name()); err != nil {
				t.Logf("WARN: failed to remove temp file %s: %v", tmpFile.Name(), err)
			}
		})

		if _, err := tmpFile.Write([]byte(content)); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		if err := tmpFile.Close(); err != nil {
			t.Fatalf("Failed to close temp file: %v", err)
		}

		loadedConfig, err := Load(tmpFile.Name())

		if err != nil {
			t.Errorf("Load() returned an unexpected error: %v", err)
		}
		if loadedConfig == nil {
			t.Fatal("Load() returned a nil config, but expected a valid one.")
		}
		if !reflect.DeepEqual(expectedConfig, loadedConfig) {
			t.Errorf("Loaded config does not match expected.\nGot: %+v\nWant: %+v", loadedConfig, expectedConfig)
		}
	})

	// --- Test Case 2: File Does Not Exist ---
	t.Run("File Not Found", func(t *testing.T) {
		loadedConfig, err := Load("non-existent-file.yaml")

		if err == nil {
			t.Error("Load() should have returned an error for a non-existent file, but it did not.")
		}
		if loadedConfig != nil {
			t.Errorf("Load() should have returned a nil config on error, but it did not.")
		}
	})

	// --- Test Case 3: Invalid YAML Content ---
	t.Run("Invalid YAML", func(t *testing.T) {
		content := `
project_name: "Some Project"
  invalid_indent: true
`
		tmpFile, err := os.CreateTemp("", "config-*.yaml")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		// Use t.Cleanup to handle the error.
		t.Cleanup(func() {
			if err := os.Remove(tmpFile.Name()); err != nil {
				t.Logf("WARN: failed to remove temp file %s: %v", tmpFile.Name(), err)
			}
		})

		if _, err := tmpFile.Write([]byte(content)); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		if err := tmpFile.Close(); err != nil {
			t.Fatalf("Failed to close temp file: %v", err)
		}

		loadedConfig, err := Load(tmpFile.Name())

		if err == nil {
			t.Error("Load() should have returned an error for invalid YAML, but it did not.")
		}
		if loadedConfig != nil {
			t.Errorf("Load() should have returned a nil config on error, but it did not.")
		}
	})
}
