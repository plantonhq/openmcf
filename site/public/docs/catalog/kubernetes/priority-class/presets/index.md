---
title: "Presets"
description: "Ready-to-deploy configuration presets for Priority Class"
type: "preset-list"
componentSlug: "priority-class"
componentTitle: "Priority Class"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-critical-services"
    rank: "01"
    title: "Critical Services"
    excerpt: "This preset creates the top rung of the user-definable importance ladder: value 1,000,000 with preemption enabled. Pods that reference this class (via the workload pod spec's `priority_class_name`)..."
  - slug: "02-standard-default"
    rank: "02"
    title: "Standard Default"
    excerpt: "This preset creates the middle rung of the ladder AND makes it the cluster-wide default: every pod that names no priority class receives value 1000 instead of the bare-Kubernetes default of 0. This..."
  - slug: "03-preemptable-batch"
    rank: "03"
    title: "Preemptable Batch"
    excerpt: "This preset creates the bottom rung of the ladder: a negative-value, non-preempting class for work that should run only on spare capacity. Pods of this class yield to everything — including unmarked..."
---

# Priority Class Presets

Ready-to-deploy configuration presets for Priority Class. Each preset is a complete manifest you can copy, customize, and deploy.
