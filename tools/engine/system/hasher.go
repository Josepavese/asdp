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
	fs          *RealFileSystem
	ignoreDirs  []string
	ignoreFiles []string
}

func NewSHA256ContentHasher(ignoreDirs, ignoreFiles []string) *SHA256ContentHasher {
	return &SHA256ContentHasher{
		fs:          NewRealFileSystem(),
		ignoreDirs:  ignoreDirs,
		ignoreFiles: ignoreFiles,
	}
}

// HashDir calculates a deterministic hash of the semantic content of a directory (Non-Recursive).
// It only considers regular files in the root folder, ignoring dependencies and hidden items.
func (h *SHA256ContentHasher) HashDir(root string) (string, error) {
	var files []string

	isIgnored := func(name string) bool {
		name = strings.ToLower(name)
		for _, p := range h.ignoreDirs {
			if strings.Contains(name, p) {
				return true
			}
		}
		return false
	}

	isIgnoredFile := func(name string) bool {
		name = strings.ToLower(name)
		for _, p := range h.ignoreFiles {
			// suffix matching for extensions, exact/pattern for others?
			// Simplification: if pattern starts with *, suffix match.
			if strings.HasPrefix(p, "*") {
				if strings.HasSuffix(name, p[1:]) {
					return true
				}
			} else if name == p {
				return true
			} else if strings.Contains(name, p) {
				// Fallback to contains for safety if unsure of format
				return true
			}
		}
		return false
	}

	// 1. Walk directory RECURSIVELY (Boundary-Aware)
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		name := d.Name()
		if d.IsDir() {
			if path == root {
				return nil
			}
			// Skip ignored directories
			if isIgnored(name) {
				return filepath.SkipDir
			}
			// Boundary Check: If this directory is a separate ASDP module, skip it.
			// We check for codespec.md or codemodel.md
			if _, err := os.Stat(filepath.Join(path, "codespec.md")); err == nil {
				return filepath.SkipDir
			}
			if _, err := os.Stat(filepath.Join(path, "codemodel.md")); err == nil {
				return filepath.SkipDir
			}
			return nil
		}

		// Filtering rules for files:
		// Filtering rules for files:
		if isIgnored(name) || isIgnoredFile(name) || strings.HasPrefix(name, ".") {
			return nil
		}

		// Extra safety check for regular files
		info, err := d.Info()
		if err != nil || !info.Mode().IsRegular() {
			return nil
		}

		files = append(files, path)
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to walk directory %s: %w", root, err)
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
