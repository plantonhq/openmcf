---
title: "Free Instance"
description: "Provisions the project's one zero-cost Cloud Spanner instance — full Spanner semantics (strong consistency, SQL, the same client libraries) with limited capacity (about 10 GB of storage). Nothing..."
type: "preset"
rank: "01"
presetSlug: "01-free-instance"
componentSlug: "spanner-instance"
componentTitle: "Spanner Instance"
provider: "gcp"
icon: "package"
order: 1
---

# Free Instance

Provisions the project's one zero-cost Cloud Spanner instance — full Spanner semantics (strong consistency, SQL, the same client libraries) with limited capacity (about 10 GB of storage). Nothing about the application code changes when you later move to a provisioned instance.

## When to Use

- Learning Spanner or prototyping a schema before committing to capacity
- Development and integration-test environments where cost must be zero
- Evaluating Spanner's SQL dialects and client libraries against real behavior

## Key Configuration

- **FREE_INSTANCE** — no capacity fields, no edition, no automatic backups (all rejected by GCP for free instances; the spec validates this pre-deploy)
- **One per billing account** — GCP allows a single free instance per billing account
- **One-way upgrade** — a free instance can be upgraded to PROVISIONED in place; there is no downgrade path

## Customization Notes

- `metadata.name` doubles as the instance name when `instanceName` is omitted; instance names must be 6-30 characters
- `project_id` falls back to the provider's default project; set `projectId` (value or `valueFrom` ref to a `GcpProject`) to target another project

## Related Presets

- **02-regional-production** — fixed-capacity production instance
- **03-autoscaling-production** — autoscaling multi-region instance
