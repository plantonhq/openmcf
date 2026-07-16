---
title: "HA Read Pool"
description: "This preset creates a regional READ_POOL with two nodes for higher read availability."
type: "preset"
rank: "02"
presetSlug: "02-read-pool-ha"
componentSlug: "alloydb-instance"
componentTitle: "AlloyDB Instance"
provider: "gcp"
icon: "package"
order: 2
---

# HA Read Pool

This preset creates a regional READ_POOL with two nodes for higher read availability.

## When to Use

- Production read scaling where queries must survive a zone loss
- Analytics or reporting workloads that fan out across multiple read nodes

## Key Configuration Choices

- **nodeCount: 2 + REGIONAL** — nodes span zones; minimum for regional read pools
- **cpuCount: 4** — moderate compute per node

## Related Presets

- **01-read-pool-basic** — single-node zonal pool for lower cost
- **03-read-pool-production** — hardened production settings
