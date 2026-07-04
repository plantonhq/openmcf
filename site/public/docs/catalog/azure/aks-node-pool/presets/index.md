---
title: "Presets"
description: "Ready-to-deploy configuration presets for AKS Node Pool"
type: "preset-list"
componentSlug: "aks-node-pool"
componentTitle: "AKS Node Pool"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-on-demand-general"
    rank: "01"
    title: "On-Demand General Purpose Node Pool"
    excerpt: "This preset creates a general-purpose AKS user node pool with on-demand (regular) VMs, autoscaling from 2 to 10 nodes across 3 availability zones. This is the standard configuration for production..."
  - slug: "02-spot-cost-optimized"
    rank: "02"
    title: "Spot Cost-Optimized Node Pool"
    excerpt: "This preset creates a cost-optimized AKS user node pool using Azure Spot VMs, which provide 30-90% savings over on-demand pricing. The pool scales to zero when idle and up to 20 nodes under load...."
  - slug: "03-gpu-or-windows"
    rank: "03"
    title: "GPU Node Pool Reserved by Taint"
    excerpt: "This preset creates a single-node GPU pool (NVIDIA T4) with the AKS-managed GPU driver installed, labeled and tainted so only GPU workloads schedule onto it. Scale the fixed count -- or switch to..."
---

# AKS Node Pool Presets

Ready-to-deploy configuration presets for AKS Node Pool. Each preset is a complete manifest you can copy, customize, and deploy.
