---
title: "Presets"
description: "Ready-to-deploy configuration presets for Elastic File System"
type: "preset-list"
componentSlug: "elastic-file-system"
componentTitle: "Elastic File System"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-general-purpose-regional"
    rank: "01"
    title: "General Purpose Regional EFS"
    excerpt: "Regional, encrypted, bursting throughput, backup enabled. Simplest production-safe starting point."
  - slug: "02-one-zone-dev"
    rank: "02"
    title: "One Zone Dev EFS"
    excerpt: "One Zone storage (us-east-1a), encrypted, bursting throughput. Lower cost for dev/test. Single subnet."
  - slug: "03-production-elastic-tiered"
    rank: "03"
    title: "Production Elastic EFS with Lifecycle Tiering and DR Replication"
    excerpt: "Regional, encrypted, elastic throughput, full lifecycle tiering (IA + Archive + warm-on-access), daily backups, and a cross-region replica for disaster recovery."
---

# Elastic File System Presets

Ready-to-deploy configuration presets for Elastic File System. Each preset is a complete manifest you can copy, customize, and deploy.
