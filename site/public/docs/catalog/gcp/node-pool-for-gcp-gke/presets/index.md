---
title: "Presets"
description: "Ready-to-deploy configuration presets for Node Pool for GCP GKE"
type: "preset-list"
componentSlug: "node-pool-for-gcp-gke"
componentTitle: "Node Pool for GCP GKE"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-on-demand-autoscaling"
    rank: "01"
    title: "On-Demand Autoscaling Pool"
    excerpt: "This preset creates the workhorse node pool most clusters run first: on-demand VMs, per-zone autoscaling, surge upgrades with zero capacity dip, a dedicated node service account, and secure boot on."
  - slug: "02-spot-cost-optimized"
    rank: "02"
    title: "Spot Cost-Optimized Pool"
    excerpt: "This preset creates a scale-to-zero Spot pool for fault-tolerant workloads: deeply discounted capacity (60-91%) that costs nothing while idle, fenced off from general workloads by a taint."
  - slug: "03-gpu-accelerated"
    rank: "03"
    title: "GPU-Accelerated Pool"
    excerpt: "This preset creates a scale-to-zero GPU pool for ML workloads: NVIDIA L4 nodes with GKE-managed drivers, image streaming for large ML images, and an explicit GPU taint."
---

# Node Pool for GCP GKE Presets

Ready-to-deploy configuration presets for Node Pool for GCP GKE. Each preset is a complete manifest you can copy, customize, and deploy.
