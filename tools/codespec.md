```
---
asdp_version: "0.1.15"
id: "tools"
type: "library"
title: "tools"
summary: "Container for all ASDP tooling and executables."
capabilities:
  - "Engine Library"
  - "MCP Server Application"
dependencies: []
requirements: []
exports: []
---
# Tools Specification

## Context
This directory acts as the root for all executable tools and core libraries that make up the ASDP implementation. It enforces a separation between the abstract protocol definition (in `core`) and the concrete implementation (in `tools`).

## Requirements
- Must contain all Go code.
```
