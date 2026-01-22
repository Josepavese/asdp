package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Josepavese/asdp/engine/domain"
	"github.com/Josepavese/asdp/engine/system"
	"github.com/Josepavese/asdp/engine/usecase"
	"github.com/Josepavese/asdp/internal/adapter/mcp"
	"github.com/Josepavese/asdp/validate/check"
)

func main() {
	// Simple CLI args for testing
	queryPath := flag.String("query", "", "Path to query context for (e.g. ./tools/mcp-server)")
	flag.Parse()

	// Load Configuration
	cfg, err := system.LoadConfig("") // Load from default locations
	if err != nil {
		log.Printf("Warning: Failed to load config, using defaults: %v", err)
		cfg = domain.DefaultConfig()
	}

	// Dependency Injection
	fs := system.NewRealFileSystem()
	configLoader := system.NewConfigurationLoader()
	hasher := system.NewSHA256ContentHasher(cfg.Hasher)
	parser := system.NewPolyglotParser(*cfg) // Switched to Polyglot

	queryUC := usecase.NewQueryContextUseCase(fs, hasher, *cfg)
	syncUC := usecase.NewSyncModelUseCase(fs, parser, hasher, cfg.Sync.Model)
	scaffoldUC := usecase.NewScaffoldUseCase(fs, cfg.Scaffold)
	initAgentUC := usecase.NewInitAgentUseCase(fs, *cfg)
	syncTreeUC := usecase.NewSyncTreeUseCase(fs, cfg.Sync.Tree)
	functionUC := usecase.NewGetFunctionInfoUseCase(fs, parser, hasher, *cfg)

	// Mode 1: Query CLI (Testing)
	if *queryPath != "" {
		resp, err := queryUC.Execute(*queryPath)
		if err != nil {
			log.Fatalf("Error querying context: %v", err)
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		encoder.Encode(resp)
		return
	}

	initProjectUC := usecase.NewInitProjectUseCase(initAgentUC, syncTreeUC, scaffoldUC)
	validateUC := check.NewValidateProjectUseCase(fs, parser, hasher, configLoader, cfg)

	// Mode 2: MCP Server (Default)
	fmt.Fprintf(os.Stderr, "ASDP MCP Server v%s started.\n", domain.Version)
	mcpServer := mcp.NewServer(queryUC, syncUC, scaffoldUC, initAgentUC, syncTreeUC, initProjectUC, validateUC, functionUC, *cfg)
	mcpServer.Serve()
}
