---
title: "General-Purpose v2 Account"
description: "This preset creates a StorageV2 account on the standard tier with local redundancy and blob data protection -- the baseline for application assets, uploads, and scratch data. Containers are added as..."
type: "preset"
rank: "01"
presetSlug: "01-general-purpose-v2"
componentSlug: "storage-account"
componentTitle: "Storage Account"
provider: "azure"
icon: "package"
order: 1
---

# General-Purpose v2 Account

This preset creates a StorageV2 account on the standard tier with local
redundancy and blob data protection -- the baseline for application
assets, uploads, and scratch data. Containers are added as separate
AzureStorageContainer resources referencing this account.

## When to Use

- The default storage account for an application environment
- Dev/test environments where local redundancy is acceptable
- The backend a Function App or Web App binds to

## Key Configuration Choices

- **StorageV2 / Standard / LRS** (the spec defaults) -- every data
  service available, cheapest redundancy
- **Blob versioning + soft delete** -- overwrites and deletes stay
  recoverable for 7 days (Azure's default window)
- **No firewall block** -- reachable from all networks (Azure's
  default); add `networkRules` when locking down

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | The Azure region, e.g. `eastus` | Your region strategy |
| `<resource-group-resource-name>` | The AzureResourceGroup's Planton resource name | Your resource-group composition |
| `<accountname>` | 3-24 lowercase letters/digits, globally unique | Your naming convention (no hyphens!) |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

## Downstream Wiring

Add containers and bind app services through the account's outputs:

```yaml
# On an AzureStorageContainer
storageAccountId:
  valueFrom:
    kind: AzureStorageAccount
    name: app-storage
    fieldPath: status.outputs.storage_account_id
```
