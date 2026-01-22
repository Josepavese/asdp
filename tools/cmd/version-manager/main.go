package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func performBump(cfg *VersionConfig, path string) error {
	// Parse current version
	parts := strings.Split(cfg.KitVersion, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid semantic version format: %s", cfg.KitVersion)
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("invalid patch version: %s", parts[2])
	}

	// Increment patch
	newVersion := fmt.Sprintf("%s.%s.%d", parts[0], parts[1], patch+1)
	fmt.Printf("Bumping version matching patch requirements: %s -> %s\n", cfg.KitVersion, newVersion)

	// Update Struct
	cfg.KitVersion = newVersion

	// Write back to yaml
	// We read the raw file again to replace just the version string to preserve comments
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Safe regex replacement for the YAML file
	re := regexp.MustCompile(`kit_version:\s*"?[\d\.]+"?`)
	newContent := re.ReplaceAll(content, []byte(fmt.Sprintf("kit_version: %s", newVersion)))

	if err := os.WriteFile(path, newContent, 0644); err != nil {
		return err
	}

	fmt.Println("Updated version.yaml")

	// Now sync
	return performSync(cfg, false)
}

type VersionConfig struct {
	KitVersion         string `yaml:"kit_version"`
	SchemaVersion      string `yaml:"schema_version"`
	McpProtocolVersion string `yaml:"mcp_protocol_version"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: version-manager <sync|check|bump>")
		os.Exit(1)
	}
	command := os.Args[1]

	// 1. Load Master Record
	cfg, err := loadVersionConfig("version.yaml")
	if err != nil {
		fmt.Printf("Failed to load version.yaml: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded Master Record:\n  Kit: %s\n  Schema: %s\n  MCP: %s\n",
		cfg.KitVersion, cfg.SchemaVersion, cfg.McpProtocolVersion)

	switch command {
	case "sync":
		if err := performSync(cfg, false); err != nil {
			fmt.Printf("Sync failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Sync completed successfully.")
	case "check":
		if err := performSync(cfg, true); err != nil {
			fmt.Printf("Check failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Check passed: All versions are in sync.")
	case "bump":
		if err := performBump(cfg, "version.yaml"); err != nil {
			fmt.Printf("Bump failed: %v\n", err)
			os.Exit(1)
		}
		// performBump handles the valid sync execution after updating the file
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func loadVersionConfig(path string) (*VersionConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg VersionConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func performSync(cfg *VersionConfig, dryRun bool) error {
	changesCount := 0

	// 1. Update Go Constant (tools/engine/domain/version.go)
	changed, err := updateGoVersion("tools/engine/domain/version.go", cfg.KitVersion, dryRun)
	if err != nil {
		return fmt.Errorf("failed to update domain/version.go: %w", err)
	}
	if changed {
		changesCount++
	}

	// 2. Update Go Constant (tools/mcp-server/internal/adapter/mcp/server.go) - Protocol Version
	changed, err = updateMcpProtocol("tools/mcp-server/internal/adapter/mcp/server.go", cfg.McpProtocolVersion, dryRun)
	if err != nil {
		// Non-fatal if file structure changed, but good to know
		fmt.Printf("Warning: Could not update MCP protocol version: %v\n", err)
	}
	if changed {
		changesCount++
	}

	// 3. Walk and Update Frontmatter (codespec/codemodel)
	err = filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && (d.Name() == ".git" || d.Name() == "node_modules") {
			return fs.SkipDir
		}
		if !d.IsDir() && (strings.HasSuffix(path, "codespec.md") || strings.HasSuffix(path, "codemodel.md")) {
			changed, err := updateFrontmatterVersion(path, cfg.SchemaVersion, dryRun)
			if err != nil {
				return err
			}
			if changed {
				changesCount++
			}
			return nil
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk directories: %w", err)
	}

	if dryRun && changesCount > 0 {
		return fmt.Errorf("version mismatch detected: %d files out of sync", changesCount)
	}

	return nil
}

func updateGoVersion(path string, version string, dryRun bool) (bool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	// Regex for: const Version = "..."
	re := regexp.MustCompile(`const Version = "[^"]+"`)
	newContent := re.ReplaceAll(content, []byte(fmt.Sprintf(`const Version = "%s"`, version)))

	if !bytes.Equal(content, newContent) {
		if dryRun {
			fmt.Printf("[CHECK] %s would be updated to %s\n", path, version)
			return true, nil
		}
		fmt.Printf("Updating %s -> %s\n", path, version)
		return true, os.WriteFile(path, newContent, 0644)
	}
	return false, nil
}

func updateMcpProtocol(path string, version string, dryRun bool) (bool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	// Look for ProtocolVersion: "..." in struct initialization
	re := regexp.MustCompile(`ProtocolVersion:\s*".*"`)
	newContent := re.ReplaceAll(content, []byte(fmt.Sprintf(`ProtocolVersion: "%s"`, version)))

	if !bytes.Equal(content, newContent) {
		if dryRun {
			fmt.Printf("[CHECK] ProtocolVersion in %s would be updated to %s\n", path, version)
			return true, nil
		}
		fmt.Printf("Updating ProtocolVersion in %s -> %s\n", path, version)
		return true, os.WriteFile(path, newContent, 0644)
	}
	return false, nil
}

func updateFrontmatterVersion(path string, version string, dryRun bool) (bool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	// Simple YAML Frontmatter replacement using regex is safer than parsing/marshaling
	// because we want to preserve comments and layout of the body.

	strContent := string(content)

	// Pattern: asdp_version: ... (at start of line)
	re := regexp.MustCompile(`(?m)^asdp_version:\s*.*$`)

	if re.MatchString(strContent) {
		newContent := re.ReplaceAllString(strContent, fmt.Sprintf("asdp_version: %s", version))
		if newContent != strContent {
			if dryRun {
				fmt.Printf("[CHECK] %s would be updated to asdp_version: %s\n", path, version)
				return true, nil
			}
			fmt.Printf("Updating %s -> %s\n", path, version)
			return true, os.WriteFile(path, []byte(newContent), 0644)
		}
	} else {
		// If missing, we might want to inject it?
		// For now, only update if present to avoid breaking non-standard files.
	}

	return false, nil
}
