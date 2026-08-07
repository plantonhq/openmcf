---
title: "Development Project"
description: "This preset creates a lightweight project for development environments: billing linked, a minimal API set, and the default DELETE deletion policy so teardown is one command."
type: "preset"
rank: "02"
presetSlug: "02-development"
componentSlug: "project-on-google-cloud"
componentTitle: "Project on Google Cloud"
provider: "gcp"
icon: "package"
order: 2
---

# Development Project

This preset creates a lightweight project for development environments:
billing linked, a minimal API set, and the default DELETE deletion policy
so teardown is one command.

## When to Use

- Per-team or per-feature development projects
- Sandboxes that are created and destroyed regularly

## Key Configuration Choices

- **Default deletion policy (DELETE)** — destroy shuts the project down
  (30-day recovery window applies).
- **Unique IDs are yours to choose** — project IDs are globally unique and
  reserved for ~30 days after deletion; bake uniqueness into the ID
  itself (team, ticket, or sequence suffix).
- **Minimal API set** — components enable the APIs they need on their own.

## Placeholders to Replace

The sample values for `projectId`, `parentId`, and `billingAccountId` are
realistic placeholders for pattern-validated fields — replace them with
your own.

## Related Presets

- **01-standard-production** — hardened foundation project with destroy
  protection
