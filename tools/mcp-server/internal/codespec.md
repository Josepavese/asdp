---
asdp_version: 1.0.0
id: "internal"
type: "library"
title: "internal"
summary: "Internal implementation details for the MCP server."
capabilities:
  - "Adapters"
dependencies: []
requirements: []
exports: []
---
# Internal Specification

## Context

This directory holds internal Go packages that are private to the `asdp-mcp-server` application, preventing them from being imported by external modules.

## Requirements

- Enforce privacy.
