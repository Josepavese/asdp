---
trigger: always_on
description: Rules for excluding folders from ASDP analysis.
---

# Managing Exclusions (Tree Pruning)

**Rule**: "Keep the Tree Clean and Focused."

You have the power to exclude folders or entire branches from the ASDP protocol. Use this to reduce noise and save context window.

## When to Exclude

Exclude a folder if:
1. **It is Temporary**: Build artifacts (`dist`, `build`, `bin`), temp files (`tmp`), or cache (`.cache`).
2. **It is Self-Explanatory**: Standard frameworks structures that don't need deep context (e.g., standard Angular `node_modules`, simple asset folders).
3. **It is Legacy/Irrelevant**: Old code that is not being touched and shouldn't pollute the context.
4. **It is Generated**: Auto-generated code that shouldn't be edited manually.

## How to Exclude

**NEVER** edit `codetree.md` manually to add exclusions.
**ALWAYS** use the `asdp_manage_exclusions` tool.

Example:
`asdp_manage_exclusions` with arguments: `path="/abs/root", target="dist", action="add"`

This tool automatically updates the configuration and regenerates the `codetree.md`.

## Handling "Excluded" Folders

Once excluded, a folder is invisible to `asdp_sync_codetree` and `asdp_query_context`.
If you need to work on it later, you must `remove` the exclusion first.
