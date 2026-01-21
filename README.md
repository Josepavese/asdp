# ASDP: Agentic Spec Driven Programming

ASDP is a protocol designed to bridge the gap between AI Agents and Codebases. It enforces a documentation-first workflow where "The Spec is the Truth."

> **Status**: Core Implementation Complete (v0.1.13)

## Project Structure

- **`core/`**: The "Standard Library" of ASDP (Protocol definitions, Rules, Skills).
  - **`spec/`**: Formal definitions (`codespec`, `codemodel`, `codetree`).
  - **`agent/`**: Rules and Workflows for AI agents.
- **`tools/`**: The reference implementation ("The Kit").
  - **`engine/`**: The shared Go library implementing the ASDP logic.
  - **`mcp-server/`**: The MCP server that exposes ASDP capabilities.
- **`installer/`**: Cross-platform installation scripts. These scripts now include an **interactive initialization** feature that optionally sets up the `.agent/` folder in your current directory.

### Method 1: Automatic (Native Tools)

**Linux/macOS**:

```bash
# Set your token and run the installer
GITHUB_TOKEN=your_token_here ./installer/install.sh
```

**Windows (PowerShell)**:

```powershell
# Set your token and run the installer
$env:GITHUB_TOKEN='your_token_here'; .\installer\install.ps1
```

### Method 2: Manual (GitHub CLI)

If you have the `gh` CLI installed and authenticated:

```bash
gh repo clone Josepavese/asdp
cd asdp
./installer/install.sh
```

## Features

The `asdp` binary (MCP Server) provides the following tools:

1. **`asdp_query_context`**: Reads `codespec.md` and `codemodel.md`, verifying their freshness.
2. **`asdp_sync_codemodel`**: Automatically parses source code (Go + Polyglot Ctags) and updates `codemodel.md`.
3. **`asdp_sync_codetree`**: Automatically scans the project structure and generates/updates `codetree.md`.
4. **`asdp_scaffold`**: Creates new ASDP-compliant modules with standard templates.

## Contributing

This project itself follows ASDP.

- `tools/codetree.md`: The root map of the tools.
- `tools/mcp-server/codespec.md`: The spec for the server.
