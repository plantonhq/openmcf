---
title: "Presets"
description: "Ready-to-deploy configuration presets for Batch Job Queue"
type: "preset-list"
componentSlug: "batch-job-queue"
componentTitle: "Batch Job Queue"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-single-environment"
    rank: "01"
    title: "Single Environment Queue"
    excerpt: "The standard starting point: one queue mapped onto one compute environment, with stuck-job protection so a mis-sized job cannot block the queue forever."
  - slug: "02-spot-overflow"
    rank: "02"
    title: "Spot-First With On-Demand Overflow"
    excerpt: "The canonical Batch cost pattern: jobs try the Spot environment first (order 1) and spill to On-Demand (order 2) only when Spot capacity runs out. Most work rides ~90%-discounted compute; nothing..."
---

# Batch Job Queue Presets

Ready-to-deploy configuration presets for Batch Job Queue. Each preset is a complete manifest you can copy, customize, and deploy.
