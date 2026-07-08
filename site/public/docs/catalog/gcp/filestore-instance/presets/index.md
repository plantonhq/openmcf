---
title: "Presets"
description: "Ready-to-deploy configuration presets for Filestore Instance"
type: "preset-list"
componentSlug: "filestore-instance"
componentTitle: "Filestore Instance"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-dev-basic"
    rank: "01"
    title: "Dev Basic"
    excerpt: "The minimal working Filestore instance for development, testing, and CI: SSD-backed, on the default VPC, easy to tear down."
  - slug: "02-production-enterprise"
    rank: "02"
    title: "Production Enterprise"
    excerpt: "The production posture: a regional tier that survives zone failures, deletion protection as the destroy guard, private-services networking, and locked-down NFS exports."
  - slug: "03-high-performance-zonal"
    rank: "03"
    title: "High Performance Zonal"
    excerpt: "The throughput posture: the modern ZONAL tier with IOPS provisioned per terabyte, so performance scales automatically as the share grows, plus customer-managed encryption."
---

# Filestore Instance Presets

Ready-to-deploy configuration presets for Filestore Instance. Each preset is a complete manifest you can copy, customize, and deploy.
