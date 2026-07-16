---
title: "Deployer actAs Grant"
description: "This preset grants `roles/iam.serviceAccountUser` on a runtime service account to a deployer identity. This is the actAs permission that Cloud Run, GCE, Cloud Functions, and Dataflow all check at..."
type: "preset"
rank: "03"
presetSlug: "03-deployer-act-as"
componentSlug: "service-account-iam-member"
componentTitle: "Service Account IAM Member"
provider: "gcp"
icon: "package"
order: 3
---

# Deployer actAs Grant

This preset grants `roles/iam.serviceAccountUser` on a runtime service account to a deployer identity. This is the actAs permission that Cloud Run, GCE, Cloud Functions, and Dataflow all check at deploy time: whoever attaches a service account to a workload must be allowed to act as it. Without this grant, deploys fail with a `iam.serviceaccounts.actAs` permission error even when the deployer has full admin on the workload service itself.

## When to Use

- A CI pipeline deploys Cloud Run services (or GCE instances, functions, jobs) that run as a dedicated runtime account
- A human operator needs to deploy workloads attached to a specific service account
- You are tightening a project-level serviceAccountUser grant down to specific runtime accounts

## Key Configuration Choices

- **Grant per runtime account** — one grant per (deployer, runtime SA) pair keeps the blast radius explicit; a project-level serviceAccountUser grant would allow acting as EVERY account in the project
- **Deployer as a reference** — when the deployer is a Planton-managed GcpServiceAccount, referencing its `member` output keeps the edge visible; for human deployers use a literal `user:<email>` member instead
- **Pairs with workload kinds** — the workload manifest (Cloud Run, GCE) references the runtime account by email; this grant is the permission that lets the deploy succeed

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<runtime-service-account-resource-name>` | The account workloads run as | Your GcpServiceAccount manifest's `metadata.name` |
| `<deployer-service-account-resource-name>` | The deploying identity | Your GcpServiceAccount manifest's `metadata.name` |

## Related Presets

- **01-github-workload-identity-impersonation** — Keyless CI/CD impersonation from GitHub
- **02-token-creator-grant** — Short-lived token minting between accounts
