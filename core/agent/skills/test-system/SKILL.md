---
name: asdp-test-system
description: Run full system tests for ASDP protocol using the validation suite.
---

# ASDP System Test Skill

This skill runs the functional test suite located in `tools/validate`.
This suite spins up the MCP server in-process and tests all tools against a sandbox.

## Steps

1. **Clean Test Cache**:
   Ensure no stale test results.
   ```bash
   go clean -testcache
   ```

2. **Run Functional Tests**:
   Execute the Go test suite in `tools/validate`.
   ```bash
   cd tools/validate
   go test -v ./...
   ```

3. **Interpret Results**:
   - **PASS**: All scenarios (Scaffold, Sync, Exclusions, etc.) are working.
   - **FAIL**: Check the logs. If a specific tool fails, fix the logic in `tools/engine`.

## When to use
- After creating a new feature.
- Before verifying a task complete.
- When the user asks to "test everything".
