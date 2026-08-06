---
title: "MSSQL Failover Group"
description: "MSSQL Failover Group deployment documentation"
icon: "package"
order: 100
componentName: "azuremssqlfailovergroup"
---

# Azure MSSQL Failover Group

Creates an Azure SQL Failover Group — a disaster-recovery grouping that replicates databases from a primary logical server to partner servers in other regions, behind a single listener endpoint that follows the primary through a failover. Applications connect to the listener, not a server, so a failover needs no connection-string change.

## What Gets Created

When you deploy an AzureMssqlFailoverGroup resource, Planton provisions:

- **Failover Group** — an `azurerm_mssql_failover_group` on the primary server, replicating the listed databases to the partner servers with the configured failover policy

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A primary logical server** (an `AzureMssqlServer`) with the databases to protect
- **One or more partner servers** in different regions (each an `AzureMssqlServer`)
- **SQL write rights**: `Microsoft.Sql/servers/failoverGroups/write`

## Quick Start

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureMssqlFailoverGroup
metadata:
  name: prod-sql-fog
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureMssqlFailoverGroup.prod-sql-fog
spec:
  name: prod-sql-fog
  serverId:
    valueFrom:
      name: prod-sql-primary
  partnerServers:
    - serverId:
        valueFrom:
          name: prod-sql-dr
  databaseIds:
    - valueFrom:
        name: orders-db
  readWriteEndpointFailoverPolicy:
    mode: AUTOMATIC
    graceMinutes: 60
```

Deploy:

```shell
planton apply -f failover-group.yaml
```

After deployment, point applications at `status.outputs.read_write_listener_endpoint`.

## Configuration Reference

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `name` | `string` | Group name; globally unique (the listener DNS label). Fixed at creation. |
| `serverId` | `StringValueOrRef` | The primary logical server. Defaults to an `AzureMssqlServer` reference. Fixed at creation. |
| `partnerServers` | `object[]` | Partner servers (each `{ serverId }`), one or more, each in a different region. |
| `readWriteEndpointFailoverPolicy` | `object` | `mode` (`AUTOMATIC` / `MANUAL`) and `graceMinutes` (≥ 60 for AUTOMATIC, omitted for MANUAL). |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `databaseIds` | `StringValueOrRef[]` | `[]` | Databases on the primary to replicate. Default `AzureMssqlDatabase` references. |
| `readonlyEndpointFailoverPolicyEnabled` | `bool` | `false` | Fail over the read-only listener too (disabled when unset). |
| `tags` | `map(string)` | `{}` | User tags, merged over Planton-derived tags (user wins). |

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `failover_group_id` | `string` | Full ARM ID of the group |
| `failover_group_name` | `string` | The group's name (listener DNS label) |
| `read_write_listener_endpoint` | `string` | `{name}.database.windows.net` |
| `read_only_listener_endpoint` | `string` | `{name}.secondary.database.windows.net` |

## Related Components

- [AzureMssqlServer](/docs/catalog/azure/mssql-server) — the primary and partner logical servers
- [AzureMssqlDatabase](/docs/catalog/azure/mssql-database) — the databases the group replicates
- [AzureResourceGroup](/docs/catalog/azure/resource-group) — provides the resource groups for the servers
