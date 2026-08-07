---
title: "Presets"
description: "Ready-to-deploy configuration presets for FSx OpenZFS File System"
type: "preset-list"
componentSlug: "fsx-openzfs-file-system"
componentTitle: "FSx OpenZFS File System"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-single-az-development"
    rank: "01"
    title: "Preset: Single-AZ Development"
    excerpt: "**Use case**: Development and testing environments where cost efficiency matters more than high availability or performance."
  - slug: "02-single-az-production"
    rank: "02"
    title: "Preset: Single-AZ Production"
    excerpt: "**Use case**: Production workloads in a single availability zone where NFS performance, data compression, encryption, and daily backups are required."
  - slug: "03-multi-az-high-availability"
    rank: "03"
    title: "Preset: Multi-AZ High Availability"
    excerpt: "**Use case**: Mission-critical production workloads requiring automatic failover across availability zones, provisioned IOPS, storage quotas, and extended backup retention."
  - slug: "04-multi-az-intelligent-tiering"
    rank: "04"
    title: "Multi-AZ Intelligent-Tiering FSx OpenZFS"
    excerpt: "MULTI_AZ_1 file system on the INTELLIGENT_TIERING storage class — elastic, pay-for-what-you-store NFS capacity with 1280 MB/s of provisioned throughput, an AWS-sized SSD read cache, and automatic..."
---

# FSx OpenZFS File System Presets

Ready-to-deploy configuration presets for FSx OpenZFS File System. Each preset is a complete manifest you can copy, customize, and deploy.
