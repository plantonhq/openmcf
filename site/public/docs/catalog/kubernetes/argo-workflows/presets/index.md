---
title: "Presets"
description: "Ready-to-deploy configuration presets for Argo Workflows"
type: "preset-list"
componentSlug: "argo-workflows"
componentTitle: "Argo Workflows"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-dev-pipelines"
    rank: "01"
    title: "Dev pipelines preset"
    excerpt: "The smallest useful Argo Workflows: the engine, the UI, and a runner identity — submit a Workflow CR (or use the UI over a port-forward; the command lands in the stack outputs) and it runs. The..."
  - slug: "02-durable-pipelines"
    rank: "02"
    title: "Durable pipelines preset"
    excerpt: "Argo Workflows with its two durability seams filled: an S3-compatible artifact store so steps pass files and archived logs outlive pods, and a Postgres archive so run history outlives the Workflow..."
  - slug: "03-multi-team"
    rank: "03"
    title: "Multi-team preset"
    excerpt: "One cluster, several workflow engines, zero crosstalk. The instance ID is the mechanism: this controller reconciles only Workflows carrying its `controller-instanceid` label and ignores the rest, so..."
---

# Argo Workflows Presets

Ready-to-deploy configuration presets for Argo Workflows. Each preset is a complete manifest you can copy, customize, and deploy.
