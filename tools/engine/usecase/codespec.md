---
asdp_version: "0.1.15"
id: "usecase"
type: "library"
title: "usecase"
summary: "Implements the Application Business Rules."
capabilities:
  - "InitProjectUseCase"
  - "InitAgentUseCase"
  - "SyncModelUseCase"
  - "SyncTreeUseCase"
  - "ScaffoldUseCase"
dependencies:
  - module: "domain"
    reason: "Business logic"
  - module: "system"
    reason: "Injected infrastructure"
requirements: []
exports:
  - "NewInitProjectUseCase"
  - "NewSyncTreeUseCase"
---
# UseCase Specification

## Context

This module acts as the Application Layer (Clean Architecture). It orchestrates the flow of data to and from the `domain` entities, and directs those entities to use their critical business rules to achieve the goals of the use case. It depends on `domain` interfaces and is injected with `system` implementations.

## Requirements

- Must implement specific user operations (e.g. Sync, Init, Scaffold).
- Must remain independent of UI/Delivery mechanisms (MCP).
