---
title: "Conditional Grant (Time-Bound Access)"
description: "This preset grants a role that only applies while a CEL condition evaluates true — here, an expiry timestamp that makes the access self-revoking. Conditions also support resource-name prefixes and..."
type: "preset"
rank: "03"
presetSlug: "03-conditional-grant"
componentSlug: "project-iam-member"
componentTitle: "Project IAM Member"
provider: "gcp"
icon: "package"
order: 3
---

# Conditional Grant (Time-Bound Access)

This preset grants a role that only applies while a CEL condition evaluates true — here, an expiry timestamp that makes the access self-revoking. Conditions also support resource-name prefixes and other request attributes for scoping a grant below the project level.

## When to Use

- Temporary human access for a migration, incident, or review that must auto-expire
- Scoping a project-level grant to a subset of resources by name prefix (e.g. only `prod-` buckets)
- Any grant a security policy requires to be time-bound

## Key Configuration Choices

- **Expiry via `request.time`** — the grant stops applying at the timestamp with no follow-up action; remove the resource at leisure
- **Condition is part of the grant's identity** — the same role granted unconditionally is a separate, independent grant; adding a condition never narrows an existing unconditional grant
- **`user:` member** — conditions most often gate human access; service accounts work identically

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<gcp-project-id>` | Project whose IAM policy receives the grant | GCP Console or `GcpProject` outputs |
| `<user-email>` | The human user receiving temporary access | Your identity provider / Google Workspace |
| `<expiry-timestamp>` | RFC3339 expiry, e.g. `2027-01-01T00:00:00Z` | Your access policy |

## Related Presets

- **01-service-account-grant** — The standard unconditional workload grant
- **02-custom-role-grant** — Grant a custom role defined as a first-class node
