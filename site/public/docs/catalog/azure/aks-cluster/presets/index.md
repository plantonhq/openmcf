---
title: "Presets"
description: "Ready-to-deploy configuration presets for AKS Cluster"
type: "preset-list"
componentSlug: "aks-cluster"
componentTitle: "AKS Cluster"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard"
    rank: "01"
    title: "Standard Production AKS Cluster"
    excerpt: "This preset deploys a production-ready AKS cluster with a public API endpoint, Azure CNI Overlay networking, a 3-zone autoscaling system pool tainted for system pods only, and the recommended add-ons..."
  - slug: "02-private"
    rank: "02"
    title: "Private AKS Cluster with Workload Identity"
    excerpt: "This preset deploys a private AKS cluster with no public API server endpoint and the workload-identity loop enabled. The Kubernetes API is accessible only from within the VNet or via peered networks..."
  - slug: "03-hardened-enterprise"
    rank: "03"
    title: "Hardened Enterprise AKS Cluster"
    excerpt: "This preset deploys a compliance-posture AKS cluster: Azure AD RBAC with local accounts disabled, API-server access restricted to authorized ranges, Cilium-compatible Azure network policy, host..."
---

# AKS Cluster Presets

Ready-to-deploy configuration presets for AKS Cluster. Each preset is a complete manifest you can copy, customize, and deploy.
