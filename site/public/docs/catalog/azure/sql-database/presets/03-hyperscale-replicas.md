---
title: "Hyperscale Database with Readable Replicas"
description: "This preset creates a Hyperscale database: elastic storage to 100 TB, 4 vCores of independent compute, two readable high-availability replicas spread across availability zones, and zone-redundant..."
type: "preset"
rank: "03"
presetSlug: "03-hyperscale-replicas"
componentSlug: "sql-database"
componentTitle: "SQL Database"
provider: "azure"
icon: "package"
order: 3
---

# Hyperscale Database with Readable Replicas

This preset creates a Hyperscale database: elastic storage to 100 TB,
4 vCores of independent compute, two readable high-availability
replicas spread across availability zones, and zone-redundant backup
storage.

## When to Use

- Databases expected to grow beyond the 4 TB vCore ceiling
- Read-heavy workloads that want read-intent connections served by
  replicas instead of the primary
- Workloads that need near-instant failover (a replica promotes in
  seconds)

## Key Configuration Choices

- **`HS_Gen5_4`** -- Hyperscale's log-service architecture makes storage
  elastic and backups near-instant snapshots; note leaving Hyperscale
  later REPLACES the database (ARM's contract)
- **`readReplicaCount: 2`** -- connections with
  `ApplicationIntent=ReadOnly` load-balance across the replicas;
  the count is also the failover pool
- **`ZONE_REDUNDANT` backup storage** -- restores survive a zone outage
  without paying for paired-region geo-replication; switch to
  `GEO_ZONE_REDUNDANT` when cross-region restore matters too

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<server-resource-name>` | The AzureMssqlServer's Planton resource name | Your server composition |
| `<database-name>` | The database's name on the server | Your application |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |
