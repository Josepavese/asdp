## RULE — WHEN_CREATING_A_FOLDER

When a new folder is created:

1. Add the folder entry to `codetree.md`.
2. Create `<new-folder>/codespec.md` using `spec/asdp/codespec.md`.
3. Create `<new-folder>/codemodel.md` using `spec/asdp/codemodel.md` when the folder has code elements.
4. Update parent `codespec.md` and `codemodel.md` recursively up to root if module boundaries, responsibilities, or elements change.
5. Append changelog entries to all updated documents.
