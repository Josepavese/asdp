---
description: "Workflow to synchronize ASDP specifications (Spec -> Model -> Tree)"
---

# ASDP Spec Sync Workflow

This workflow enforces the standard lifecycle for ASDP modules: **Spec -> Model -> Tree**.

## Prerequisite: Anchoring

1. Check if the project is anchored (has `.agent` and `codetree.md` at root).
2. If not, refer to `asdp-init.md` first.

## Phase 1: SPECIFICATION (The Intent)

**Goal**: Ensure every **Significant** module has a valid `codespec.md`.

1. **Analyze Structure**: Walk the directory tree.
2. **Apply Significance Test**:
    - **Is it a LEAF?** (Has `.go` files?) -> **YES**.
    - **Is it a HUB?** (Has >1 sub-folders?) -> **YES**.
    - **Is it a PASS-THROUGH?** (Has 1 sub-folder, no code?) -> **NO (SKIP)**.
3. **Scaffold (If Significant & Missing)**:
    - Run `asdp_scaffold`.
4. **Contextualize (CRITICAL)**:
    - **IMMEDIATE ACTION**: Read the new `codespec.md`.
    - **EDIT**: Use `replace_file_content` to fill `id`, `summary`, and `dependencies`.
    - **RULE**: Do NOT leave "TODO" values.
    - **SYNTAX**: Ensure dependencies are objects: `- { module: "x", reason: "y" }`.

## Phase 2: MODELING (The Reality)

**Goal**: Update `codemodel.md` to reflect the actual source code.

1. **Sync Model**:
    - Run `asdp_sync_codemodel(path="/abs/path/to/module")`.
    - This will update symbols, integrity hashes, and `last_modified` timestamps.

## Phase 3: MAPPING (The Map)

**Goal**: Update the project-wide `codetree.md`.

1. **Sync Tree**:
    - Run `asdp_sync_codetree(path="/abs/path/to/project_root")`.
    - This aggregates the rich `summary` from Phase 1 and `last_modified` from Phase 2.

## Verification

- Read `codetree.md`.
- Verify that the module appears with a correct type and a rich, human-readable description (not just its name).
