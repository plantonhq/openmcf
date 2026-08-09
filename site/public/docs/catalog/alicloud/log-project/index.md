---
title: "Log Project"
description: "Log Project deployment documentation"
icon: "package"
order: 100
componentName: "alicloudlogproject"
---

# AliCloud Log Project

Deploys an Alibaba Cloud Simple Log Service (SLS) project with bundled log stores and full-text search indexes. The component provisions the project, creates each specified log store with configurable retention and sharding, and enables full-text indexing per store by default so logs are immediately searchable after ingestion. The project integrates with Planton's Provider Connections for AliCloud credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **SLS Project** -- an `alicloud_log_project` with the specified name, description, resource group, and tags
- **Log Stores** -- one `alicloud_log_store` per entry in the `logStores` list, each with configurable retention period, shard count, auto-split behavior, and metadata enrichment
- **Full-Text Store Indexes** -- one `alicloud_log_store_index` per log store where `enableIndex` is `true` (the default), configured with case-insensitive matching and standard tokenization
- **AliCloud Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) merged with user-provided `tags`

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### Alibaba Cloud Account

- **A globally unique project name** -- SLS project names must be unique across all Alibaba Cloud accounts within a region. Use a descriptive name that includes your organization or environment prefix.

## Deploy

### Console

Open the deployment store, find **AliCloud Log Project**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Production Multi-Store** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudLogProject
metadata:
  name: app-logging
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  projectName: acme-prod-logs
  logStores:
    - name: app-logs
    - name: audit-logs
      retentionDays: 365
```

```shell
planton apply -f log-project.yaml
```

This creates an SLS project named `acme-prod-logs` in `cn-hangzhou` with two log stores: `app-logs` with 30-day retention (default) and `audit-logs` with 365-day retention. Both stores have full-text indexing enabled, 2 shards with auto-split, and metadata enrichment turned on.

## Key Configuration

These are the most important decisions when configuring a Log Project. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Log store retention** -- Set `retentionDays` per log store to control how long data is retained. Defaults to 30 days. Set to 3650 for permanent retention. Use shorter retention for development stores and longer retention for audit or compliance stores.

**Shard count and auto-split** -- Each shard supports approximately 5 MB/s write throughput. Set `shardCount` based on expected ingestion volume. Enable `autoSplit` (default: `true`) to let SLS automatically split shards when write throughput exceeds capacity, up to `maxSplitShardCount` (default: 64).

**Full-text indexing** -- When `enableIndex` is `true` (default), a full-text search index is created on the log store, making logs immediately queryable after ingestion. Disable indexing (`enableIndex: false`) for archive-only stores where query capability is not needed, eliminating index storage costs.

**Metadata enrichment** -- When `appendMeta` is `true` (default), SLS automatically appends log receive time and client IP as metadata fields on each log entry. Useful for debugging ingestion issues and correlating logs with source hosts.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `project_name` | SLS project name (also serves as the project identifier in SLS APIs) | Downstream components referencing this project for log ingestion |
| `project_id` | SLS project resource ID | Monitoring dashboards, audit references |
| `log_store_names` | Map of log store names created within the project | Downstream components referencing specific stores via ValueFromRef |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Production multi-store** -- Separate log stores for application logs, audit trails, and access logs with distinct retention periods and shard configurations. Tags enable cost attribution and organizational filtering. Start from the **Production Multi-Store** preset.

**Development** -- A single log store with short retention and minimal shards for development and testing environments. Start from the **Development** preset.

## Works With

This component operates independently and does not reference other components.