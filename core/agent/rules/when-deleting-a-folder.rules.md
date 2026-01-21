## RULE — WHEN_DELETING_A_FOLDER

When a folder is deleted:

1. Remove the folder entry from `codetree.md`.
2. Delete `<deleted-folder>/codespec.md` and `<deleted-folder>/codemodel.md` if present.
3. Update parent `codespec.md` and `codemodel.md` recursively up to root if module boundaries, responsibilities, or elements change.
4. Remove or update any paths in `codemodel.md` that referenced deleted elements.
5. Append changelog entries to all updated documents.
