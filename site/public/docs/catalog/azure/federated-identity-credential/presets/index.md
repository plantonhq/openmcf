---
title: "Presets"
description: "Ready-to-deploy configuration presets for Federated Identity Credential"
type: "preset-list"
componentSlug: "federated-identity-credential"
componentTitle: "Federated Identity Credential"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-github-actions-oidc"
    rank: "01"
    title: "GitHub Actions Keyless CI"
    excerpt: "This preset trusts a GitHub repository's `main` branch to authenticate as a managed identity -- the standard shape for secretless CI/CD. Workflows request an OIDC token (`permissions: id-token:..."
  - slug: "02-aks-workload-identity"
    rank: "02"
    title: "AKS Workload Identity"
    excerpt: "This preset trusts a Kubernetes service account in an AKS cluster to authenticate as a managed identity -- the modern way for pods to reach Key Vault, Storage, or any RBAC-granted Azure resource..."
  - slug: "03-external-oidc-issuer"
    rank: "03"
    title: "Generic External OIDC Issuer"
    excerpt: "This preset trusts any OIDC-compliant external system -- GitLab, a self-hosted CI, another cloud's workload identity, an internal token service -- to authenticate as a managed identity. Azure AD only..."
---

# Federated Identity Credential Presets

Ready-to-deploy configuration presets for Federated Identity Credential. Each preset is a complete manifest you can copy, customize, and deploy.
