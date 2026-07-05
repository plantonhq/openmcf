---
title: "General Purpose Production Database"
description: "This preset creates a provisioned General Purpose database (2 vCores, 128 GB) with a real backup posture: a 14-day point-in-time restore window plus long-term weekly and monthly backups. Transparent..."
type: "preset"
rank: "01"
presetSlug: "01-general-purpose"
componentSlug: "sql-database"
componentTitle: "SQL Database"
provider: "azure"
icon: "package"
order: 1
---

# General Purpose Production Database

This preset creates a provisioned General Purpose database (2 vCores,
128 GB) with a real backup posture: a 14-day point-in-time restore
window plus long-term weekly and monthly backups. Transparent data
encryption and geo-redundant backup storage ride Azure's defaults.

## When to Use

- Steady production workloads where predictable performance and cost
  beat serverless elasticity
- Any database that needs restores beyond the 35-day PITR ceiling

## Key Configuration Choices

- **`GP_Gen5_2`** -- the balanced production tier; scale vCores in place
  by changing the SKU (Hyperscale transitions replace the database)
- **Long-term retention** -- weekly backups kept 12 weeks, monthly kept
  12 months; restore them into a NEW AzureMssqlDatabase with
  `createMode: RESTORE_LONG_TERM_RETENTION_BACKUP`
- **`maxSizeGb: 128`** -- the storage ceiling, independent of compute in
  vCore tiers

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<server-resource-name>` | The AzureMssqlServer's Planton resource name | Your server composition |
| `<database-name>` | The database's name on the server | Your application |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

## Downstream Wiring

```text
Server={server fqdn},1433;Database=<database-name>;User ID={admin};Password={password};Encrypt=True;
```
