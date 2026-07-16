---
title: "CI Federation Pool"
description: "This preset creates the standard keyless-auth trust boundary for CI/CD: a FEDERATION_ONLY pool that GitHub Actions, GitLab CI, or any other OIDC issuer attaches to via a..."
type: "preset"
rank: "01"
presetSlug: "01-ci-federation-pool"
componentSlug: "workload-identity-pool"
componentTitle: "Workload Identity Pool"
provider: "gcp"
icon: "package"
order: 1
---

# CI Federation Pool

This preset creates the standard keyless-auth trust boundary for CI/CD: a FEDERATION_ONLY pool that GitHub Actions, GitLab CI, or any other OIDC issuer attaches to via a GcpWorkloadIdentityPoolProvider. No service-account keys are created anywhere in this path.

## When to Use

- CI pipelines currently authenticate with downloaded service-account JSON keys
- You are setting up deployment automation for a new project and want keyless from day one
- Security review flagged long-lived credentials in CI secrets

## Key Configuration Choices

- **FEDERATION_ONLY mode (default, omitted)** — the token-exchange mode; the only mode that can hold providers
- **A stable, boundary-named pool ID** — the pool name is embedded in every IAM principal, so name it after the trust boundary (`ci-federation`), not after one issuer; one pool can trust several issuers through separate providers
- **No `disabled` flag** — the pool starts active; `disabled: true` is the reversible kill switch to keep in your incident-response runbook (deleting a pool reserves its ID for ~30 days and cannot be undone by re-creating)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<gcp-project-id>` | GCP project that owns the pool | GCP Console or `GcpProject` outputs |
| `ci-federation` (workloadIdentityPoolId) | Stable pool ID (4-32 chars; lowercase, digits, hyphens) | Name it after the trust boundary |

## Related Presets

- **02-locked-down-pool** — The same pool pre-disabled, for staged rollouts
