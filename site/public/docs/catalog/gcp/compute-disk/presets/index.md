---
title: "Presets"
description: "Ready-to-deploy configuration presets for Compute Disk"
type: "preset-list"
componentSlug: "compute-disk"
componentTitle: "Compute Disk"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-data-volume"
    rank: "01"
    title: "Data Volume"
    excerpt: "The default posture for stateful VMs: an empty pd-balanced data disk that exists as its own resource, so the data outlives whatever instance mounts it."
  - slug: "02-encrypted-database-volume"
    rank: "02"
    title: "Encrypted Database Volume"
    excerpt: "The posture for regulated data: a high-IOPS SSD volume under customer-managed encryption, with a destroy-time snapshot so even a mistaken teardown cannot lose the data."
  - slug: "03-hyperdisk-high-iops"
    rank: "03"
    title: "Hyperdisk High IOPS"
    excerpt: "Storage for workloads whose performance needs will change: a hyperdisk-balanced volume where IOPS and throughput are dials, not consequences of size."
---

# Compute Disk Presets

Ready-to-deploy configuration presets for Compute Disk. Each preset is a complete manifest you can copy, customize, and deploy.
