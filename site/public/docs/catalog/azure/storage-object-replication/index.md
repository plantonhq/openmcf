---
title: "Storage Object Replication"
description: "Storage Object Replication deployment documentation"
icon: "package"
order: 100
componentName: "azurestorageobjectreplication"
---

# Azure Storage Object Replication

Creates an object replication policy between two AzureStorageAccounts -- asynchronous, rule-driven copying of block blobs from source containers to destination containers. The storage-level building block for cross-region DR, read-local data distribution, and archival fan-out.

## What Gets Created

When you deploy an AzureStorageObjectReplication resource, Planton provisions:

- **Object Replication Policy** -- an `azurerm_storage_object_replication` spanning the account pair; Azure materializes it on BOTH accounts (destination first, which assigns rule IDs, then the source mirror) under one shared policy GUID

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **Two AzureStorageAccounts**: the source with `blobProperties.versioningEnabled: true` and `changeFeedEnabled: true`, the destination with `versioningEnabled: true` -- Azure rejects the policy otherwise (which also rules out hierarchical-namespace accounts)
- **The mapped containers** on their respective accounts

## Quick Start

Create a file `replication.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureStorageObjectReplication
metadata:
  name: invoices-dr
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureStorageObjectReplication.invoices-dr
spec:
  sourceStorageAccountId:
    valueFrom:
      kind: AzureStorageAccount
      name: primary-storage
      fieldPath: status.outputs.storage_account_id
  destinationStorageAccountId:
    valueFrom:
      kind: AzureStorageAccount
      name: dr-storage
      fieldPath: status.outputs.storage_account_id
  rules:
    - sourceContainerName:
        valueFrom:
          kind: AzureStorageContainer
          name: invoices
          fieldPath: status.outputs.container_name
      destinationContainerName:
        valueFrom:
          kind: AzureStorageContainer
          name: invoices-replica
          fieldPath: status.outputs.container_name
      copyBlobsCreatedAfter: Everything
```

Deploy:

```shell
planton apply -f replication.yaml
```

Replication is asynchronous with no default RPO guarantee -- monitor it via the policy (`az storage account or-policy show --policy-id <policy_id>`).

## Key Outputs

| Output | Purpose |
|--------|---------|
| `source_object_replication_id` | The policy's ARM id on the source account |
| `destination_object_replication_id` | The policy's ARM id on the destination account (the authoritative copy) |
| `policy_id` | The shared policy GUID monitoring and the CLI key on |

## Related Resources

- [Azure Storage Account](/docs/catalog/azure/storage-account) -- both ends of the pair (versioning/change-feed prerequisites)
- [Azure Storage Container](/docs/catalog/azure/storage-container) -- what the rules map between
