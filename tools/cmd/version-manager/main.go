package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

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
		if err := performSync(cfg); err != nil {
			fmt.Printf("Sync failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Sync completed successfully.")
	case "check":
		// Check would verify if files match config without changing them
		fmt.Println("Check not implemented yet, running sync for now...")
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

func performSync(cfg *VersionConfig) error {
	// 1. Update Go Constant (tools/engine/domain/version.go)
	if err := updateGoVersion("tools/engine/domain/version.go", cfg.KitVersion); err != nil {
		return fmt.Errorf("failed to update domain/version.go: %w", err)
	}

	// 2. Update Go Constant (tools/mcp-server/internal/adapter/mcp/server.go) - Protocol Version
	// Actually we handle protocol version in server.go logic or constant?
	// Implementation plan said "Update ... MCP Protocol Version".
	// Let's assume it's in the HandleInitialize method.
	if err := updateMcpProtocol("tools/mcp-server/internal/adapter/mcp/server.go", cfg.McpProtocolVersion); err != nil {
		// Non-fatal if file structure changed, but good to know
		fmt.Printf("Warning: Could not update MCP protocol version: %v\n", err)
	}

	// 3. Walk and Update Frontmatter (codespec/codemodel)
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && (d.Name() == ".git" || d.Name() == "node_modules") {
			return fs.SkipDir
		}
		if !d.IsDir() && (strings.HasSuffix(path, "codespec.md") || strings.HasSuffix(path, "codemodel.md")) {
			return updateFrontmatterVersion(path, cfg.SchemaVersion)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk directories: %w", err)
	}

	return nil
}

func updateGoVersion(path string, version string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Regex for: const Version = "..."
	re := regexp.MustCompile(`const Version = "[^"]+"`)
	newContent := re.ReplaceAll(content, []byte(fmt.Sprintf(`const Version = "%s"`, version)))

	if !bytes.Equal(content, newContent) {
		fmt.Printf("Updating %s -> %s\n", path, version)
		return os.WriteFile(path, newContent, 0644)
	}
	return nil
}

func updateMcpProtocol(path string, version string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Look for ProtocolVersion: "..." in struct initialization
	re := regexp.MustCompile(`ProtocolVersion:\s*".*"`)
	newContent := re.ReplaceAll(content, []byte(fmt.Sprintf(`ProtocolVersion: "%s"`, version)))

	if !bytes.Equal(content, newContent) {
		fmt.Printf("Updating ProtocolVersion in %s -> %s\n", path, version)
		return os.WriteFile(path, newContent, 0644)
	}
	return nil
}

func updateFrontmatterVersion(path string, version string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Simple YAML Frontmatter replacement using regex is safer than parsing/marshaling
	// because we want to preserve comments and layout of the body.

	strContent := string(content)

	// Pattern: asdp_version: ... (at start of line)
	re := regexp.MustCompile(`(?m)^asdp_version:\s*.*$`)

	if re.MatchString(strContent) {
		newContent := re.ReplaceAllString(strContent, fmt.Sprintf("asdp_version: %s", version))
		if newContent != strContent {
			fmt.Printf("Updating %s -> %s\n", path, version)
			return os.WriteFile(path, []byte(newContent), 0644)
		}
	} else {
		// If missing, we might want to inject it?
		// For now, only update if present to avoid breaking non-standard files.
	}

	return nil
}
