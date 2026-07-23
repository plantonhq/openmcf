---
title: "Presets"
description: "Ready-to-deploy configuration presets for Karpenter Node Pool"
type: "preset-list"
componentSlug: "karpenter-node-pool"
componentTitle: "Karpenter Node Pool"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-general-purpose-on-demand"
    rank: "01"
    title: "General Purpose On-Demand"
    excerpt: "This preset declares the default fleet most clusters start with: amd64 on-demand nodes from current-generation compute-, general- and memory-optimized instance families, with consolidation active and..."
  - slug: "02-spot-diversified"
    rank: "02"
    title: "Spot Diversified"
    excerpt: "This preset declares a cost-optimized spot pool with deliberate instance-family diversity: a wide `In` list of families plus `minValues` forcing the launched fleet to span at least four of them, so a..."
  - slug: "03-gpu-dedicated"
    rank: "03"
    title: "GPU Dedicated"
    excerpt: "This preset declares a tainted, on-demand GPU pool constrained to the g5 instance family. The `nvidia.com/gpu` taint is the standard dedicated-pool pattern: only pods that tolerate it schedule onto —..."
---

# Karpenter Node Pool Presets

Ready-to-deploy configuration presets for Karpenter Node Pool. Each preset is a complete manifest you can copy, customize, and deploy.
