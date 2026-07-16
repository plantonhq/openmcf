---
title: "Manual Failover Group"
description: "This preset creates a failover group with manual failover -- an operator (or an automation pipeline) initiates every failover. This is the choice when you want a human decision in the loop before..."
type: "preset"
rank: "02"
presetSlug: "02-manual-failover"
componentSlug: "mssql-failover-group"
componentTitle: "MSSQL Failover Group"
provider: "azure"
icon: "package"
order: 2
---

# Manual Failover Group

This preset creates a failover group with manual failover -- an operator (or
an automation pipeline) initiates every failover. This is the choice when
you want a human decision in the loop before redirecting production traffic
to another region, for example when regional failover has downstream
coordination costs (DNS, caches, dependent services) you want to sequence
deliberately.

The read-only listener failover is disabled here so read-only workloads stay
pinned to their secondary until you decide otherwise.

## When to Use

- DR where failover must be a deliberate, coordinated action
- Compliance regimes that require an explicit failover authorization
- Testing failover procedures on a schedule you control

## Key Configuration Choices

- **`mode: MANUAL`** -- no `graceMinutes` (it is only meaningful for
  automatic failover); Azure never fails over on its own
- **`readonlyEndpointFailoverPolicyEnabled: false`** -- the read-only
  listener does not move with the primary

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<primary-server>` | The primary AzureMssqlServer | Your primary server's Planton resource name |
| `<partner-server>` | The partner AzureMssqlServer (different region) | Your DR server's Planton resource name |
| `<database>` | An AzureMssqlDatabase on the primary to protect | Your database's Planton resource name |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |
