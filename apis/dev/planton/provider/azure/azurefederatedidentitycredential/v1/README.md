# AzureFederatedIdentityCredential

## Overview

`AzureFederatedIdentityCredential` provisions a federated identity credential:
a keyless OIDC trust rule on a user-assigned managed identity that lets an
external workload -- a GitHub Actions workflow, a Kubernetes service account,
any OIDC-issuing system -- exchange its own token for the identity's Azure
credentials. No client secret is created, stored, rotated, or leaked.

## Why a First-Class Resource?

Federated credentials are real infrastructure with their own lifecycle:

- **Many per identity** -- each credential trusts exactly one external
  workload (one issuer + subject + audience triple); an identity accumulates
  one per branch, environment, or service account that acts as it (Azure
  allows up to 20)
- **Independent lifecycle** -- trust rules come and go as pipelines and
  workloads change, without touching the identity or its grants
- **The keyless unlock** -- this is the resource that makes secretless CI/CD
  and AKS workload identity possible; modeling it first-class makes "which
  external systems can act as this identity" reviewable infrastructure

## Key Features

- **Three-way trust match** -- Azure AD exchanges a token only when its
  `iss`, `sub`, and `aud` claims equal the credential's issuer, subject, and
  audience exactly
- **GitHub Actions ready** -- trust a repo's branch, environment, or tag with
  the documented subject formats
- **AKS workload identity ready** -- trust a cluster service account against
  the cluster's OIDC issuer URL, referenced directly from an
  `AzureAksCluster`'s `oidc_issuer_url` output so the trust always matches
  the cluster it is deployed beside
- **Sensible audience default** -- `api://AzureADTokenExchange` (what every
  standard client requests) unless explicitly overridden
- **Composable** -- the parent identity is referenced by ARM ID, defaulting
  to an `AzureUserAssignedIdentity`'s output in composed environments

## When to Use

- Granting a GitHub Actions workflow keyless access to deploy Azure
  infrastructure (no stored service-principal secret)
- Enabling AKS workload identity: pods authenticate as a managed identity
  through their service account's projected token
- Trusting any external OIDC issuer (GitLab, self-hosted systems) to act as
  an Azure identity

## Spec Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | The credential's name under the parent identity (3-120 chars, unique per identity) |
| `user_assigned_identity` | StringValueOrRef | Yes | ARM ID of the parent identity (defaults to an AzureUserAssignedIdentity reference) |
| `issuer` | StringValueOrRef | Yes | OIDC issuer URL the token's `iss` claim must equal -- a literal URL, or a reference defaulting to an AzureAksCluster's `oidc_issuer_url` output |
| `subject` | string | Yes | Workload identifier the token's `sub` claim must equal |
| `audience` | string | No | The token's required `aud` claim; defaults to `api://AzureADTokenExchange` |

Federated credentials do not support ARM tags (they are untagged child
resources of the identity), so the spec carries none.

## Outputs

| Output | Description |
|--------|-------------|
| `federated_identity_credential_id` | Full ARM ID of the credential |
| `name` | The credential's name as deployed |
| `user_assigned_identity_id` | ARM ID of the parent identity |
| `issuer` | The trusted issuer as deployed |
| `subject` | The trusted subject as deployed |
| `audience` | The required audience as deployed |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFederatedIdentityCredential
metadata:
  name: ci-deployer-trust
  org: mycompany
  env: production
spec:
  name: github-main-branch
  userAssignedIdentity:
    valueFrom:
      name: ci-deployer-identity
  issuer:
    value: https://token.actions.githubusercontent.com
  subject: repo:mycompany/platform:ref:refs/heads/main
```

The complete keyless-CI story composes three kinds: the identity
(`AzureUserAssignedIdentity`), its permissions (`AzureRoleAssignment`
granting, say, Contributor on a resource group), and this trust rule letting
the workflow act as the identity.

## Required Permissions

The deploying credential needs
`Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials/write`
on the parent identity -- held via Managed Identity Contributor, Contributor,
or Owner on the identity's scope.

## Subject Format Reference

| External workload | Subject format |
|-------------------|----------------|
| GitHub branch | `repo:{owner}/{repo}:ref:refs/heads/{branch}` |
| GitHub environment | `repo:{owner}/{repo}:environment:{env}` |
| GitHub tag | `repo:{owner}/{repo}:ref:refs/tags/{tag}` |
| GitHub pull requests | `repo:{owner}/{repo}:pull_request` |
| Kubernetes service account | `system:serviceaccount:{namespace}:{serviceaccount}` |

No wildcards -- one credential trusts one subject. Create one credential per
branch, environment, or service account.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
