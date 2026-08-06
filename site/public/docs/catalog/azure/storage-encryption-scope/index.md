---
title: "Storage Encryption Scope"
description: "Storage Encryption Scope deployment documentation"
icon: "package"
order: 100
componentName: "azurestorageencryptionscope"
---

# Azure Storage Encryption Scope

Creates an encryption scope inside an AzureStorageAccount -- a named encryption boundary that lets different data in one account encrypt under different keys. Containers and blobs opt in by name: the multi-tenant and mixed-sensitivity key-isolation patterns without per-tenant accounts.

## What Gets Created

When you deploy an AzureStorageEncryptionScope resource, Planton provisions:

- **Encryption Scope** -- an `azurerm_storage_encryption_scope` on the referenced account, sourcing its key from Microsoft.Storage (platform-managed) or your Key Vault (customer-managed)

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureStorageAccount** to create the scope in (referenced through `storageAccountId`)
- **For the Key Vault source**: an AzureKeyVaultKey, and the account must carry a managed identity with wrap/unwrap access on the key's vault (the same plumbing as the account-level customer-managed key)

## Quick Start

Create a file `scope.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureStorageEncryptionScope
metadata:
  name: tenant42-scope
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureStorageEncryptionScope.tenant42-scope
spec:
  storageAccountId:
    valueFrom:
      kind: AzureStorageAccount
      name: my-app-storage
      fieldPath: status.outputs.storage_account_id
  scopeName: tenant42scope
  source: MICROSOFT_STORAGE
```

Deploy:

```shell
planton apply -f scope.yaml
```

Note Azure's unusual delete semantics: destroying a scope soft-disables it (ARM has no true delete), and the scope's name stays reserved within the account -- recreating the same name re-enables it.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `encryption_scope_id` | The scope's ARM id |
| `encryption_scope_name` | What containers (`default_encryption_scope`), ADLS filesystems, and per-blob upload options reference |
| `storage_account_name` | The account/scope pair, without a second reference |

## Related Resources

- [Azure Storage Account](/docs/catalog/azure/storage-account) -- the parent account
- [Azure Storage Container](/docs/catalog/azure/storage-container) -- pins a scope as its default
- [Azure Key Vault Key](/docs/catalog/azure/key-vault-key) -- the customer-managed key source
