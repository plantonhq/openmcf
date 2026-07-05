---
title: "Re-Link After Failover"
description: "This preset closes the DR loop: after a regional failover promoted the secondary (by deleting the original link), it re-establishes replication in the opposite direction once the failed region..."
type: "preset"
rank: "02"
presetSlug: "02-relink-after-failover"
componentSlug: "redis-linked-server"
componentTitle: "Redis Linked Server"
provider: "azure"
icon: "package"
order: 2
---

# Re-Link After Failover

This preset closes the DR loop: after a regional failover promoted the
secondary (by deleting the original link), it re-establishes replication
in the opposite direction once the failed region recovers.

## When to Use

- The recovery half of every geo-replication runbook: fail over
  (delete the link), then re-protect (create this one)
- Planned region evacuations, where the "failover" is a controlled
  migration

## Key Configuration Choices

- **The roles have swapped** -- the promoted cache is now the target
  (primary) and the recovered region's cache joins as SECONDARY; the
  recovered cache's data is overwritten by replication
- **The geo hostname follows the pair** -- applications pointed at
  `geo_replicated_primary_host_name` keep working through both the
  failover and the re-link
- **Fail back later by repeating the cycle** -- delete this link,
  repoint, and link the other way again

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<promoted-cache-resource-name>` | The cache promoted during the failover | Your DR runbook |
| `<recovered-cache-resource-name>` | The recovered region's cache rejoining as replica | Your DR runbook |
