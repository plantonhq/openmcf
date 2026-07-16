---
title: "Tenant-Scoped Encryption Container"
description: "This preset creates a private container pinned to an encryption scope: every blob written to it encrypts under the scope's key, and individual writes cannot override the scope -- per-tenant key..."
type: "preset"
rank: "03"
presetSlug: "03-tenant-scoped-encryption"
componentSlug: "storage-container"
componentTitle: "Storage Container"
provider: "azure"
icon: "package"
order: 3
---

# Tenant-Scoped Encryption Container

This preset creates a private container pinned to an encryption scope:
every blob written to it encrypts under the scope's key, and individual
writes cannot override the scope -- per-tenant key isolation inside a
shared storage account.

## When to Use

- Multi-tenant platforms giving each tenant a container encrypted under
  the tenant's own key
- Compliance regimes that require sub-account key boundaries without
  the cost of an account per tenant

## Key Configuration Choices

- **`defaultEncryptionScope`** -- references an AzureStorageEncryptionScope
  on the same account; that scope's key covers this container's blobs
- **`encryptionScopeOverrideEnabled: false`** -- the scope is mandatory;
  a blob write naming a different scope is rejected
- Both are **fixed at creation** -- key-isolation boundaries should not
  drift after the fact

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<storage-account-resource-name>` | The AzureStorageAccount's Planton resource name | Your storage composition |
| `tenant-data` | 3-63 lowercase letters/digits/hyphens | Your naming convention |
| `<encryption-scope-resource-name>` | The AzureStorageEncryptionScope's Planton resource name (must live on the same account) | Your storage composition |
| `<tenant-id>` | The tenant this container isolates | Your tenancy model |
