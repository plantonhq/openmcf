---
title: "GitLab CI OIDC Provider"
description: "This preset attaches GitLab CI as a trusted issuer: pipelines exchange their GitLab-minted ID tokens for short-lived Google credentials. Works for gitlab.com and (with the issuer URI changed)..."
type: "preset"
rank: "03"
presetSlug: "03-gitlab-ci-oidc"
componentSlug: "workload-identity-pool-provider-on-google-cloud"
componentTitle: "Workload Identity Pool Provider on Google Cloud"
provider: "gcp"
icon: "package"
order: 3
---

# GitLab CI OIDC Provider

This preset attaches GitLab CI as a trusted issuer: pipelines exchange their GitLab-minted ID tokens for short-lived Google credentials. Works for gitlab.com and (with the issuer URI changed) self-managed GitLab instances.

## When to Use

- GitLab CI pipelines deploy to or read from Google Cloud
- You are removing service-account JSON keys from GitLab CI/CD variables
- You want per-project or per-branch control over what CI can touch in GCP

## Key Configuration Choices

- **`attributeCondition` scoped to your group — mandatory on gitlab.com.** The shared issuer signs tokens for every project on the instance; the condition restricts federation to your namespace. Tighten to a project (`assertion.project_path == "group/project"`) or protected branches (`assertion.ref_protected == "true"`) as needed.
- **`attribute.project_path` and `attribute.ref` mappings** — enable per-project and per-branch grants via `principalSet://.../attribute.project_path/<group>/<project>`.
- **Self-managed instances** — set `issuerUri` to your instance URL; if the instance is not internet-reachable, provide its signing keys inline via `oidc.jwksJson`.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `ci-federation-pool` | Your GcpWorkloadIdentityPool resource name | The pool manifest's `metadata.name` |
| `gitlab-oidc` (providerId) | Stable provider ID (4-32 chars) | Name it after the issuer |
| `<gitlab-group>` | Your GitLab group/namespace path | GitLab group page |

## Related Presets

- **01-github-actions-oidc** — Trust GitHub Actions workflows
- **02-aws-account** — Trust workloads in one AWS account
