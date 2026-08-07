---
title: "Token Creator Grant (Cross-Account Impersonation)"
description: "This preset grants `roles/iam.serviceAccountTokenCreator` on one service account to another — letting the caller mint short-lived access and ID tokens AS the target. This is the building block of..."
type: "preset"
rank: "02"
presetSlug: "02-token-creator-grant"
componentSlug: "service-account-iam-member-on-google-cloud"
componentTitle: "Service Account IAM Member on Google Cloud"
provider: "gcp"
icon: "package"
order: 2
---

# Token Creator Grant (Cross-Account Impersonation)

This preset grants `roles/iam.serviceAccountTokenCreator` on one service account to another — letting the caller mint short-lived access and ID tokens AS the target. This is the building block of token-broker designs and privilege elevation on demand: the caller holds modest standing permissions and impersonates a more privileged account only for specific operations.

## When to Use

- A workload needs to occasionally act as a more privileged identity without holding its permissions permanently
- Generating signed ID tokens for service-to-service authentication (e.g. invoking an authenticated Cloud Run service)
- Centralizing privileged access behind one auditable impersonation edge

## Key Configuration Choices

- **Both sides referenced** — the target's `name` output and the caller's `member` output both come from GcpServiceAccount nodes, so the impersonation edge is fully visible in the resource graph
- **Account-scoped, not project-scoped** — a project-level tokenCreator grant would let the caller mint tokens as EVERY account in the project
- **Tokens are short-lived** — minted credentials expire (default 1 hour); no long-lived key material is ever created

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<target-service-account-resource-name>` | The account whose tokens may be minted | Your GcpServiceAccount manifest's `metadata.name` |
| `<caller-service-account-resource-name>` | The account allowed to mint them | Your GcpServiceAccount manifest's `metadata.name` |

## Related Presets

- **01-github-workload-identity-impersonation** — Federated (keyless) impersonation from GitHub
- **03-deployer-act-as** — The actAs permission Cloud Run/GCE deployments require
