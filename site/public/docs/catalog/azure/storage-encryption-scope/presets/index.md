---
title: "Presets"
description: "Ready-to-deploy configuration presets for Storage Encryption Scope"
type: "preset-list"
componentSlug: "storage-encryption-scope"
componentTitle: "Storage Encryption Scope"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-platform-managed-scope"
    rank: "01"
    title: "Platform-Managed Encryption Scope"
    excerpt: "This preset creates an encryption scope with a Microsoft-managed key -- a distinct encryption boundary inside a shared account with ZERO key management: Azure creates and rotates the key."
  - slug: "02-customer-managed-key-scope"
    rank: "02"
    title: "Customer-Managed-Key Encryption Scope"
    excerpt: "This preset creates an encryption scope backed by YOUR Key Vault key -- customer key custody for just the data that opts into the scope, while the rest of the account stays on platform-managed keys."
  - slug: "03-double-encryption-scope"
    rank: "03"
    title: "Double-Encryption Scope"
    excerpt: "This preset creates a scope with infrastructure (double) encryption: two independent encryption layers with independent keys and algorithms for the data that opts in -- WITHOUT enabling account-wide..."
---

# Storage Encryption Scope Presets

Ready-to-deploy configuration presets for Storage Encryption Scope. Each preset is a complete manifest you can copy, customize, and deploy.
