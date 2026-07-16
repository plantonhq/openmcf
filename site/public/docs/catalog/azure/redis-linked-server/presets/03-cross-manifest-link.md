---
title: "Cross-Manifest Geo Link"
description: "This preset links to a secondary cache that is NOT managed in the same manifest set -- another team's cache, or one provisioned outside Planton -- by passing its ARM id and region as literal values."
type: "preset"
rank: "03"
presetSlug: "03-cross-manifest-link"
componentSlug: "redis-linked-server"
componentTitle: "Redis Linked Server"
provider: "azure"
icon: "package"
order: 3
---

# Cross-Manifest Geo Link

This preset links to a secondary cache that is NOT managed in the same
manifest set -- another team's cache, or one provisioned outside
Planton -- by passing its ARM id and region as literal values.

## When to Use

- The DR cache belongs to a different team or composition
- Migrating an existing, externally-created geo pair under Planton
  management one side at a time

## Key Configuration Choices

- **Literals replace references only where composition is impossible**
  -- when both caches live in one manifest set, prefer the fully
  referenced form (preset 01) so ids and regions can never drift
- **The literal region must match the linked cache's actual region** --
  with literals there is no reference keeping them honest; ARM rejects a
  mismatch at link time

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<primary-cache-resource-name>` | The primary AzureRedisCache's Planton resource name | Your cache composition |
| `<linkedRedisCacheArmId>` | The external secondary cache's ARM id | `az redis show --name <cache> --resource-group <rg> --query id` |
| `<linkedRedisCacheRegion>` | The external secondary cache's region | The same `az redis show` output (`location`) |
