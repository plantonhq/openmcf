---
title: "Presets"
description: "Ready-to-deploy configuration presets for Ray Cluster"
type: "preset-list"
componentSlug: "ray-cluster"
componentTitle: "Ray Cluster"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-lab"
    rank: "01"
    title: "Lab preset"
    excerpt: "The smallest useful Ray cluster: a head that also RUNS tasks (`scheduleTasksOnHead: true` — the deliberate lab inversion of the production default, where the head advertises zero CPUs to keep..."
  - slug: "02-production-autoscaling"
    rank: "02"
    title: "Production autoscaling preset"
    excerpt: "The production shape, all four load-bearing choices made:"
---

# Ray Cluster Presets

Ready-to-deploy configuration presets for Ray Cluster. Each preset is a complete manifest you can copy, customize, and deploy.
