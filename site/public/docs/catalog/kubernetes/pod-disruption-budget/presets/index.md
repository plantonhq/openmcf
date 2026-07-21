---
title: "Presets"
description: "Ready-to-deploy configuration presets for Pod Disruption Budget"
type: "preset-list"
componentSlug: "pod-disruption-budget"
componentTitle: "Pod Disruption Budget"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-protect-workload"
    rank: "01"
    title: "Protect Workload"
    excerpt: "This preset is the standard availability floor for one workload: at least one of its pods must survive any voluntary disruption. Node drains, cluster upgrades, and autoscaler consolidation all go..."
  - slug: "02-tier-percentage"
    rank: "02"
    title: "Tier Percentage"
    excerpt: "This preset spans several workloads with one budget: every pod labelled `tier: web` or `tier: api` is covered, and at most a quarter of them may be down at once during voluntary disruptions. It is..."
  - slug: "03-crashloop-tolerant"
    rank: "03"
    title: "Crashloop Tolerant"
    excerpt: "This preset is the availability floor that cannot wedge a node drain. Under the default eviction policy, a budget counts only READY pods as available — so a crash-looping application (running, never..."
---

# Pod Disruption Budget Presets

Ready-to-deploy configuration presets for Pod Disruption Budget. Each preset is a complete manifest you can copy, customize, and deploy.
