---
title: "Storage Data Lake Gen2 Filesystem"
description: "Storage Data Lake Gen2 Filesystem deployment documentation"
icon: "package"
order: 100
componentName: "azurestoragedatalakegen2filesystem"
---

# Azure Storage Data Lake Gen2 Filesystem

Creates a Data Lake Storage Gen2 filesystem inside an AzureStorageAccount -- the namespace unit of an analytics data lake. Analytics engines address data as `abfss://{filesystem}@{account}.dfs.core.windows.net/path`, and the filesystem is the grant boundary for data-plane RBAC and root-path POSIX ACLs: the one-filesystem-per-zone (raw, curated, gold) pattern.

## What Gets Created

When you deploy an AzureStorageDataLakeGen2Filesystem resource, Planton provisions:

- **Data Lake Gen2 Filesystem** -- an `azurerm_storage_data_lake_gen2_filesystem` on the referenced account, with optional root-path ownership (owner/group), a POSIX ACL, a default encryption scope, and filesystem properties

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureStorageAccount with `isHnsEnabled: true`** (referenced through `storageAccountId`) -- POSIX access control is rejected on flat-namespace accounts
- **Data-plane reachability**: filesystems are created through the account's dfs endpoint, so a data-plane firewall that blocks the deploy runner blocks the create even though ARM would allow it

## Quick Start

Create a file `filesystem.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureStorageDataLakeGen2Filesystem
metadata:
  name: raw-zone
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureStorageDataLakeGen2Filesystem.raw-zone
spec:
  storageAccountId:
    valueFrom:
      kind: AzureStorageAccount
      name: my-lake-storage
      fieldPath: status.outputs.storage_account_id
  filesystemName: raw-zone
  aces:
    - type: USER
      permissions: rwx
    - type: GROUP
      permissions: r-x
    - type: OTHER
      permissions: "---"
```

Deploy:

```shell
planton apply -f filesystem.yaml
```

## Key Outputs

| Output | Purpose |
|--------|---------|
| `filesystem_id` | The ARM container-proxy id -- what data-plane role assignments (Storage Blob Data Reader/Contributor/Owner) scope to |
| `filesystem_name` | The container segment of every abfss:// and dfs URL |
| `storage_account_name` | The account/filesystem pair, without a second reference |

## Related Resources

- [Azure Storage Account](/docs/catalog/azure/storage-account) -- the parent (HNS-enabled) account
- [Azure Storage Encryption Scope](/docs/catalog/azure/storage-encryption-scope) -- sub-account key isolation the filesystem can pin
- [Azure Role Assignment](/docs/catalog/azure/role-assignment) -- grants data-plane roles at the `filesystem_id` scope
