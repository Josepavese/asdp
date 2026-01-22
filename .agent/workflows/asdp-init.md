---
description: Initialize a project with ASDP, identifying the best code root (anchoring) and setting up agent assets.
---

1. **Analyze Project Structure**:
   - Locate the repository root (containing `.git` or project config).
   - Identify the **Code Root**: Look for the primary directory where actual code resides (e.g., `src/`, `tools/`, `pkg/`, `libs/`).
   - **Goal**: Avoid the repository root if it contains mostly noise (config, tests, CI/CD). Answering "Where does the implementation really start?" is key.

2. **Initialize Project Anchor**:
   - Use the `asdp_init_project` tool to set up the protocol.
   - **Parameters**:
     - `path`: The absolute repository root.
     - `code_path`: The absolute path to the identified **Code Root**.
   > **Tool**: `asdp_init_project(path="/repo", code_path="/repo/tools")`

3. **Verify Grounding**:
   - Check that `.agent/` is in the repository root.
   - Check that `codetree.md`, `codespec.md`, and `codemodel.md` are initialized in the **Code Root**.
   - The code root should now have `root: true` in its `codetree.md`.

4. **Onboard AI Agent**:
   - Instruct the agent to start analysis from the **Code Root**'s context.
   > **Tool**: `asdp_query_context(path="/repo/tools")`
