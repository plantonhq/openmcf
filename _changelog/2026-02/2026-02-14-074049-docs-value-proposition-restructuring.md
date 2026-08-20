# Documentation Restructuring: Value-Proposition Oriented Navigation

**Date**: February 14, 2026
**Type**: Refactoring
**Components**: Documentation

## Summary

Restructured the entire Planton documentation from internal domain model organization (Infra Hub, Service Hub, Config Manager, Cloud Ops, Connect) to user-oriented vocabulary (Infrastructure, CI/CD, Secrets, Variables, Operations, Connections). Split "Secrets and Config" into two standalone first-class sections. Updated all cross-references across 55 pages.

## Problem Statement / Motivation

The documentation was organized around Planton's internal domain modeling — product names like "Infra Hub," "Service Hub," "Secrets and Config," and "Cloud Ops" that a DevOps engineer would never search for. When a reader arrives at a documentation site, they bring their existing vocabulary: "infrastructure," "CI/CD," "secrets," "operations." The old structure required readers to learn Planton's internal terminology before they could navigate the docs.

Additionally, "Secrets and Config" bundled two distinct concerns (secrets management and configuration variables) under a single section named after an internal domain ("Config Manager"), diminishing the value proposition of Planton's built-in secrets manager as a first-class, Vault-competitive product.

### Pain Points

- DevOps engineers searching for "infrastructure" found "Infra Hub" instead
- "Service Hub" communicated nothing to someone looking for "CI/CD"
- Secrets management was buried under a combined "Secrets and Config" section
- "Cloud Ops" was a Planton product name with no external recognition
- The documentation structure mirrored implementation domains, not user tasks

## Solution / What's New

### Directory Restructuring (6 sections renamed)

| Old | New | Rationale |
|-----|-----|-----------|
| `connect/` | `connections/` | Matches web console sidebar. Noun, not verb. |
| `infra-hub/` | `infrastructure/` | Universal term every DevOps engineer knows. |
| `service-hub/` | `ci-cd/` | Most searchable term. Tekton extensibility makes it genuine CI/CD. |
| `secrets-and-config/` | `secrets/` | Split: secrets become standalone first-class section. |
| *(new)* | `variables/` | Split: variables get their own section. |
| `cloud-ops/` | `operations/` | Simpler, more standard. |

### Secrets and Config Split

The old 4-page "Secrets and Config" section became:

**Secrets** (4 pages): `index.md`, `managing-secrets.md`, `backends.md`, `encryption.md`
- Backends and encryption split into separate pages for focused treatment
- Index page rewritten with secrets-manager-as-product positioning

**Variables** (2 pages): `index.md`, `variable-groups.md`
- Variable groups extracted to their own page

### Sidebar Result

```
Platform, Connections, Infrastructure, CI/CD, Secrets, Variables,
Operations, Runner, Security, Teams and Access
```

Every item is a clean, concrete noun. No item has an unnecessary qualifier.

### Vocabulary Approach

Section names use user-oriented vocabulary. Inside pages, Planton product names are introduced in context: "In the web console, navigate to **Infra Hub** in the sidebar."

## Implementation Details

- **55 documentation pages** across 10 sections (up from 53 across 9)
- **~300 cross-reference links** updated across all pages
- **6 section index pages** updated with new titles and frontmatter
- **2 new pages** created from content splits (backends.md, encryption.md, variable-groups.md, variables/index.md — net +2 pages)
- **0 broken links** — verified via full site build
- **Pagefind search index** rebuilt (57 pages indexed, 2,214 words)
- Build passes cleanly

## Benefits

- Readers find documentation sections using vocabulary they already know
- Secrets management is promoted to a first-class, standalone product section
- Encryption and backends each get focused treatment
- The sidebar reads as a clean, scannable list of concrete nouns
- Web console product names are still introduced where relevant

## Impact

- All existing internal links updated — zero broken references
- No external links exist yet (site is pre-launch) — no redirect infrastructure needed
- Pagefind search index automatically includes new structure

## Related Work

- Session 1-20: Documentation overhaul project (content creation and quality)
- This session: Structural restructuring (navigation and vocabulary)

---

**Status**: Live
**Timeline**: Single session
