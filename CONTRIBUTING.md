# Contributing to the ASDP Protocol

Thank you for your interest in contributing to the Agentic Spec Driven Programming (ASDP) protocol. This is an experimental research project dedicated to refining the interaction between autonomous AI agents and complex codebases.

> [!IMPORTANT]
> **Protocol-First workflow**: ASDP is not just a repository; it is a system of boundaries. All contributions must adhere strictly to the "Speccing before Coding" paradigm.

## Scientific Rationale

The core goal of ASDP is to minimize the semantic entropy in software development. For our research to be valid, every change must be documented in a way that is deterministic for both humans and AI agents.

## The Contribution Lifecycle

Every feature or fix must follow this exact state transition:

1. **Spec Initialization**: Before writing logic, you must update or create the relevant `codespec.md`. This defines the *Contract of Intent*.
2. **Structural Modeling**: Use the `asdp_sync_codemodel` tool to ensure the `codemodel.md` reflects the architectural changes.
3. **Implementation**: Write the source code, ensuring it fulfills the specifications defined in step 1.
4. **Verification**: Execute the global validation suite:

    ```bash
    go run tools/validate/main.go
    ```

5. **Synchronization**: Ensure all versioning metadata is consistent:

    ```bash
    go run tools/cmd/version-manager/main.go sync
    ```

## PR Submission Requirements

A Pull Request will only be considered for review if it meets the following "Definition of Done":

- [ ] **ASDP Compliance**: Any new or modified module must have a corresponding `codespec.md` and `codemodel.md`.
- [ ] **Freshness Integrity**: The `codemodel.md` must be synchronized with the source code.
- [ ] **Validation Pass**: The PR must pass all existing functional tests.
- [ ] **Architectural Alignment**: The change must be reflected in the global `codetree.md` if it alters the project hierarchy.

## Standards of Excellence

- **Aesthetics Matters**: Documentation should be clean, professional, and use the provided ASDP Markdown templates.
- **Scientific Tone**: Avoid informal language; describe logic and architecture with precision.
- **Commit Messages**: Use [Conventional Commits](https://www.conventionalcommits.org/) (e.g., `feat:`, `fix:`, `chore:`).

---
*By contributing to this repository, you agree to uphold the experimental integrity of the ASDP protocol.*
