---
title: "Presets"
description: "Ready-to-deploy configuration presets for Spanner Instance"
type: "preset-list"
componentSlug: "spanner-instance"
componentTitle: "Spanner Instance"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-free-instance"
    rank: "01"
    title: "Free Instance"
    excerpt: "Provisions the project's one zero-cost Cloud Spanner instance — full Spanner semantics (strong consistency, SQL, the same client libraries) with limited capacity (about 10 GB of storage). Nothing..."
  - slug: "02-regional-production"
    rank: "02"
    title: "Regional Production"
    excerpt: "Provisions a production-ready Cloud Spanner instance with fixed capacity (one node), ENTERPRISE edition, and automatic backup schedules for new databases. Suitable for production workloads with..."
  - slug: "03-autoscaling-production"
    rank: "03"
    title: "Autoscaling Production (Multi-Region)"
    excerpt: "Provisions a multi-region Cloud Spanner instance with Spanner's managed autoscaler: capacity follows utilization within explicit bounds, including an asymmetric override that scales one read-heavy..."
---

# Spanner Instance Presets

Ready-to-deploy configuration presets for Spanner Instance. Each preset is a complete manifest you can copy, customize, and deploy.
