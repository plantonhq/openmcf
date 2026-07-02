---
title: "Presets"
description: "Ready-to-deploy configuration presets for EKS Node Group"
type: "preset-list"
componentSlug: "eks-node-group"
componentTitle: "EKS Node Group"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-on-demand-general"
    rank: "01"
    title: "On-Demand General Pool"
    excerpt: "This preset runs the workhorse node pool of a typical cluster: On-Demand AL2023 nodes across two availability zones, surge-enabled version rollouts, and managed node auto-repair. Everything composes..."
  - slug: "02-spot-cost-optimized"
    rank: "02"
    title: "Spot Cost-Optimized Pool"
    excerpt: "This preset runs an interruptible batch/burst pool on Spot capacity at a steep discount, with the two Spot survival practices built in: several similar instance types for pool diversity, and a taint..."
  - slug: "03-launch-template"
    rank: "03"
    title: "Launch-Template Pool"
    excerpt: "This preset composes the node group onto a first-class `AwsLaunchTemplate`: the template owns the launch mechanics (instance type, IMDSv2 enforcement, encrypted volumes, instance tags), and promoting..."
---

# EKS Node Group Presets

Ready-to-deploy configuration presets for EKS Node Group. Each preset is a complete manifest you can copy, customize, and deploy.
