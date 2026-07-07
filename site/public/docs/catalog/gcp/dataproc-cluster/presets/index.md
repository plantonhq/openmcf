---
title: "Presets"
description: "Ready-to-deploy configuration presets for Dataproc Cluster"
type: "preset-list"
componentSlug: "dataproc-cluster"
componentTitle: "Dataproc Cluster"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-dev-jupyter"
    rank: "01"
    title: "Dev Jupyter"
    excerpt: "A lightweight development cluster with JupyterLab for interactive data exploration and prototyping Spark jobs, wired to delete itself after 30 minutes of inactivity."
  - slug: "02-ha-production"
    rank: "02"
    title: "HA Production"
    excerpt: "A high-availability Dataproc cluster for production Spark workloads: 3 masters, SSD storage with NVMe local SSDs, Shielded VMs, CMEK encryption, private networking, and OSS metrics into Cloud..."
  - slug: "03-cost-optimized-batch"
    rank: "03"
    title: "Cost-Optimized Batch"
    excerpt: "An ephemeral Dataproc cluster tuned for batch Spark jobs: a small on-demand base, a Spot secondary group with machine-type flexibility and a standard/spot capacity mix, an attached autoscaling..."
  - slug: "04-spark-on-gke"
    rank: "04"
    title: "Spark on GKE"
    excerpt: "A Dataproc-on-GKE virtual cluster: Spark workloads run as Kubernetes pods on an existing GKE cluster instead of dedicated Compute Engine VMs, sharing the GKE cluster's capacity, autoscaling, and..."
---

# Dataproc Cluster Presets

Ready-to-deploy configuration presets for Dataproc Cluster. Each preset is a complete manifest you can copy, customize, and deploy.
