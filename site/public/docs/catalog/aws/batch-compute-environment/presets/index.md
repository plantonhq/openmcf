---
title: "Presets"
description: "Ready-to-deploy configuration presets for Batch Compute Environment"
type: "preset-list"
componentSlug: "batch-compute-environment"
componentTitle: "Batch Compute Environment"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-fargate-batch"
    rank: "01"
    title: "Fargate Batch"
    excerpt: "A serverless Fargate compute environment — zero instances to manage, per-second billing, scale-to-zero when queues are empty. The right default for containerized batch jobs that fit Fargate's sizing..."
  - slug: "02-ec2-managed-batch"
    rank: "02"
    title: "EC2 Managed Batch"
    excerpt: "An EC2 On-Demand compute environment with `optimal` instance selection and the BEST_FIT_PROGRESSIVE strategy — the configuration that keeps day-2 infrastructure changes cheap (in-place updates..."
  - slug: "03-spot-cost-optimized-batch"
    rank: "03"
    title: "Spot Cost-Optimized Batch"
    excerpt: "An EC2 Spot compute environment using the SPOT_PRICE_CAPACITY_OPTIMIZED strategy — up to ~90% below On-Demand pricing for retry-tolerant batch workloads, with instance selection that balances price..."
---

# Batch Compute Environment Presets

Ready-to-deploy configuration presets for Batch Compute Environment. Each preset is a complete manifest you can copy, customize, and deploy.
