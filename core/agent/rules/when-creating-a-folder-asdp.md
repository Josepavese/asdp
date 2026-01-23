---
trigger: always_on
description: Protocol to follow whenever creating a new folder or ASDP module.
---

## RULE — WHEN_CREATING_A_FOLDER

When a new folder is created:

1. Add the folder entry to `codetree.md`.
2. Create `<new-folder>/codespec.md` and `<new-folder>/codemodel.md` from spec.
    > **Tool**: `asdp_scaffold` (creates folder + templated docs).
3. Update parent `codespec.md` and `codemodel.md` recursively up to root if module boundaries, responsibilities, or elements change.
4. Append changelog entries to all updated documents.
