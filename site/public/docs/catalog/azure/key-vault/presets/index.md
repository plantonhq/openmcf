---
title: "Presets"
description: "Ready-to-deploy configuration presets for Key Vault"
type: "preset-list"
componentSlug: "key-vault"
componentTitle: "Key Vault"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard-rbac"
    rank: "01"
    title: "Standard RBAC Vault"
    excerpt: "This preset creates the baseline modern vault: Standard SKU, Azure RBAC authorization (the spec default), public endpoint, and Azure's own deletion-safety defaults (90-day soft delete, purge..."
  - slug: "02-premium-network-restricted"
    rank: "02"
    title: "Premium Network-Restricted CMK Vault"
    excerpt: "This preset creates the production posture for a vault holding customer-managed keys: Premium SKU (unlocks the HSM-backed key types for the `AzureKeyVaultKey` resources inside), purge protection ON..."
  - slug: "03-legacy-access-policy"
    rank: "03"
    title: "Legacy Access-Policy Vault"
    excerpt: "This preset runs the vault in the legacy access-policy authorization mode: `rbacAuthorizationEnabled: false` with explicit per-principal permission lists carried on the vault itself. ARM stores but..."
---

# Key Vault Presets

Ready-to-deploy configuration presets for Key Vault. Each preset is a complete manifest you can copy, customize, and deploy.
