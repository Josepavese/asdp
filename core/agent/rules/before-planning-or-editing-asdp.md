---
trigger: always_on
description: Mandatory checks before planning or editing work.
---

# ASDP RULES — Before Planning or Editing (v1.0)

These steps MUST ALWAYS be followed before planning or editing.

- Read root and target `codespec.md`/`codemodel.md`.
    > **Tool**: `asdp_query_context` (checks freshness vs src_hash).
- If misaligned, stop and correct docs.
    > **Skill**: `asdp-doc-rebuilder` (uses `asdp_sync_codemodel`).
- Do not start until you fully understand module role, APIs, structure, and dependencies.
