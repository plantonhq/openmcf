---
title: "Two-Region Active Pair"
description: "This preset links two Managed Redis instances in different regions into an active geo-replication group: both accept writes, Azure merges the datasets conflict-free, and applications read and write..."
type: "preset"
rank: "01"
presetSlug: "01-two-region-pair"
componentSlug: "managed-redis-geo-replication"
componentTitle: "Managed Redis Geo Replication"
provider: "azure"
icon: "package"
order: 1
---

# Two-Region Active Pair

This preset links two Managed Redis instances in different regions into
an active geo-replication group: both accept writes, Azure merges the
datasets conflict-free, and applications read and write their local
region.

## When to Use

- Active-active disaster recovery -- losing a region loses no writes
  accepted in the other
- Two-region applications that want local read/write latency everywhere
- The starting shape for a global group (add members later, up to five)

## Key Configuration Choices

- **One resource for the whole group** -- linking is reciprocal, so
  never create one per member
- **Both members declare the same `geoReplicationGroupName`** and are
  `BALANCED_B3` or larger -- Azure enforces both at link time
- **Deleting this resource unlinks the members** -- each keeps its own
  copy of the data and becomes an independent instance; removing one ID
  from the linked list evacuates just that region

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<primary-region-member-resource-name>` | The AzureManagedRedis managed through | Your per-region member manifests |
| `<secondary-region-member-resource-name>` | The other region's AzureManagedRedis | Your per-region member manifests |
