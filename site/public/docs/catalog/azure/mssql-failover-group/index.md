---
title: "MSSQL Failover Group"
description: "MSSQL Failover Group deployment documentation"
icon: "package"
order: 100
componentName: "azuremssqlfailovergroup"
---

# Azure MSSQL Failover Group

Deploys a cross-region failover group pairing an Azure SQL logical server (the primary) with one or more partner servers. The group replicates the databases it lists and exposes two LISTENER endpoints — a read-write listener that always points at the current primary and a read-only listener pointed at the secondary — so applications survive a regional failover with no connection-string change. That indirection is the group's whole value.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **SQL Failover Group** -- created ON the primary server, spanning the partner server(s), with the group name becoming the listener DNS prefix (`{name}.database.windows.net`)
- **Database Replication** -- a maintained replica of every listed database on every partner server
- **Failover Policy** -- automatic (Azure promotes the partner after the grace window) or manual (an operator decides), plus the read-only listener's failover behavior
- **Azure Tags** -- resource metadata tags applied to the group

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Two AzureMssqlServers** -- the primary and the partner, usually in different regions (Azure's paired-region guidance applies). Both referenced through their `server_id` outputs.
- **AzureMssqlDatabases** (optional at creation) -- the databases that replicate, referenced through their `database_id` outputs. An empty group is legal; membership edits in place.

### Azure Subscription

- **Keep both servers' postures in sync** -- the partner is a full logical server with its own authentication and firewall; a failed-over application hits ITS door.
- **Data-loss planning** -- automatic failover loses whatever replication lag exists when the grace window (≥ 60 minutes) expires; manual failover keeps that acceptance with your operators.

## Deploy

### Console

Open the deployment store, find **Azure MSSQL Failover Group**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields — with a live topology diagram resolving the primary and partner as you pick them. Start from the **automatic-failover** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureMssqlFailoverGroup
metadata:
  name: appdb-dr
  org: acme-corp
  env: prod
spec:
  name: appdb-dr
  serverId:
    valueFrom:
      kind: AzureMssqlServer
      name: app-sql
      fieldPath: status.outputs.server_id
  partnerServers:
    - serverId:
        valueFrom:
          kind: AzureMssqlServer
          name: app-sql-dr
          fieldPath: status.outputs.server_id
  databaseIds:
    - valueFrom:
        kind: AzureMssqlDatabase
        name: app-database
        fieldPath: status.outputs.database_id
  readWriteEndpointFailoverPolicy:
    mode: AUTOMATIC
    graceMinutes: 60
```

```shell
planton apply -f mssql-failover-group.yaml
```

This creates an automatic-failover group replicating one database to the partner, with listeners at `appdb-dr.database.windows.net` (read-write) and `appdb-dr.secondary.database.windows.net` (read-only). A Stack Job tracks the provisioning in real time.

### InfraChart

The group references both servers and every listed database — the InfraPipeline deploys the servers first, then the databases, then the group, in one resolved graph.

## Key Configuration

These are the most important decisions when configuring a failover group. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**The name is the listener** -- `name` becomes the DNS prefix applications connect to. Point every connection string at the listener, never a server directly. Fixed at creation.

**Failover mode** -- `readWriteEndpointFailoverPolicy.mode`: AUTOMATIC requires `graceMinutes` ≥ 60 (the data-loss dial — Azure waits it out before promoting, and whatever replication lag exists is lost); MANUAL forbids the grace window and keeps the decision with your operators.

**Read-only listener failover** -- `readonlyEndpointFailoverPolicyEnabled` decides whether read-intent traffic follows a failover (sharing the surviving server) or goes dark until failback (Azure's default).

**Membership** -- `databaseIds` lists exactly what replicates; nothing joins implicitly. Replicas bill at the partner's rates, and pooled databases need a matching pool on the partner.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureMssqlServer** | `serverId` (primary) | `status.outputs.server_id` |
| **AzureMssqlServer** | `partnerServers[].serverId` | `status.outputs.server_id` |
| **AzureMssqlDatabase** | `databaseIds` | `status.outputs.database_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `read_write_listener_endpoint` | The listener that follows the primary | Application connection strings |
| `read_only_listener_endpoint` | The listener pointed at the secondary | Reporting/read-intent connection strings |
| `failover_group_id` | Azure resource ID of the group | Diagnostics, automation |
| `failover_group_name` | Name of the group | Monitoring, dashboards |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Automatic failover** -- The hands-off DR posture: Azure promotes the partner after the grace window. Start from the **automatic-failover** preset.

**Manual failover** -- Tier-1 estates with a runbook: a human confirms the region is really gone before accepting the lag loss. Start from the **manual-failover** preset.

## Works With

- [**Azure MSSQL Server**](/cloud-catalog/azure-mssql-server) -- the primary and every partner
- [**Azure MSSQL Database**](/cloud-catalog/azure-mssql-database) -- the databases that replicate through the group
