---
description: Verify ASDP compliance and rebuild missing CodeTree/CodeSpec/CodeModel docs across the project.
---
1. Stop coding; goal is documentation conformance only.
2. Read `~/.asdp/core/spec/codespec.md`, `~/.asdp/core/spec/codemodel.md`, `~/.asdp/core/spec/codetree.md` to follow formats exactly.
3. From repo root, inspect existing docs: `README.md`, `rules/`, `codetree.md`, `codespec.md`, `codemodel.md` (recreate root docs per spec if missing).
4. Map project folders to include; exclude vendor/derived dirs (`node_modules`, `.git`, `.svn`, `.hg`, `.venv`, `.tox`, `__pycache__`, `.pytest_cache`, `dist`, `build`, `target`, `.turbo`, `coverage`, `*.egg-info`).
5. For each project folder (depth-first):
   - Read parent `codespec.md` and `codemodel.md` to understand intent/scope.
   - Create/update that folder’s `codespec.md` per `~/.asdp/core/spec/codespec.md`, reflecting actual role and boundaries only.
   - Create/update that folder’s `codemodel.md` per `~/.asdp/core/spec/codemodel.md`, listing real elements only (no invention).
     > **Tip**: Use `asdp_sync_codemodel` to automatically generate the symbol list and source hash for this step.
   - If folder structure changed, update root `codetree.md`.
6. After traversal, review all generated docs together for coherence; ensure code/structure and docs match exactly.
7. Append `YYYY-MM-DD – agent – change` changelog entries in every updated document.
