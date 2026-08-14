---
title: "Presets"
description: "Ready-to-deploy configuration presets for Machine Learning Compute Cluster"
type: "preset-list"
componentSlug: "machine-learning-compute-cluster"
componentTitle: "Machine Learning Compute Cluster"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-cpu-training-cluster"
    rank: "01"
    title: "CPU Training Cluster"
    excerpt: "This preset creates a scale-to-zero dedicated CPU cluster -- the everyday shared training pool: free between jobs, four general-purpose nodes at peak, a system identity for credential-free data..."
  - slug: "02-gpu-training-cluster"
    rank: "02"
    title: "GPU Training Cluster"
    excerpt: "This preset creates a scale-to-zero GPU cluster for deep-learning training -- two V100 nodes at peak, released ten minutes after their last job because GPU idle time is the most expensive waste in an..."
  - slug: "03-low-priority-batch-cluster"
    rank: "03"
    title: "Low-Priority Batch Cluster"
    excerpt: "This preset creates a wide, cheap, evictable cluster for fault-tolerant batch work -- eight spot-class nodes at a fraction of dedicated cost, releasing five minutes after their last job."
---

# Machine Learning Compute Cluster Presets

Ready-to-deploy configuration presets for Machine Learning Compute Cluster. Each preset is a complete manifest you can copy, customize, and deploy.
