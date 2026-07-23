---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cluster Autoscaler"
type: "preset-list"
componentSlug: "cluster-autoscaler"
componentTitle: "Cluster Autoscaler"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-eks-autodiscovery"
    rank: "01"
    title: "EKS Autodiscovery"
    excerpt: "This preset installs the Cluster Autoscaler on EKS in the recommended posture: tag-based ASG auto-discovery, keyless AWS access via IRSA, the `least-waste` expander, and balanced multi-AZ node..."
  - slug: "02-cluster-api"
    rank: "02"
    title: "Cluster API"
    excerpt: "This preset installs the Cluster Autoscaler with the Cluster API provider in the self-managed mode: both the workload cluster and the CAPI management objects live in the same cluster, machine..."
  - slug: "03-azure-vmss"
    rank: "03"
    title: "Azure VMSS"
    excerpt: "This preset installs the Cluster Autoscaler against Azure VM scale sets with tag-based auto-discovery and federated workload identity — the keyless credential posture. It is the arm for AKS clusters..."
---

# Cluster Autoscaler Presets

Ready-to-deploy configuration presets for Cluster Autoscaler. Each preset is a complete manifest you can copy, customize, and deploy.
