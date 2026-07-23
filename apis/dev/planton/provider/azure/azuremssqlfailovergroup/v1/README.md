# AzureMssqlFailoverGroup

## Overview

`AzureMssqlFailoverGroup` provisions an Azure SQL Failover Group: a
disaster-recovery grouping that replicates a set of databases from a primary
logical server to one or more partner servers in other regions, behind a
single listener endpoint that follows the primary through a failover.
Applications connect to `{name}.database.windows.net` instead of a specific
server, so a failover redirects traffic without a connection-string change.

## Why a First-Class Resource?

A failover group is a composable DR node with its own lifecycle:

- **Spans servers** -- it references a primary server, one or more partner
  servers, and a set of databases, each an independently-managed resource
- **The listener is the contract** -- the read-write endpoint always points
  at the current primary; the group is what makes that indirection real
- **Add/remove without rebuild** -- databases and the failover policy change
  in place

## Key Features

- **Cross-region replication** -- partner servers in different regions
  receive continuous geo-replication of every database in the group
- **Automatic or manual failover** -- `AUTOMATIC` with a grace period (≥ 60
  minutes, the window for an outage to self-heal without data loss), or
  `MANUAL` operator-initiated
- **Read-only listener** -- an optional second endpoint routes read-only
  workloads to a secondary
- **Composable** -- primary/partner servers and databases are all references
  defaulting to their Planton kinds

## Spec Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Group name (globally unique; the listener DNS label); fixed at creation |
| `server_id` | StringValueOrRef | Yes | Primary server (defaults to AzureMssqlServer); fixed at creation |
| `partner_servers` | repeated message | Yes (≥1) | Partner servers, each with a `server_id` |
| `database_ids` | repeated StringValueOrRef | No | Databases on the primary to replicate (default AzureMssqlDatabase) |
| `read_write_endpoint_failover_policy` | message | Yes | `mode` (AUTOMATIC/MANUAL) + `grace_minutes` |
| `readonly_endpoint_failover_policy_enabled` | bool | No | Fail over the read-only listener too (disabled when unset) |
| `tags` | map | No | User tags, merged over Planton-derived tags |

## Outputs

| Output | Description |
|--------|-------------|
| `failover_group_id` | Full ARM ID of the group |
| `failover_group_name` | The group's name (listener DNS label) |
| `read_write_listener_endpoint` | `{name}.database.windows.net` -- the failover-following endpoint |
| `read_only_listener_endpoint` | `{name}.secondary.database.windows.net` |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureMssqlFailoverGroup
metadata:
  name: prod-sql-fog
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

## Lifecycle Notes

- `name` and `server_id` are **fixed at creation** -- the group is anchored
  to its primary
- Databases must live on the **primary** server
- `AUTOMATIC` requires `grace_minutes` ≥ 60; `MANUAL` must omit it
- Partner servers must each be in a **different region** than the primary

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
