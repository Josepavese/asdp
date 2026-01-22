# CodeSpec Format Definition

**Filename**: `codespec.md`
**Location**: In the root of every logical module/folder.
**Purpose**: Defines the *Intent*, *Requirements*, and *Public Interface* of a module.

## Format: Hybrid YAML+Markdown

The file MUST consist of two parts:

1. **Header**: Strict YAML Frontmatter containing machine-readable metadata.
2. **Body**: Free-form Markdown for human reasoning and context.

### 1. YAML Frontmatter (Strict Schema)

The frontmatter MUST adhere to the following structure.

```yaml
---
# ASDP Protocol Version
asdp_version: 1.0.0

# Unique Identifier for this module (dot notation recommended)
id: "pkg.network.http"

# Module Type: 'library', 'application', 'service', 'interface'
type: "library"

# Short, one-line summary of what this module does
title: "HTTP Client Wrapper"

# High-level abstract/summary for AI Context (Max 200 chars)
summary: "Provides a resilient HTTP client with auto-retry and logging middleware."

# List of precise capabilities provided by this module
capabilities:
  - "http-request-handling"
  - "retry-logic"

# Dependencies (Logical, not just package imports)
dependencies:
  - module: "core/domain"
    reason: "Uses domain types"
  - module: "core/system"
    reason: "Uses file system"
  - module: "pkg.logger"
    reason: "Logging failures"

# Functional Requirements
requirements:
  - id: "REQ-001"
    desc: "Must retry 3 times on 5xx errors"
    priority: "high"
  - id: "REQ-002"
    desc: "Must support custom headers"

# Public Interface Contract (High Level)
exports:
  - "Client"
  - "NewClient"
---
```

### 2. Markdown Body (Human Context)

After the `---` separator, the file continues as standard Markdown. This section is for:

- Architecture Decisions Records (ADR).
- Usage Examples.
- nuance that is hard to capture in YAML.
- "Thinking" process of the Agent/Developer.

#### Example Body

```markdown
# Context & Rationale

We chose to wrap the standard `http.Client` because we needed consistent retry logic across all microservices.

## Architecture
The client functions as a middleware chain...

## Usage
...
```
