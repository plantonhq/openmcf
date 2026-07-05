---
title: "Data Lake Storage Gen2 Account"
description: "This preset creates a hierarchical-namespace (ADLS Gen2) account: real directories, POSIX ACLs, and the dfs endpoint that analytics engines -- Spark, Databricks, Synapse -- address. SFTP ingestion is..."
type: "preset"
rank: "03"
presetSlug: "03-data-lake-gen2"
componentSlug: "storage-account"
componentTitle: "Storage Account"
provider: "azure"
icon: "package"
order: 3
---

# Data Lake Storage Gen2 Account

This preset creates a hierarchical-namespace (ADLS Gen2) account: real
directories, POSIX ACLs, and the dfs endpoint that analytics engines --
Spark, Databricks, Synapse -- address. SFTP ingestion is enabled for
upstream systems that deliver files over SFTP.

## When to Use

- The storage layer of an analytics platform or lakehouse
- Workloads that need directory semantics and POSIX ACLs on blob data
- SFTP-based data ingestion onto cloud storage

## Key Configuration Choices

- **`isHnsEnabled: true`** -- the create-time switch that makes this a
  data lake; it cannot be toggled later, and it excludes blob versioning
  (Azure's contract -- data protection on a lake comes from snapshots
  and soft delete instead)
- **`sftpEnabled: true`** -- the SFTP endpoint onto the blob service
  (bills per enabled hour; remove if unused)
- **ZRS replication** -- analytics data is expensive to regenerate
- **A raw-zone lifecycle rule** -- `raw/` data cools after 30 days and
  ages out after a year; add per-zone rules as the lake grows

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | The Azure region, e.g. `eastus` | Your region strategy |
| `<resource-group-resource-name>` | The AzureResourceGroup's Planton resource name | Your resource-group composition |
| `<accountname>` | 3-24 lowercase letters/digits, globally unique | Your naming convention (no hyphens!) |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

## Downstream Wiring

Analytics engines consume the dfs endpoint:

```yaml
# The output analytics configuration points at
status.outputs.primary_dfs_endpoint  # https://{account}.dfs.core.windows.net/
```
