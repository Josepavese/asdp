---
asdp_version: "0.1.15"
id: "mcp"
type: "library"
title: "mcp"
summary: "Adapts ASDP UseCases to the Model Context Protocol (MCP)."
capabilities:
  - "MCP Server Implementation"
  - "Tool Handlers (ServeTool)"
  - "JSON-RPC Communication"
dependencies:
  - module: "engine/usecase"
    reason: "Calls business logic"
requirements: []
exports:
  - "NewServer"
---
# MCP Adapter Specification

## Context

This module acts as the Interface Adapter Layer. It converts calls from the external world (via MCP JSON-RPC) into calls to the internal `usecase` layer, and converts the results back. It isolates the core logic from the transport mechanism.

## Requirements

- Must implement the MCP specification.
- Must handle protocol errors correctly.
