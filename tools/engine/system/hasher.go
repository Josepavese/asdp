package system

import (
	"crypto/sha256"
	"encoding/hex"
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

// HashDir calculates a deterministic hash of the semantic content of a directory.
// It ignores dependencies, build artifacts, environments, and hidden files.
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

	// 1. Walk and collect relevant files
	err := h.fs.Walk(root, func(path string, isDir bool) error {
		name := filepath.Base(path)

		if isDir {
			// Skip ignored directories entirely for performance
			if isIgnored(name) || (strings.HasPrefix(name, ".") && path != root) {
				return filepath.SkipDir
			}
			return nil
		}

		// Filtering rules for files:
		// - Skip if parent/name is ignored
		// - Skip markdown (docs shouldn't invalidate code hash)
		// - Skip tests (optional, but good for "Interface Stability")
		if isIgnored(name) ||
			strings.HasPrefix(name, ".") ||
			strings.HasSuffix(name, ".md") ||
			strings.HasSuffix(name, "_test.go") {
			return nil
		}

		// Extra safety: Check if it's a regular file or a symlink to a file
		info, err := os.Stat(path)
		if err != nil {
			return nil // Skip files we can't stat
		}
		if !info.Mode().IsRegular() {
			return nil // Skip directories (if missed by choice), symlinks to dirs, devices, etc.
		}

		files = append(files, path)
		return nil
	})
	if err != nil {
		return "", err
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
