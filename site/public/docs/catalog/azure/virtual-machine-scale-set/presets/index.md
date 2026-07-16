---
title: "Presets"
description: "Ready-to-deploy configuration presets for Virtual Machine Scale Set"
type: "preset-list"
componentSlug: "virtual-machine-scale-set"
componentTitle: "Virtual Machine Scale Set"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-stateless-web-flexible"
    rank: "01"
    title: "Stateless Web Fleet (Flexible)"
    excerpt: "This preset creates a FLEXIBLE-orchestration Linux fleet built for stateless web workloads: three zone-spread instances on ephemeral OS disks, accelerated networking, load-balancer pool membership..."
  - slug: "02-spot-batch"
    rank: "02"
    title: "Spot Batch Fleet (Flexible, Mixed SKUs)"
    excerpt: "This preset creates a FLEXIBLE-orchestration spot fleet for interruption-tolerant batch work: ten instances drawn from three interchangeable VM sizes (capacity-optimized), a two-instance guaranteed..."
  - slug: "03-windows-uniform-rolling"
    rank: "03"
    title: "Windows Fleet with Automatic OS Upgrades (Uniform)"
    excerpt: "This preset creates a UNIFORM-orchestration Windows Server fleet that patches itself at the image level: automatic OS image upgrades roll new `latest` releases across zone-balanced instances in..."
---

# Virtual Machine Scale Set Presets

Ready-to-deploy configuration presets for Virtual Machine Scale Set. Each preset is a complete manifest you can copy, customize, and deploy.
