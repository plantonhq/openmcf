---
title: "Pure Identity with First-Class Grants"
description: "This preset creates the service account as a pure identity node with no inline role lists. Every grant lives as its own GcpProjectIamMember resource referencing this account's `member` output — the..."
type: "preset"
rank: "03"
presetSlug: "03-identity-with-first-class-grants"
componentSlug: "service-account"
componentTitle: "Service Account"
provider: "gcp"
icon: "package"
order: 3
---

# Pure Identity with First-Class Grants

This preset creates the service account as a pure identity node with no inline role lists. Every grant lives as its own GcpProjectIamMember resource referencing this account's `member` output — the fully-composed pattern where identity, roles, and access are independent nodes in the resource graph.

## When to Use

- Grants use custom roles (GcpIamCustomRole) or IAM conditions, which the inline role lists cannot express
- Different charts or teams own the identity and its grants (e.g. a platform chart creates the identity; an app chart grants it access)
- You want access changes reviewable and revertible independently of the identity's lifecycle

## Key Configuration Choices

- **No `projectIamRoles`/`orgIamRoles`** — deliberate; the grant nodes carry all access
- **Keyless** — `createKey` stays at its default (false); pair with Workload Identity or federation
- **`description` names the pattern** — the auditor who finds this identity sees immediately where its access is defined

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `my-app-prod` (serviceAccountId) | Service account ID (6-30 chars, lowercase) | Replace with a descriptive ID for your workload |
| `<gcp-project-id>` | GCP project ID | `GcpProject` outputs |

## Related Presets

- **01-workload-identity** — Identity with the common inline role list for GKE workloads
- **02-ci-cd-pipeline** — CI/CD identity with an exported key for non-federating systems
