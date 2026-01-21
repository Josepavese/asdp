---
asdp_version: 0.1.13
root: true
components:
    - name: core
      type: module
      path: ./core
      description: ""
      last_modified: 2026-01-21T16:40:29.034512609+01:00
      has_spec: false
      has_model: false
      children:
        - name: agent
          type: module
          path: ./core/agent
          description: ""
          last_modified: 2026-01-22T00:11:15.457518557+01:00
          has_spec: false
          has_model: false
          children:
            - name: rules
              type: module
              path: ./core/agent/rules
              description: ""
              last_modified: 2026-01-21T19:55:21.081941349+01:00
              has_spec: false
              has_model: false
            - name: skills
              type: module
              path: ./core/agent/skills
              description: ""
              last_modified: 2026-01-21T16:40:29.034512609+01:00
              has_spec: false
              has_model: false
              children:
                - name: asdp-doc-rebuilder
                  type: module
                  path: ./core/agent/skills/asdp-doc-rebuilder
                  description: ""
                  last_modified: 2026-01-21T16:40:59.070404591+01:00
                  has_spec: false
                  has_model: false
            - name: workflows
              type: module
              path: ./core/agent/workflows
              description: ""
              last_modified: 2026-01-22T00:11:22.38254873+01:00
              has_spec: false
              has_model: false
        - name: spec
          type: library
          path: ./core/spec
          description: HTTP Client Wrapper
          last_modified: 2026-01-22T00:13:20.517014808+01:00
          has_spec: true
          has_model: true
    - name: installer
      type: module
      path: ./installer
      description: ""
      last_modified: 2026-01-21T18:16:41.528092019+01:00
      has_spec: false
      has_model: false
    - name: sandbox
      type: module
      path: ./sandbox
      description: ""
      last_modified: 2026-01-22T00:09:18.279951949+01:00
      has_spec: false
      has_model: false
      children:
        - name: home
          type: module
          path: ./sandbox/home
          description: ""
          last_modified: 2026-01-21T17:02:37.966217791+01:00
          has_spec: false
          has_model: false
          children:
            - name: go
              type: module
              path: ./sandbox/home/go
              description: ""
              last_modified: 2026-01-21T17:01:37.808354027+01:00
              has_spec: false
              has_model: false
              children:
                - name: pkg
                  type: module
                  path: ./sandbox/home/go/pkg
                  description: ""
                  last_modified: 2026-01-21T17:02:37.967217789+01:00
                  has_spec: false
                  has_model: false
                  children:
                    - name: mod
                      type: module
                      path: ./sandbox/home/go/pkg/mod
                      description: ""
                      last_modified: 2026-01-21T17:02:37.967217789+01:00
                      has_spec: false
                      has_model: false
                      children:
                        - name: gopkg.in
                          type: module
                          path: ./sandbox/home/go/pkg/mod/gopkg.in
                          description: ""
                          last_modified: 2026-01-21T17:01:38.625352046+01:00
                          has_spec: false
                          has_model: false
                          children:
                            - name: check.v1@v0.0.0-20161208181325-20d25e280405
                              type: module
                              path: ./sandbox/home/go/pkg/mod/gopkg.in/check.v1@v0.0.0-20161208181325-20d25e280405
                              description: ""
                              last_modified: 2026-01-21T17:01:38.625352046+01:00
                              has_spec: false
                              has_model: false
                            - name: yaml.v3@v3.0.1
                              type: module
                              path: ./sandbox/home/go/pkg/mod/gopkg.in/yaml.v3@v3.0.1
                              description: ""
                              last_modified: 2026-01-21T17:01:38.340352736+01:00
                              has_spec: false
                              has_model: false
        - name: project_v2
          type: module
          path: ./sandbox/project_v2
          description: ""
          last_modified: 2026-01-22T00:09:18.280951955+01:00
          has_spec: false
          has_model: false
          children:
            - name: src
              type: module
              path: ./sandbox/project_v2/src
              description: src
              last_modified: 2026-01-22T00:09:18.280951955+01:00
              has_spec: true
              has_model: true
    - name: tools
      type: module
      path: ./tools
      description: ""
      last_modified: 2026-01-22T00:13:28.363042779+01:00
      has_spec: false
      has_model: false
      children:
        - name: engine
          type: module
          path: ./tools/engine
          description: ""
          last_modified: 2026-01-22T00:08:04.42953289+01:00
          has_spec: false
          has_model: false
          children:
            - name: domain
              type: module
              path: ./tools/engine/domain
              description: ""
              last_modified: 2026-01-22T00:13:30.641050836+01:00
              has_spec: false
              has_model: true
            - name: system
              type: module
              path: ./tools/engine/system
              description: ""
              last_modified: 2026-01-21T23:59:15.802222674+01:00
              has_spec: false
              has_model: true
            - name: usecase
              type: module
              path: ./tools/engine/usecase
              description: ""
              last_modified: 2026-01-22T00:13:08.279970493+01:00
              has_spec: false
              has_model: true
        - name: mcp-server
          type: application
          path: ./tools/mcp-server
          description: ASDP MCP Server
          last_modified: 2026-01-22T00:13:28.408042939+01:00
          has_spec: true
          has_model: true
          children:
            - name: cmd
              type: module
              path: ./tools/mcp-server/cmd
              description: ""
              last_modified: 2026-01-21T23:53:20.999306413+01:00
              has_spec: false
              has_model: false
              children:
                - name: asdp-mcp-server
                  type: module
                  path: ./tools/mcp-server/cmd/asdp-mcp-server
                  description: ""
                  last_modified: 2026-01-22T00:08:40.653745046+01:00
                  has_spec: false
                  has_model: true
            - name: internal
              type: module
              path: ./tools/mcp-server/internal
              description: ""
              last_modified: 2026-01-21T17:01:54.184315023+01:00
              has_spec: false
              has_model: false
              children:
                - name: adapter
                  type: module
                  path: ./tools/mcp-server/internal/adapter
                  description: ""
                  last_modified: 2026-01-21T23:53:22.482351812+01:00
                  has_spec: false
                  has_model: false
                  children:
                    - name: mcp
                      type: module
                      path: ./tools/mcp-server/internal/adapter/mcp
                      description: ""
                      last_modified: 2026-01-22T00:09:01.226859814+01:00
                      has_spec: false
                      has_model: true
verification:
    scan_time: 2026-01-22T00:13:30.643650484+01:00
---

# Project Hierarchy

Auto-generated by ASDP SyncTree.
