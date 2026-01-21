# PROPOSAL: ASDP MCP Server (Agent Systems Documentation Protocol)

**Date**: 2026-01-21
**Status**: DRAFT
**Context**: Nido Project Improvement
**Target Audience**: Developers building the ASDP automation suite.

---

## 1. Context & Problem Statement

### The Challenge: "Agentic Amnesia" & "Context Drift"

In large codebases, AI Agents face significant challenges:

1. **Limited Context Window**: Agents cannot "see" the whole project at once. They rely on expensive and slow retrieval mechanisms (grep, file reading).
2. **Drift**: Documentation (High-level intent) and Code (Low-level implementation) naturally drift apart over time.
3. **Bureaucractic Overhead**: The ASDP protocol (Codespec, Codemodel, Codetree) solves the context problem by creating a predictable "semantic filesystem". However, manually maintaining these files is tedious, error-prone, and consumes valuable Agent tokens/steps.

### The Solution

Build a dedicated **Model Context Protocol (MCP) Server** focused on **ASDP Automation**.
This server will provide a suite of **deterministic, script-based tools** (Zero AI) that handle the structural and mechanical aspects of documentation maintenance.

By offloading the "mechanical" validaton and synchronization to fast, local scripts, the AI Agent is freed to focus on "reasoning" and "semantic" content.

---

## 2. Objectives

1. **Reduce Agent Step Count**: Convert multi-step manual checks (list dir -> read file -> compare) into single atomic tool calls (`asdp_validate`).
2. **Guarantee Structural Integrity**: Ensure `codetree.md` is always recursive and `codemodel.md` signatures exactly match the source code AST (Abstract Syntax Tree).
3. **Zero-Friction Compliance**: Make it easier to follow the rules than to break them.

---

## 3. Architecture & Tech Stack

- **Language**: Go (Recommended for native integration with Nido's codebase and AST parsing capabilities) or Python.
- **Nature**: Stateless, deterministic CLI wrappers exposed via MCP Standard.
- **Intelligence**: **None**. The tools rely on file system operations, regex, and AST parsing. They do not infer intent.

---

## 4. Proposed Toolset (The API)

The MCP Server should expose the following tools to the Agent:

### A. Discovery Tools (Read-Only)

#### 1. `asdp_scan_structure`

- **Purpose**: Fast replacement for `codetree.md` walking.
- **Input**: `root_path` (optional).
- **Logic**:
  - Recursively walks the file system (ignoring `.git`, `node_modules`, `dist`, etc.).
  - Returns a JSON tree of the project structure.
  - Annotates each node with "HasCodeSpec: boolean", "HasCodeModel: boolean".
- **Agent Use Case**: "I need to understand the project layout to find where `Auth` logic resides."

#### 2. `asdp_read_context`

- **Purpose**: "One-shot" context loading for a specific module.
- **Input**: `path` (e.g., `internal/net`).
- **Logic**:
  - Reads `codespec.md` and `codemodel.md` of the target folder.
  - Reads the `Purpose` section of the parent's `codespec.md`.
  - Returns a consolidated JSON/Markdown blob.
- **Agent Use Case**: "I'm about to edit `internal/net`. load me the full context."

### B. Maintenance Tools (Write/Edit)

#### 3. `asdp_scaffold_module`

- **Purpose**: Create a new module with compliant boilerplate.
- **Input**: `path`, `purpose_one_liner`.
- **Logic**:
  - `mkdir -p <path>`
  - Creates `codespec.md` with standard template.
  - Creates `codemodel.md` with standard template.
  - Updates root `codetree.md` inserting the new line in the correct alphabetical/hierarchical position.
- **Agent Use Case**: "Create a new module `cmd/new-tool`."

#### 4. `asdp_sync_codemodel` (The "Auto-Doctor")

- **Purpose**: Sync code signatures with documentation. This is the **most critical tool**.
- **Input**: `path` (folder to sync).
- **Logic**:
  - **Parse Code**: Uses Go/Python AST parsers to find all exported Functions, Structs, Classes, Interfaces in the `.go` / `.py` files.
  - **Parse Doc**: Parses existing `codemodel.md`.
  - **Merge**:
    - *Match*: Update signature in Doc if changed in Code. Keep existing Description.
    - *New in Code*: Add to Doc (Description: "TODO: Agent fill this").
    - *Deleted in Code*: Mark as `[DEPRECATED]` or remove from Doc.
  - **Write**: Save updated `codemodel.md`.
- **Agent Use Case**: "I just finished refactoring `internal/tui`. Sync the docs."

#### 5. `asdp_validate_integrity`

- **Purpose**: CI/CD style check.
- **Logic**:
  - Checks if every folder has specs.
  - Checks if `codetree` matches FS.
  - Checks if `codemodel` matches FS.
- **Output**: List of violations (Errors/Warnings).
- **Agent Use Case**: "Run a health check before I submit my PR."

---

## 5. Implementation Roadmap

### Phase 1: The "Linter" (Read-Only)

- Implement `asdp_validate_integrity` logic.
- Result: An MCP tool that tells the agent *what* is broken, but doesn't fix it.

### Phase 2: The "Scaffolder" (FS Ops)

- Implement `asdp_scaffold_module`.
- Implement `codetree` auto-updating logic.

### Phase 3: The "Parser" (AST Integration)

- Implement `asdp_sync_codemodel`.
- Requires integrating `go/ast` (if project is Go) or generic regex parsers.

---

## 6. Example Interaction Flow

**User**: "Refactor the `NewServer` function in `pkg/server` to accept a Context."

**Agent**:

1. **Call `asdp_read_context(pkg/server)`**. Understands current design.
2. **Edit `pkg/server/server.go`**. Changes the code signature.
3. **Call `asdp_sync_codemodel(pkg/server)`**. The script automatically updates the signature in MD.
4. **Agent reads diff**. Sees the signature updated, adds a short description to the changelog.
5. **Done**.

Context kept in sync with near-zero effort.
