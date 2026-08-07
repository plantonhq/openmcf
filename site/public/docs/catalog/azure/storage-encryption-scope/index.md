---
title: "Storage Encryption Scope"
description: "Storage Encryption Scope deployment documentation"
icon: "package"
order: 100
componentName: "azurestorageencryptionscope"
---

# Azure Storage Encryption Scope

Deploys an encryption scope inside an Azure Storage Account -- a named encryption boundary that lets different data in ONE account encrypt under DIFFERENT keys. The account's own encryption settings cover everything by default; a scope overrides that for the blobs, containers, and ADLS filesystems that opt into it. That is the multi-tenant pattern (one account, per-tenant keys -- revoking a tenant's key renders exactly that tenant's data unreadable) and the mixed-sensitivity pattern (platform-managed keys for ordinary data, a customer-managed key for the regulated subset) without splitting into per-tenant accounts. Scopes are many-per-account with independent lifecycles and are referenced BY NAME from containers and filesystems, which is why they are a first-class kind rather than a list folded into the account's spec.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Encryption Scope** -- a named scope on the referenced storage account (by ARM ID -- the control-plane path), with your chosen key ownership model (platform-managed or your Key Vault key) and optional infrastructure (double) encryption

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### Azure Subscription

- **An AzureStorageAccount** the scope will live in, referenced through `storageAccountId`. The parent is fixed at creation: a scope cannot move between accounts.
- **For a customer-managed key** (`source: MICROSOFT_KEY_VAULT`): an AzureKeyVaultKey, and the ACCOUNT must carry an identity with wrap/unwrap access on the key's vault -- the same plumbing as the account-level customer-managed key.

## Deploy

### Console

Open the deployment store, find **Azure Storage Encryption Scope**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Platform-Managed Scope** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureStorageEncryptionScope
metadata:
  name: tenant-a-scope
  org: acme-corp
  env: prod
spec:
  storageAccountId:
    valueFrom:
      kind: AzureStorageAccount
      name: app-storage
      fieldPath: status.outputs.storage_account_id
  scopeName: tenantAScope
  source: MICROSOFT_STORAGE
```

```shell
planton apply -f scope.yaml
```

This creates a platform-managed scope -- a distinct encryption boundary with zero key management: Azure creates and rotates the key.

### InfraChart

When deploying as part of a multi-resource environment, the ValueFromRef above wires the scope to its account: the InfraPipeline resolves the dependency graph, deploys the storage account first, then provisions the scope with the resolved ARM ID.

## Key Configuration

These are the most important decisions when configuring an encryption scope. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Scope name** -- `scopeName` is the scope's identity, not a label: containers (`defaultEncryptionScope`) and ADLS filesystems reference the scope by this exact name. 4-63 letters and digits; hyphens are NOT allowed (stricter than most storage names). Renaming replaces the scope and strands every reference to the old name.

**Key ownership** -- `source` is a required, explicit choice. `MICROSOFT_STORAGE` uses a platform-managed key Azure creates and rotates -- zero key management, right when the boundary exists for organization rather than control. `MICROSOFT_KEY_VAULT` encrypts under YOUR Key Vault key -- you control rotation, revocation, and access -- and requires `keyVaultKeyId`. The spec rejects a Key Vault choice without a key at validation time.

**The key reference** -- `keyVaultKeyId` should reference an AzureKeyVaultKey's VERSIONLESS ID so rotation propagates to the scope with zero intervention; pin a versioned URI only when a compliance regime demands a frozen version.

**Double encryption** -- `infrastructureEncryptionRequired` adds a second, independent platform-keyed layer underneath the scope's key, for regimes that demand defense against a single-algorithm or single-key compromise. Independent of the ACCOUNT's infrastructure-encryption switch, so a scope can add the second layer for just the regulated subset. Fixed at creation.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureStorageAccount** | `storageAccountId` | `status.outputs.storage_account_id` |
| **AzureKeyVaultKey** | `keyVaultKeyId` | `status.outputs.versionless_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `encryption_scope_name` | The scope's name within the account | AzureStorageContainer `defaultEncryptionScope`, ADLS filesystem `defaultEncryptionScope` -- scopes compose BY NAME |
| `encryption_scope_id` | Azure Resource Manager ID of the scope | Diagnostics and policy audits |
| `storage_account_name` | The parent account's name, parsed from the resolved account ID | The account/scope pair without a second reference |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Per-tenant key isolation** -- one platform-managed scope per tenant, each tenant's container pinned to its scope with the per-blob override blocked. Start from the **Platform-Managed Scope** preset.

**Regulated subset under a customer key** -- `source: MICROSOFT_KEY_VAULT` with a versionless key reference; revoking the vault key is the kill switch. Start from the **Customer-Managed-Key Scope** preset.

**Double encryption for one data class** -- `infrastructureEncryptionRequired: true` on the scope while the account's own switch stays off. Start from the **Double-Encryption Scope** preset.

## Works With

- [**Azure Storage Account**](/cloud-catalog/azure-storage-account) -- the parent account; with a customer-managed key it must carry an identity with wrap/unwrap access on the key's vault
- [**Azure Storage Container**](/cloud-catalog/azure-storage-container) -- pins the scope through `defaultEncryptionScope`, optionally blocking the per-blob override for hard isolation
- [**Azure Key Vault Key**](/cloud-catalog/azure-key-vault-key) -- the customer-managed key, referenced versionless so rotation is transparent
