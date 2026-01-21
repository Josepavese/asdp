package usecase

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func validateAndExpandPath(path string) (string, error) {
	// 1. Handle Empty
	if path == "" {
		return "", fmt.Errorf("path cannot be empty. Please provide an absolute path")
	}

	// 2. Expand ~/
	if strings.HasPrefix(path, "~/") || path == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		if path == "~" {
			path = home
		} else {
			path = filepath.Join(home, path[2:])
		}
	}

	// 3. Force Absolute
	if !filepath.IsAbs(path) {
		return "", fmt.Errorf("invalid path: '%s'. ASDP tools require an ABSOLUTE path or a path starting with '~/' to avoid wandering outside the project", path)
	}

	return filepath.Clean(path), nil
}
