# CodeTree Format Definition

**Filename**: `codetree.md`
**Location**: Only at the **Root** of the repository (or major sub-projects).
**Purpose**: Defines the recursive hierarchy of the entire project. Maps the filesystem to logical components.

## Format: YAML/JSON (Strict SBOM)

Unlike `codespec` and `codemodel` which are hybrid, `codetree.md` is strictly structural. It MAY be a pure `.json` or `.yaml` file, but we stick to `.md` (Frontmatter only) for consistency if preferred, OR we can mandate `codetree.yaml`.

**Decision**: We currently define it as **`codetree.md` (Frontmatter Only)** to keep the `.md` consistency across the repo.

### 1. YAML Frontmatter (CycloneDX Subset)

We use a simplified schema inspired by **CycloneDX** (Component Object).

```yaml
---
asdp_version: "1.0.0"

# Root Marker: DEFINES an ASDP Island.
# If true, tools stop traversing up when they hit this file.
# This allows having multiple ASDP roots in a monorepo (e.g. legacy/ and new/).
root: true

# Tree Structure
components:
  - name: "cmd"
    type: "application"
    path: "./cmd"
    children:
      - name: "server"
        path: "./cmd/server"
        description: "Main HTTP Server Entrypoint"
        has_spec: true
        has_model: true

  - name: "pkg"
    type: "library"
    path: "./pkg"
    children:
      - name: "auth"
        path: "./pkg/auth"
        description: "Core authentication logic"
        # ASDP Meta-properties
        has_spec: true   # If true, expects ./auth/codespec.md
        has_model: true  # If true, expects ./auth/codemodel.md
        components:
          - name: "oauth"
            has_spec: true
            has_model: true
      
      - name: "db"
        path: "./pkg/db"
        description: "Postgres Connection Pool"
        has_spec: true
        has_model: true

# Verification
verification:
  scan_time: "2023-10-27T10:00:00Z"

# Exclusions
# List of glob patterns or path prefixes to exclude from the ASDP tree and context.
excludes:
  - "dist"
  - "tmp"
  - "legacy/generated"
---
```

### 2. Markdown Body (Optional)

Usually empty, or contains a rendered tree view (ASCII) for human delight.

```markdown
# Project Tree
.
├── cmd
│   └── server
└── pkg
    ├── auth
    └── db
```
