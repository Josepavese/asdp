---
asdp_version: 1.0.0
id: "adapter"
type: "library"
title: "adapter"
summary: "Interface Adapters for the MCP server."
capabilities:
  - "MCP Adapter"
dependencies: []
requirements: []
exports: []
---
# Adapter Specification

## Context

This directory contains the Interface Adapters that connect the Clean Architecture UseCases to the outside world (validating inputs, formatting outputs).

## Requirements

- Decouple UseCases from Transport.
