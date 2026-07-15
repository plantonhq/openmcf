---
title: "Presets"
description: "Ready-to-deploy configuration presets for Planton Runner"
type: "preset-list"
componentSlug: "planton-runner"
componentTitle: "Planton Runner"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-private-vpc-worker"
    rank: "01"
    title: "Private VPC Worker"
    excerpt: "This preset deploys the standard runner appliance: a pull-based worker on two private subnets that receives deploy operations through its queue and executes them from inside the VPC. The 30-second..."
  - slug: "02-dual-mode"
    rank: "02"
    title: "Dual Mode (Deploys + Live CloudOps)"
    excerpt: "This preset runs the runner in `dual` mode: everything the private VPC worker does, plus the real-time CloudOps channel -- live browsing of the resources behind the runner (pods, services, cluster..."
  - slug: "03-high-capacity"
    rank: "03"
    title: "High Capacity (Production Hardened)"
    excerpt: "This preset is the production posture for heavy workloads: a sized-up runner (1 vCPU / 4 GiB) with a pinned version, a first-class runtime IAM role, and extended log retention. Everything is..."
---

# Planton Runner Presets

Ready-to-deploy configuration presets for Planton Runner. Each preset is a complete manifest you can copy, customize, and deploy.
