# ASDP RULES — After Editing or Task Complete (v1.0)

These steps MUST ALWAYS be followed at task end.

- Always update `codespec.md` and `codemodel.md` in every folder touched by the task (do this at task end).
- From each changed folder, walk up to root and update parent `codespec.md` and `codemodel.md` if roles or elements are impacted.

1. Update the `codemodel.md` file for any modified folders to reflect API changes, new symbols, or dependency changes.
    > **Tip**: Use `asdp_sync_codemodel` (e.g., `{"path": "./my-folder"}`) to auto-parse the code and update the `symbols` and `src_hash` fields.
2. If the folder structure changed (files added/removed/renamed), update `codespec.md` and parent `codetree.md`. and also follow `WHEN_CREATING_A_FOLDER` or `WHEN_DELETING_A_FOLDER` as applicable.

- If functions/classes/structures changed (add/remove/rename/signature-change), update the relevant `codemodel.md`.
- A task is NOT complete until all documentation is fully aligned and consistent.
- Every updated document MUST receive a `YYYY-MM-DD – agent – change` entry.
