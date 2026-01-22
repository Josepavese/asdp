---
asdp_version: "0.1.15"
id: "engine"
type: "library"
title: "engine"
summary: "The implementation core of ASDP logic."
capabilities:
  - "Clean Architecture Layers"
dependencies:
  - module: "core"
    reason: "Implements spec"
requirements: []
exports: []
---
# Engine Specification

## Context

This directory houses the "Engine", which is the Go implementation of the ASDP protocol capabilities. It follows Clean Architecture, divided into `domain`, `usecase`, and `system`.

## Requirements

- Must be importable by other tools (e.g. MCP Server, CLI).
