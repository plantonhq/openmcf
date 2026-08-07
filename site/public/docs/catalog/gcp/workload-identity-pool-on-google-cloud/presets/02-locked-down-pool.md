---
title: "Locked-Down (Staged) Pool"
description: "This preset provisions the pool in a disabled state: the resource exists, its name is stable, IAM bindings and providers can be prepared against it — but every token exchange is rejected until..."
type: "preset"
rank: "02"
presetSlug: "02-locked-down-pool"
componentSlug: "workload-identity-pool-on-google-cloud"
componentTitle: "Workload Identity Pool on Google Cloud"
provider: "gcp"
icon: "package"
order: 2
---

# Locked-Down (Staged) Pool

This preset provisions the pool in a disabled state: the resource exists, its name is stable, IAM bindings and providers can be prepared against it — but every token exchange is rejected until `disabled` flips to `false`. Use it to land federation infrastructure ahead of a go-live, or to template the incident-response posture.

## When to Use

- Rolling out federation in stages: infrastructure first, trust enabled at cutover
- Provisioning a boundary for a partner/vendor whose access starts on a contract date
- Codifying the kill-switch state so re-disabling in an incident is a reviewed one-line diff

## Key Configuration Choices

- **`disabled: true`** — the pool's kill switch, fully reversible. This is the correct staging/lockdown lever: deleting a pool instead reserves its ID for ~30 days and cannot be undone by re-creating.
- **Grants can be created while disabled** — IAM bindings on `principal://` members from this pool are inert until the pool is enabled, so the whole access topology can be reviewed before any token is honored.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<gcp-project-id>` | GCP project that owns the pool | GCP Console or `GcpProject` outputs |
| `partner-federation` (workloadIdentityPoolId) | Stable pool ID (4-32 chars; lowercase, digits, hyphens) | Name it after the trust boundary |

## Related Presets

- **01-ci-federation-pool** — The standard active federation pool for CI/CD
