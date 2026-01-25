---
name: asdp-feature-workflow
description: Workflow for creating or modifying ASDP features. Enforces testing.
---

# ASDP Feature Workflow

Follow this process whenever you add or modify a feature in the ASDP protocol.

## 1. Plan & Implement
- Modify Specification (`core/spec/...`).
- Modify Domain/Entities (`tools/engine/domain/...`).
- Implement Logic (`tools/engine/usecase/...`).
- Update MCP Server (`tools/mcp-server/...`).

## 2. Add Test Coverage
- You **MUST** add a new test scenario in `tools/validate/functional_test.go`.
- The test must verify the new tool or logic change (e.g., check that a new flag works, or output contains new data).

## 3. Verify Local Installation
- **Rebuild & Install**:
  ```bash
  # Go to server cmd
  cd tools/mcp-server/cmd/asdp-mcp-server
  go build -o asdp
  
  # Install (assuming ~/.asdp structure)
  mv ~/.asdp/bin/asdp ~/.asdp/bin/asdp.old
  cp asdp ~/.asdp/bin/
  
  # Update Assets
  cd ../../../.. # Back to root
  cp -r core/* ~/.asdp/core/
  ```

## 4. Run System Tests
- Invoke the test skill:
  `view_file core/agent/skills/test-system/SKILL.md` -> Follow instructions.

## 5. Documentation
- Update `walkthrough.md`.
- Create/Update Agent Rules in `core/agent/rules`.
