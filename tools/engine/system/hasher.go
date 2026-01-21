package system

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type SHA256ContentHasher struct {
	fs *RealFileSystem
}

func NewSHA256ContentHasher() *SHA256ContentHasher {
	return &SHA256ContentHasher{fs: NewRealFileSystem()}
}

// HashDir calculates a deterministic hash of the semantic content of a directory (Non-Recursive).
// It only considers regular files in the root folder, ignoring dependencies and hidden items.
func (h *SHA256ContentHasher) HashDir(root string) (string, error) {
	var files []string

	ignoredPatterns := []string{
		"node_modules", "vendor", "packages", "bower_components", // Dependencies
		"venv", ".venv", "anaconda", "conda", "env", ".env", // Environments
		"bin", "obj", "dist", "target", "build", "out", // Builds
		".vscode", ".idea", ".git", ".hg", ".svn", ".cache", // IDE/System
	}

	isIgnored := func(name string) bool {
		name = strings.ToLower(name)
		for _, p := range ignoredPatterns {
			if strings.Contains(name, p) {
				return true
			}
		}
		return false
	}

	// 1. Read just the root directory (NON-RECURSIVE)
	entries, err := os.ReadDir(root)
	if err != nil {
		return "", fmt.Errorf("failed to read directory %s: %w", root, err)
	}

	for _, entry := range entries {
		name := entry.Name()
		path := filepath.Join(root, name)

		if entry.IsDir() {
			continue // Skip subdirectories
		}

		// Filtering rules for files:
		if isIgnored(name) ||
			strings.HasPrefix(name, ".") ||
			strings.HasSuffix(name, ".md") ||
			strings.HasSuffix(name, "_test.go") {
			continue
		}

		// Extra safety: Check if it's a regular file (following symlinks if any)
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		if !info.Mode().IsRegular() {
			continue
		}

		files = append(files, path)
	}

	// 2. Sort files to ensure deterministic order
	sort.Strings(files)

	// 3. Hash content
	hasher := sha256.New()
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			continue // Skip files we can't open after collection
		}

		// Hash filename first (to detect renames)
		relPath, _ := filepath.Rel(root, file)
		hasher.Write([]byte(relPath))

		// Hash content
		if _, err := io.Copy(hasher, f); err != nil {
			f.Close()
			return "", err
		}
		f.Close()
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
