---
title: "Presets"
description: "Ready-to-deploy configuration presets for FSx for ONTAP"
type: "preset-list"
componentSlug: "fsx-for-ontap"
componentTitle: "FSx for ONTAP"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-single-az-development"
    rank: "01"
    title: "Single-AZ Development FSx ONTAP"
    excerpt: "SINGLE_AZ_2 SSD file system with 1 TiB (1024 GiB) and 384 MB/s throughput. One HA pair. No automatic backups. The smallest and cheapest ONTAP configuration for development and testing."
  - slug: "02-single-az-production"
    rank: "02"
    title: "Single-AZ Production FSx ONTAP"
    excerpt: "SINGLE_AZ_2 SSD file system with 2 TiB (2048 GiB) and 768 MB/s throughput. One HA pair. Customer-managed KMS encryption. 7-day automatic backups at 05:00 UTC. Production-ready configuration for..."
  - slug: "03-multi-az-high-availability"
    rank: "03"
    title: "Multi-AZ High Availability FSx ONTAP"
    excerpt: "MULTI_AZ_2 deployment with automatic failover across two availability zones. 2 TiB SSD, 768 MB/s throughput. 7-day backups, customer-managed KMS encryption. Mission-critical configuration for..."
  - slug: "04-single-az-scale-out"
    rank: "04"
    title: "Scale-Out FSx ONTAP (Multiple HA Pairs)"
    excerpt: "SINGLE_AZ_2 scale-out file system with 4 HA pairs delivering 6 GB/s aggregate throughput (1536 MB/s per pair) and 16 TiB of SSD storage with 160,000 provisioned IOPS. The shape for high-throughput..."
---

# FSx for ONTAP Presets

Ready-to-deploy configuration presets for FSx for ONTAP. Each preset is a complete manifest you can copy, customize, and deploy.
