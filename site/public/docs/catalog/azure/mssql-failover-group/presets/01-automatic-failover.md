---
title: "Automatic Failover Group"
description: "This preset creates a failover group with automatic failover and a 60-minute grace period -- the standard production DR shape. Azure fails the group over to the partner on its own when it detects a..."
type: "preset"
rank: "01"
presetSlug: "01-automatic-failover"
componentSlug: "mssql-failover-group"
componentTitle: "MSSQL Failover Group"
provider: "azure"
icon: "package"
order: 1
---

# Automatic Failover Group

This preset creates a failover group with automatic failover and a 60-minute
grace period -- the standard production DR shape. Azure fails the group over
to the partner on its own when it detects a sustained outage, after waiting
the grace period (the window in which the outage might self-heal without
data loss). Applications connect to the read-write listener and never change
their connection string across a failover.

## When to Use

- Production databases that need cross-region disaster recovery
- Any workload where an unplanned regional outage should recover
  automatically within the RTO the grace period allows

## Key Configuration Choices

- **`mode: AUTOMATIC` + `graceMinutes: 60`** -- the minimum grace period;
  raise it if you prefer to give a transient outage more time to recover
  before a failover with potential data loss
- **Single partner** -- one partner server in another region covers the
  common DR topology; add more for multi-region

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<primary-server>` | The primary AzureMssqlServer | Your primary server's Planton resource name |
| `<partner-server>` | The partner AzureMssqlServer (different region) | Your DR server's Planton resource name |
| `<database>` | An AzureMssqlDatabase on the primary to protect | Your database's Planton resource name |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

## Downstream Wiring

Point your application's connection string at the read-write listener:

```
Server=tcp:prod-sql-fog.database.windows.net;Database=...;
```
