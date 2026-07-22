---
title: "Presets"
description: "Ready-to-deploy configuration presets for KEDA"
type: "preset-list"
componentSlug: "keda"
componentTitle: "KEDA"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-cluster-standard"
    rank: "01"
    title: "Cluster Standard"
    excerpt: "This preset installs the KEDA engine cluster-wide with the upstream defaults: the operator, the external-metrics API server, and the admission webhooks in a dedicated `keda` namespace, watching..."
  - slug: "02-eks-irsa-scalers"
    rank: "02"
    title: "EKS with IRSA for AWS Scalers"
    excerpt: "This preset installs KEDA on an EKS cluster with IAM Roles for Service Accounts (IRSA) wired for the scalers: KEDA's own service account is annotated with an IAM role, so scalers that read AWS metric..."
  - slug: "03-ha-production"
    rank: "03"
    title: "HA Production"
    excerpt: "This preset hardens KEDA for clusters where autoscaling is load-bearing: two replicas of every component (warm standbys — KEDA leader-elects, so extra replicas buy failover speed, not throughput),..."
---

# KEDA Presets

Ready-to-deploy configuration presets for KEDA. Each preset is a complete manifest you can copy, customize, and deploy.
