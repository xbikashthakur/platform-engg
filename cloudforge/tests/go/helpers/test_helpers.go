package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Helper function to find the root of the Git repository by looking for the .git folder.
func FindRepoRoot() (string, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("could not get caller information")
	}
	dir := filepath.Dir(currentFile)
	maxDepth := 10 // Prevent infinite loop
	for i := 0; i < maxDepth; i++ {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			return "", fmt.Errorf("reached filesystem root without finding .git directory")
		}
		dir = parentDir
	}

	return "", fmt.Errorf("exceeded maximum depth (%d) without finding .git directory", maxDepth)
}
