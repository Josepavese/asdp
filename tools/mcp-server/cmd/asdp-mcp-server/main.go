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
)

func main() {
	// Simple CLI args for testing
	queryPath := flag.String("query", "", "Path to query context for (e.g. ./tools/mcp-server)")
	flag.Parse()

	// Dependency Injection
	fs := system.NewRealFileSystem()
	hasher := system.NewSHA256ContentHasher()
	parser := system.NewPolyglotParser() // Switched to Polyglot

	queryUC := usecase.NewQueryContextUseCase(fs, hasher)
	syncUC := usecase.NewSyncModelUseCase(fs, parser, hasher)
	scaffoldUC := usecase.NewScaffoldUseCase(fs)
	initAgentUC := usecase.NewInitAgentUseCase(fs)
	syncTreeUC := usecase.NewSyncTreeUseCase(fs)

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

	// Mode 2: MCP Server (Default)
	fmt.Fprintf(os.Stderr, "ASDP MCP Server v%s started.\n", domain.Version)
	mcpServer := mcp.NewServer(queryUC, syncUC, scaffoldUC, initAgentUC, syncTreeUC)
	mcpServer.Serve()
}
