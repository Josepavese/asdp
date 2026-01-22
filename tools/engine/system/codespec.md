---
asdp_version: "0.1.16"
id: "<no value>"
type: "library"
title: "system"
summary: "Provides infrastructure implementations for ASDP interfaces."
capabilities:
  - "RealFileSystem implementation"
  - "SHA256 Content Hasher"
  - "Polyglot, GoAST, and Ctags parsers"
dependencies:
  - module: "domain"
    reason: "Implements interfaces"
requirements: []
exports:
  - "NewRealFileSystem"
  - "NewSHA256ContentHasher"
---
# System Specification

## Context

This module acts as the Infrastructure Layer. It implements the interfaces defined in `domain`, providing concrete capabilities like file system access (os/io), cryptographic hashing, and source code parsing code (go/parser, ctags).

## Requirements

- Must implement `domain` interfaces.
- Should handle OS-specific details.
