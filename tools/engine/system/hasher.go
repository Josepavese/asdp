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
// It ignores hidden files (dotfiles), test files (_test.go), and docs (.md).
func (h *SHA256ContentHasher) HashDir(root string) (string, error) {
	var files []string

	// 1. Walk and collect relevant files
	err := h.fs.Walk(root, func(path string, isDir bool) error {
		if isDir {
			if strings.HasPrefix(filepath.Base(path), ".") && path != root {
				return filepath.SkipDir // Skip hidden dirs like .git
			}
			return nil
		}

		name := filepath.Base(path)
		// Filtering rules:
		// - Skip hidden files
		// - Skip markdown (docs shouldn't invalidate code hash)
		// - Skip tests (optional, but good for "Interface Stability")
		if strings.HasPrefix(name, ".") ||
			strings.HasSuffix(name, ".md") ||
			strings.HasSuffix(name, "_test.go") {
			return nil
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
			return "", err
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
