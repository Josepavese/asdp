# codetree.md SPEC v1.0

File MUST follow exactly:

# CodeTree

- `<folder>/` – <1-line description>. See `<folder>/codespec.md`.
  - `<subfolder>/` – <1-line description>. See `<subfolder>/codespec.md`.
  - ...
    Rules:
- Tree MUST be recursive: include every project folder down to the deepest level.
- List ONLY folders.
- Each line: `<path>/ – <1-line functional description>. See <path>/codespec.md.`
- No files. No build/cache/output dirs.
