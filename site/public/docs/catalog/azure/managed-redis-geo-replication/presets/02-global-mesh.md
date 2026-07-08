---
title: "Global Active Mesh"
description: "This preset links four Managed Redis instances across continents into one active geo-replication group -- a write-anywhere global cache where every region serves local reads and writes and Azure..."
type: "preset"
rank: "02"
presetSlug: "02-global-mesh"
componentSlug: "managed-redis-geo-replication"
componentTitle: "Managed Redis Geo-Replication"
provider: "azure"
icon: "package"
order: 2
---

# Global Active Mesh

This preset links four Managed Redis instances across continents into
one active geo-replication group -- a write-anywhere global cache where
every region serves local reads and writes and Azure merges the
datasets conflict-free.

## When to Use

- Global applications with write traffic on every continent
  (leaderboards, presence, session stores, feature flags)
- Follow-the-sun workloads where the active region shifts through the
  day
- Whole-continent disaster tolerance

## Key Configuration Choices

- **Up to five members per group** -- the managing member plus four
  linked; this preset uses four regions total
- **One resource for the whole mesh** -- linking is reciprocal; the
  managing member is just the one the group is administered through
- **Region evacuation is a list edit** -- remove a member's ID and it
  force-unlinks, keeping its data; add it back to re-join

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<americas-member-resource-name>` | The managing member's AzureManagedRedis | Your per-region member manifests |
| `<europe-member-resource-name>` | The Europe member | Your per-region member manifests |
| `<asia-member-resource-name>` | The Asia member | Your per-region member manifests |
| `<oceania-member-resource-name>` | The Oceania member | Your per-region member manifests |
