---
asdp_version: 1.0.0
integrity:
    src_hash: e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
    algorithm: sha256
    last_modified: 2026-01-22T01:05:59.042280577+01:00
    checked_at: 2026-01-22T01:27:09.63637224+01:00
symbols: []
---

```

### 2. Markdown Body (Agent Annotations)

The body allows Agents to add *semantic understanding* to the raw symbols. While the YAML is the "What", the Markdown is the "How it actually works inside".

```markdown
# Semantic Model

## Client
The `Client` struct holds the state for...

## NewClient
Initializes the transport layer. Note that it sets `MaxIdleConns` to 100 by default.
```
