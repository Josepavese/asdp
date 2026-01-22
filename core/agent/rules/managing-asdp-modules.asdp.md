---
description: Rules for identifying, creating, and maintaining ASDP modules.
groups: ["asdp", "maintenance"]
---

# Managing ASDP Modules

These rules govern how the AI must identify, create, and maintain ASDP modules within a project.

## 1. Module Identification (Smart Analysis)

**Rule**: "Only Significant Directories are Modules."

To prevent file explosion in deep hierarchies (e.g., `a/b/c/d/file.go`), you must apply the **Significance Test**.

A directory IS an ASDP Module (and needs a `codespec.md`) ONLY IF it meets one of these criteria:

1. **The Leaf (Code Bearing)**: containing actual source code files directly (e.g., `.go`, `.ts`).
2. **The Hub (High Traffic)**: containing **MULTIPLE** immediate sub-directories that are themselves modules or significant packages.
3. **The Root**: The project root (where `.agent` resides).
4. **The Explicit Boundary**: A layer the user explicitly flagged (e.g., `domain`, `adapter`) even if empty of code.

**The Pass-Through Exception**:
If a directory contains *no code* and only *one* sub-directory (e.g., `java/com/example`), it is a **Namespace**.

* **Action**: Do NOT generate `codespec.md` or `codemodel.md`. Let the child module bubble up in the CodeTree naturally.

## 2. The Generation Lifecycle (Order Enforcement)

**Rule**: "Follow the Holy Order of Generation."

You **MUST** strictly follow this sequence when initializing or syncing a module. **NEVER** skip a step or change the order.

1. **STEP 1: CodeSpec (`codespec.md`) - The INTENT.**
    * *Why*: Defines who the module is. Without this, the module is undocumented.
    * *Action*: Run `asdp_scaffold`.
    * *Strictness*: Must be created **BEFORE** any other ASDP file.

2. **STEP 2: CodeModel (`codemodel.md`) - The REALITY.**
    * *Why*: Inspects the code to see if it matches the intent.
    * *Action*: Run `asdp_sync_codemodel`.
    * *Strictness*: Run only **AFTER** `codespec.md` exists and code is present.

3. **STEP 3: CodeTree (`codetree.md`) - The MAP.**
    * *Why*: Aggregates the system state.
    * *Action*: Run `asdp_sync_codetree` (usually at the project root).
    * *Strictness*: Run **LAST**, to capture the fresh data from Steps 1 and 2.

## 3. Contextualization (Anti-Laziness)

**Rule**: "Kill all TODOs immediately."

* **Prohibition**: You are **FORBIDDEN** from leaving default values like `id: "<no value>"` or `summary: "TODO..."` in a `codespec.md`.
* **Mandate**: Immediately after scaffolding a spec, you **MUST** call `replace_file_content` to populate:
  * `id`: The logical name (e.g., `tools.engine`).
  * `summary`: A rich description (used by CodeTree). "Container for X" is better than just "X".
  * `dependencies`: Real relationships.

## 4. Safe Editing Protocol (Parsing Integrity)

**Rule**: "Anchor your edits."

* **Problem**: Replacing only `id: ...` can accidentally create duplicate frontmatter blocks if the file is scanned multiple times.
* **Solution**: When editing YAML frontmatter, **ALWAYS** include the `---` delimiters or a large unique block of context in your `TargetContent`.
  * *Bad*: Target `id: "<no value>"`
  * *Good*: Target `asdp_version: "..."\nPd: "<no value>"` (Include surrounding lines).

## 5. Syntax Enforcements

**Rule**: "Strict Syntax Only."

* **Dependencies**: MUST use object syntax. String lists are forbidden.

    ```yaml
    # CORRECT
    dependencies:
      - module: "engine"
        reason: "Core logic"
    
    # INCORRECT (Forbidden)
    dependencies:
      - "engine"
    ```
