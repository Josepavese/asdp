---
asdp_version: 1.0.0
id: "mcp-server"
type: "application"
title: "mcp-server"
summary: "The Model Context Protocol Server implementation."
capabilities:

- "MCP Server"
dependencies:
  - module: "engine"
    reason: "Core logic"
requirements: []
exports: []

---

# MCP Server Specification

## Context

This directory contains the specific implementation of the ASDP Agent as an MCP Server. It bridges the generic `engine` logic with the MCP transport protocol.

## Requirements

- Must expose ASDP tools over Stdio.
dependencies:
  - module: "go.std"
    reason: "Core runtime"
requirements:
  - id: "REQ-CORE-02"
    desc: "Must support querying ASDP context from disk"
exports:
  - "main"

---

# Context

This server enables AI Agents to navigate the codebase using high-level semantic tools.
