# ASDP — Agentic Spec Driven Programming

Reference repo for the ASDP protocol: a documentation-first workflow for AI agents. It defines how to structure projects with canonical specs so architecture stays aligned with code.

## What this is

- Minimal spec set for agents to plan, edit, and sync code via documentation contracts.
- Templates that describe the required shape of CodeTree, CodeSpec, and CodeModel files.
- Rulebooks that enforce read-before-work, sync guarantees, and completion checks.

## Repository layout

- `spec/asdp/` — format definitions for `codetree.md`, `codespec.md`, and `codemodel.md`.
- `.agent/rules/` — process rules for planning/editing, doc synchronization, and post-work duties.
- Root docs — `codetree.md`, `codespec.md`, `codemodel.md` describing this repo’s own structure.

## How to use

1) Read `rules/before-planning-or-editing.rules.md`, then root `codetree.md`, `codespec.md`, `codemodel.md`.
2) When working in a folder, also read its `codespec.md` and `codemodel.md`.
3) Apply changes; ensure code matches the docs (docs are canonical).
4) Sync docs per `rules/after-editing-or-task-complete.rules.md` and update changelog entries.

## When to extend

- New folders: add to `codetree.md` and create matching `codespec.md`/`codemodel.md`.
- New behaviors or APIs: update the relevant `codespec.md` and `codemodel.md` in that folder and parents as needed.
