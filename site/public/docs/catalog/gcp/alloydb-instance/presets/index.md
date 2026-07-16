---
title: "Presets"
description: "Ready-to-deploy configuration presets for AlloyDB Instance"
type: "preset-list"
componentSlug: "alloydb-instance"
componentTitle: "AlloyDB Instance"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-read-pool-basic"
    rank: "01"
    title: "Basic Read Pool"
    excerpt: "This preset adds a single-node READ_POOL instance to an existing AlloyDB cluster for offloading read traffic from the bundled primary."
  - slug: "02-read-pool-ha"
    rank: "02"
    title: "HA Read Pool"
    excerpt: "This preset creates a regional READ_POOL with two nodes for higher read availability."
  - slug: "03-read-pool-production"
    rank: "03"
    title: "Production Read Pool"
    excerpt: "This preset creates a regional three-node READ_POOL with connector enforcement, TLS-only connections, and query insights enabled."
---

# AlloyDB Instance Presets

Ready-to-deploy configuration presets for AlloyDB Instance. Each preset is a complete manifest you can copy, customize, and deploy.
