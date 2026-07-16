# envpartition: exported accessor for the untaught default rule document

**Date:** 2026-07-14
**Scope:** `pkg/iac/envpartition` (new `DefaultRuleDocument` accessor + test)

## Summary

`pkg/iac/envpartition` gains **`DefaultRuleDocument()`**: a fresh copy of the
embedded untaught `EnvironmentPartitionRule` document — the same one
`DefaultRule()` compiles.

## Why

Services that host the partition engine need to tell callers WHICH rule a
partition actually applied. When a caller sends no rule, the engine applies
its embedded untaught default; without an exported document accessor, any
host wanting to echo that rule back (or any review surface wanting to seed a
teaching UI from it) would have to maintain its own copy of the default
vocabulary — a second definition that drifts. `DefaultRuleDocument()` keeps
the untaught default at exactly one definition.

## Contract

- Returns the document form (`rulev1.EnvironmentPartitionRule`), not the
  compiled rule; `CompileRule(DefaultRuleDocument())` equals `DefaultRule()`
  (pinned by `TestDefaultRuleDocumentMatchesDefaultRule`).
- Returns a fresh copy per call, so a caller mutating its copy can never
  poison another caller's.
- Same panic contract as `DefaultRule()`: the embedded document is part of
  the build, so failing to parse it is a programming error.
