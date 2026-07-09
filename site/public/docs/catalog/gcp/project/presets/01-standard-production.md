---
title: "Standard Production Project"
description: "This preset creates a production-grade project under a folder: billing linked, the default network suppressed, a hardening baseline of APIs pre-enabled, and destroy blocked by `deletionPolicy:..."
type: "preset"
rank: "01"
presetSlug: "01-standard-production"
componentSlug: "project"
componentTitle: "Project"
provider: "gcp"
icon: "package"
order: 1
---

# Standard Production Project

This preset creates a production-grade project under a folder: billing
linked, the default network suppressed, a hardening baseline of APIs
pre-enabled, and destroy blocked by `deletionPolicy: PREVENT`.

## When to Use

- Foundation projects for production workloads
- Any project whose accidental deletion would be catastrophic

## Key Configuration Choices

- **`deletionPolicy: PREVENT`** — destroy fails while set; flip to DELETE
  deliberately when decommissioning.
- **No auto-created default network** (spec default) — networks are
  explicit `GcpVpcNetwork` resources.
- **IAM grants are not part of the project** — add
  `GcpProjectIamMember` resources per grant.

## Placeholders to Replace

The sample values for `projectId`, `parentId`, and `billingAccountId` are
realistic placeholders for pattern-validated fields — replace them with
your own project ID (6-30 lowercase chars), numeric folder ID, and billing
account ID.

## Related Presets

- **02-development** — lightweight project for dev environments
