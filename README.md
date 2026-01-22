# ASDP: Agentic Spec Driven Programming Protocol

> [!WARNING]
> **EXPERIMENTAL PROTOCOL**: ASDP is currently an experimental research project. It is being utilized to evaluate the efficacy of AI agents in autonomous software development and maintenance. Expect breaking changes and architectural evolutions as we refine the protocol's foundations.

ASDP (Agentic Spec Driven Programming) is a formal protocol designed to eliminate the semantic gap between high-level architectural intent and low-level code implementation. By enforcing a documentation-centric state machine, ASDP ensures that **"The Specification is the Immutable Truth,"** providing AI agents with a deterministic context for reasoning and action.

## Theoretical Foundations

The ASDP protocol is built upon several key pillars of computer science and software engineering research:

1. **Agent-Oriented Software Engineering (AOSE)**: Drawing from the foundational work of Wooldridge and Jennings (*"Agent-Oriented Software Engineering"*, 2000), ASDP treats the development process as a set of autonomous interactions where the agent must maintain an accurate internal model of the environment (the codebase).
2. **Intent-Based Programming**: Inspired by the paradigm where high-level goals drive system behavior, ASDP mandates that every module must be preceded by a `codespec.md`, which defines the "Contract of Intent" before any logic is implemented.
3. **Context Management & Freshness**: Addressing the "Stale Context Problem" in Large Language Model (LLM) reasoning, ASDP introduces integrity checks (Freshness Status) to ensure that the agent's understanding of the code's structure (`codemodel.md`) is perfectly synchronized with the actual source.

### Scientific Evaluation

We are currently conducting empirical tests to measure how the enforcement of ASDP boundaries affects:

- **Agent Autonomy**: The ability to perform complex refactors without human intervention.
- **Error Propagation**: Reducing the rate at which "hallucinated" API calls or logic bugs are introduced.
- **Context Efficiency**: Minimizing the required token count for an agent to understand a complex module.

## The ASDP Architecture

### Project Structure

- **`core/`**: The protocol's definition layer. Contains standard rules, workflows, and schema definitions for `codespec`, `codemodel`, and `codetree`.
- **`tools/`**: The Reference Implementation ("The Kit"). A high-performance Go engine and MCP server that provides the operational interface for agents.
- **`installer/`**: Deployment vector for integrating ASDP into developer environments.

### The Toolbelt (MCP Integration)

The ASDP engine exposes a suite of specialized tools designed for agentic consumption:

1. **`asdp_query_context`**:
    - **Function**: Retrieves the unified context (Spec + Model + Freshness) for a directory.
2. **`asdp_sync_codemodel`**:
    - **Function**: Performs static analysis of the source code to update the `codemodel.md`.
3. **`asdp_sync_codetree`**:
    - **Function**: Recursively scans the project to update the global `codetree.md`.
4. **`asdp_scaffold`**:
    - **Function**: Generates compliant module structures from templates.

## Installation

ASDP can be installed via a single command. The installer will automatically configure the environment and optional agent-ready assets.

### Linux / macOS

```bash
curl -sSL https://raw.githubusercontent.com/Josepavese/asdp/main/installer/install.sh | bash
```

### Windows (PowerShell)

```powershell
powershell -ExecutionPolicy Bypass -Command "iwr -useb https://raw.githubusercontent.com/Josepavese/asdp/main/installer/install.ps1 | iex"
```

## Contributing to the Experiment

This project is self-bootstrapping and adheres strictly to the ASDP protocol.

- Refer to [codetree.md](file:///home/jose/hpdev/Libraries/asdp/codetree.md) for the global architecture.
- Every tool in `tools/` is governed by its own `codespec.md`.

---
*Developed by Josè Pavese (Experimental Agentic Research)*
