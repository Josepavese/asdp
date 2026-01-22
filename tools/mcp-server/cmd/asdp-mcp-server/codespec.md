---
asdp_version: "0.1.16"
id: "asdp-mcp-server"
type: "application"
title: "asdp-mcp-server"
summary: "The main entrypoint for the ASDP MCP Server application."
capabilities:
  - "Binary Entrypoint"
  - "Dependency Injection Wiring"
  - "Server Startup"
dependencies:
  - module: "adapter/mcp"
    reason: "Wiring"
  - module: "engine/system"
    reason: "Wiring"
  - module: "engine/usecase"
    reason: "Wiring"
requirements: []
exports: []
---
# Action: ASDP MCP Server

## Context

This module contains the `main.go` file for the MCP server. Its sole responsibility is to wire together the dependencies (System -> UseCase -> Adapter) and start the server process.

## Requirements

- Must compile to a standalone binary.
- Must configure the server to run on Stdio.
