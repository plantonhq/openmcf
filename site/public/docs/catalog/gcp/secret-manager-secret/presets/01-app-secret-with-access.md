---
title: "App Secret with Access"
description: "The standard workload-credential shape: one manifest creates the secret, stores the payload as version 1, and grants the consuming service account read access — a READABLE secret with zero follow-up..."
type: "preset"
rank: "01"
presetSlug: "01-app-secret-with-access"
componentSlug: "secret-manager-secret"
componentTitle: "Secret Manager Secret"
provider: "gcp"
icon: "package"
order: 1
---

# App Secret with Access

The standard workload-credential shape: one manifest creates the
secret, stores the payload as version 1, and grants the consuming
service account read access — a READABLE secret with zero follow-up
steps.

## What it configures

- Automatic replication (no `replication` block) — Google places the
  payload; the right default when no residency regime applies.
- `initialVersion.data` — the payload, handled as a managed secret
  (never plaintext at rest in the control plane).
- A secret-SCOPED `secretAccessor` grant — the workload reads THIS
  secret, not every secret in the project.

## Adjust before deploying

- **data** — supply through the platform's secret handling; in charts,
  prefer valueFrom from a producing resource's sensitive output.
- **member** — reference a GcpServiceAccount resource via valueFrom
  (its `member` output is exactly this value).

## When to choose something else

Data-residency regimes need the **Regional CMEK Secret** preset;
credentials with a rotation pipeline behind them start from the
**Rotated Secret** preset.
