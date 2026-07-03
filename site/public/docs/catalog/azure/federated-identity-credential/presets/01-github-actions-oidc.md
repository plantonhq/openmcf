---
title: "GitHub Actions Keyless CI"
description: "This preset trusts a GitHub repository's `main` branch to authenticate as a managed identity -- the standard shape for secretless CI/CD. Workflows request an OIDC token (`permissions: id-token:..."
type: "preset"
rank: "01"
presetSlug: "01-github-actions-oidc"
componentSlug: "federated-identity-credential"
componentTitle: "Federated Identity Credential"
provider: "azure"
icon: "package"
order: 1
---

# GitHub Actions Keyless CI

This preset trusts a GitHub repository's `main` branch to authenticate as a
managed identity -- the standard shape for secretless CI/CD. Workflows request
an OIDC token (`permissions: id-token: write`) and `azure/login` exchanges it
for the identity's credentials; no service-principal secret exists anywhere
to store, rotate, or leak.

The subject pins the trust to `main`-branch runs. A feature-branch or fork
workflow presents a different `sub` claim and is rejected -- the trust is
exactly as narrow as this string.

The workflow side that consumes it:

```yaml
permissions:
  id-token: write
steps:
  - uses: azure/login@v2
    with:
      client-id: <the identity's status.outputs.client_id>
      tenant-id: <the identity's status.outputs.tenant_id>
      subscription-id: <your subscription id>
```

## When to Use

- CI/CD pipelines that deploy Azure infrastructure or push to ACR
- Replacing a stored `AZURE_CLIENT_SECRET` repository secret with keyless auth
- Any GitHub-to-Azure automation where secret rotation is unwelcome toil

## Key Configuration Choices

- **Branch-pinned subject** -- switch to
  `repo:{owner}/{repo}:environment:production` to gate access behind a
  protected environment's required reviewers instead of branch protection
- **One credential per trusted context** -- add sibling credentials (up to 20
  per identity) for tags or other branches rather than widening this one

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<user-assigned-identity-arm-id>` | The parent identity's ARM ID (or use `valueFrom` against an `AzureUserAssignedIdentity`) | The identity's `status.outputs.identity_id` |
| `<owner>/<repo>` | The GitHub repository allowed to authenticate | The repository URL |
