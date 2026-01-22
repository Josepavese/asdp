package usecase

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Josepavese/asdp/engine/domain"
)

type InitAgentUseCase struct {
	fs     domain.FileSystem
	config domain.Config
}

func NewInitAgentUseCase(fs domain.FileSystem, config domain.Config) *InitAgentUseCase {
	return &InitAgentUseCase{fs: fs, config: config}
}

func (uc *InitAgentUseCase) Execute(projectPath string) (string, error) {
	absPath, err := validateAndExpandPath(projectPath)
	if err != nil {
		return "", err
	}
	projectPath = absPath

	// 1. Resolve source directory from Config
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	srcDir := filepath.Join(home, uc.config.System.GlobalAssetsDir)

	// Verify source exists
	if _, err := uc.fs.Stat(srcDir); err != nil {
		return "", fmt.Errorf("global ASDP assets not found at %s. Please run the installer first", srcDir)
	}

	// 2. Resolve target directory from Config
	targetDir := filepath.Join(projectPath, uc.config.System.DefaultAgentDir)

	// 3. Walk and Copy
	created := 0
	err = uc.fs.Walk(srcDir, func(path string, isDir bool) error {
		// Calculate relative path
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(targetDir, rel)

		if isDir {
			return uc.fs.MkdirAll(destPath)
		}

		// Read source
		data, err := uc.fs.ReadFile(path)
		if err != nil {
			return err
		}

		// Write destination
		if err := uc.fs.WriteFile(destPath, data); err != nil {
			return err
		}
		created++
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to initialize agent assets: %w", err)
	}

	return fmt.Sprintf("Initialization complete. Created %d files in %s", created, targetDir), nil
}
