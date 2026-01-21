# ASDP: Agentic Spec Driven Programming

ASDP is a protocol designed to bridge the gap between AI Agents and Codebases. It enforces a documentation-first workflow where "The Spec is the Truth."

> **Status**: Core Implementation Complete (v0.1.0)

## Project Structure

- **`core/`**: The "Standard Library" of ASDP (Protocol definitions, Rules, Skills).
  - **`spec/`**: Formal definitions (`codespec`, `codemodel`, `codetree`).
  - **`agent/`**: Rules and Workflows for AI agents.
- **`tools/`**: The reference implementation ("The Kit").
  - **`engine/`**: The shared Go library implementing the ASDP logic.
  - **`mcp-server/`**: The MCP server that exposes ASDP capabilities.
- **`installer/`**: Cross-platform installation scripts.

You can install the ASDP CLIs using the provided scripts.

> [!NOTE]
> **Private Repository Note**: Since this repository is private, `curl` commands might fail with a 404 unless authenticated.
> We recommend using the **GitHub CLI (`gh`)** for a zero-friction installation.

### Method 1: GitHub CLI (Recommended)

```bash
# Clone the repository
gh repo clone Josepavese/asdp
cd asdp
# Run the installer
./installer/install.sh
```

### Method 2: Manual Script Execution

**Linux/macOS**:

```bash
./installer/install.sh
```

**Windows (PowerShell)**:

```powershell
./installer/install.ps1
```

## Features

The `asdp` binary (MCP Server) provides the following tools:

1. **`asdp_query_context`**: Reads `codespec.md` and `codemodel.md`, verifying their freshness.
2. **`asdp_sync_codemodel`**: Automatically parses source code (Go + Polyglot Ctags) and updates `codemodel.md`.
3. **`asdp_scaffold`**: Creates new ASDP-compliant modules with standard templates.

## Contributing

This project itself follows ASDP.

- `tools/codetree.md`: The root map of the tools.
- `tools/mcp-server/codespec.md`: The spec for the server.
