---
asdp_version: 0.1.15
root: true
components:
    - name: core
      type: module
      path: ./core
      description: (No specification found)
      last_modified: 2026-01-22T01:27:09.635343033+01:00
      has_spec: false
      has_model: true
      children:
        - name: agent
          type: module
          path: ./core/agent
          description: (No specification found)
          last_modified: 2026-01-22T01:04:35.734932391+01:00
          has_spec: false
          has_model: false
          children:
            - name: rules
              type: module
              path: ./core/agent/rules
              description: (No specification found)
              last_modified: 2026-01-22T01:04:35.734932391+01:00
              has_spec: false
              has_model: false
            - name: skills
              type: module
              path: ./core/agent/skills
              description: (No specification found)
              last_modified: 2026-01-21T16:40:29.034512609+01:00
              has_spec: false
              has_model: false
              children:
                - name: asdp-doc-rebuilder
                  type: module
                  path: ./core/agent/skills/asdp-doc-rebuilder
                  description: (No specification found)
                  last_modified: 2026-01-21T16:40:59.070404591+01:00
                  has_spec: false
                  has_model: false
            - name: workflows
              type: module
              path: ./core/agent/workflows
              description: (No specification found)
              last_modified: 2026-01-22T01:03:19.968503609+01:00
              has_spec: false
              has_model: false
        - name: spec
          type: library
          path: ./core/spec
          description: Provides a resilient HTTP client with auto-retry and logging middleware.
          last_modified: 2026-01-22T01:27:09.635343033+01:00
          has_spec: true
          has_model: true
    - name: improvements
      type: module
      path: ./improvements
      description: (No specification found)
      last_modified: 2026-01-22T00:26:17.971158099+01:00
      has_spec: false
      has_model: false
    - name: installer
      type: module
      path: ./installer
      description: (No specification found)
      last_modified: 2026-01-21T18:16:41.528092019+01:00
      has_spec: false
      has_model: false
    - name: tools
      type: library
      path: ./tools
      description: Container for all ASDP tooling and executables.
      last_modified: 2026-01-22T01:27:09.636343012+01:00
      has_spec: true
      has_model: true
      children:
        - name: engine
          type: library
          path: ./tools/engine
          description: The implementation core of ASDP logic.
          last_modified: 2026-01-22T01:27:09.636343012+01:00
          has_spec: true
          has_model: true
          children:
            - name: domain
              type: library
              path: ./tools/engine/domain
              description: Defines the core domain entities and interfaces for ASDP.
              last_modified: 2026-01-22T01:27:09.637342992+01:00
              has_spec: true
              has_model: true
            - name: system
              type: library
              path: ./tools/engine/system
              description: Provides infrastructure implementations for ASDP interfaces.
              last_modified: 2026-01-22T01:27:09.638342972+01:00
              has_spec: true
              has_model: true
            - name: usecase
              type: library
              path: ./tools/engine/usecase
              description: Implements the Application Business Rules.
              last_modified: 2026-01-22T01:27:09.640342931+01:00
              has_spec: true
              has_model: true
        - name: mcp-server
          type: application
          path: ./tools/mcp-server
          description: The Model Context Protocol Server implementation.
          last_modified: 2026-01-22T01:27:09.64634281+01:00
          has_spec: true
          has_model: true
          children:
            - name: cmd
              type: module
              path: ./tools/mcp-server/cmd
              description: Container for executable binaries.
              last_modified: 2026-01-22T01:27:09.641342911+01:00
              has_spec: true
              has_model: true
              children:
                - name: asdp-mcp-server
                  type: application
                  path: ./tools/mcp-server/cmd/asdp-mcp-server
                  description: The main entrypoint for the ASDP MCP Server application.
                  last_modified: 2026-01-22T01:27:09.641342911+01:00
                  has_spec: true
                  has_model: true
            - name: internal
              type: library
              path: ./tools/mcp-server/internal
              description: Internal implementation details for the MCP server.
              last_modified: 2026-01-22T01:27:09.647342789+01:00
              has_spec: true
              has_model: true
              children:
                - name: adapter
                  type: library
                  path: ./tools/mcp-server/internal/adapter
                  description: Interface Adapters for the MCP server.
                  last_modified: 2026-01-22T01:27:09.64634281+01:00
                  has_spec: true
                  has_model: true
                  children:
                    - name: mcp
                      type: library
                      path: ./tools/mcp-server/internal/adapter/mcp
                      description: Adapts ASDP UseCases to the Model Context Protocol (MCP).
                      last_modified: 2026-01-22T01:27:09.647342789+01:00
                      has_spec: true
                      has_model: true
verification:
    scan_time: 2026-01-22T01:27:09.649729093+01:00
---

# Project Hierarchy

Auto-generated by ASDP SyncTree.
