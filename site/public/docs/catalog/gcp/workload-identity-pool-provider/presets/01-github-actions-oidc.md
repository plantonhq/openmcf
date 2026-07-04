---
title: "GitHub Actions OIDC Provider"
description: "This preset attaches GitHub Actions as a trusted issuer: workflows exchange their GitHub-minted OIDC tokens for short-lived Google credentials, replacing service-account keys in CI secrets entirely...."
type: "preset"
rank: "01"
presetSlug: "01-github-actions-oidc"
componentSlug: "workload-identity-pool-provider"
componentTitle: "Workload Identity Pool Provider"
provider: "gcp"
icon: "package"
order: 1
---

# GitHub Actions OIDC Provider

This preset attaches GitHub Actions as a trusted issuer: workflows exchange their GitHub-minted OIDC tokens for short-lived Google credentials, replacing service-account keys in CI secrets entirely. It is the most common federation setup and the reference shape for any OIDC issuer.

## When to Use

- GitHub Actions workflows deploy to or read from Google Cloud
- You are removing service-account JSON keys from repository/organization secrets
- You want per-repository access control over what CI can touch in GCP

## Key Configuration Choices

- **`attributeCondition` scoped to your org — mandatory.** GitHub's issuer signs tokens for every repository on github.com; without the condition, any repository could federate into your pool. Tighten further to specific repositories (`assertion.repository == "my-org/my-repo"`) or branches (`assertion.ref == "refs/heads/main"`) as your posture demands.
- **`attribute.repository` mapping** — enables per-repo grants: bind roles to `principalSet://iam.googleapis.com/<pool name>/attribute.repository/<org>/<repo>` so each repository gets exactly its own access.
- **Default audience (no `allowedAudiences`)** — tokens must be minted for the provider's own canonical resource name (this provider's `name` output), the safest default; configure your workflow's auth step with that audience.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `ci-federation-pool` | Your GcpWorkloadIdentityPool resource name | The pool manifest's `metadata.name` |
| `github-oidc` (providerId) | Stable provider ID (4-32 chars) | Name it after the issuer |
| `<github-org>` | Your GitHub organization login | github.com organization page |

## Related Presets

- **02-aws-account** — Trust workloads in one AWS account
- **03-gitlab-ci-oidc** — Trust GitLab CI pipelines
