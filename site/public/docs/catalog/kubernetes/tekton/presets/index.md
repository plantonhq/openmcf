---
title: "Presets"
description: "Ready-to-deploy configuration presets for Tekton"
type: "preset-list"
componentSlug: "tekton"
componentTitle: "Tekton"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-ci-standard"
    rank: "01"
    title: "CI standard preset"
    excerpt: "The full Tekton control plane — Pipelines, Triggers, Dashboard and Chains — with the one piece of configuration no production cluster should skip: a pruner. Completed runs keep their pods around..."
  - slug: "02-pipelines-engine"
    rank: "02"
    title: "Pipelines engine preset"
    excerpt: "Tekton as an embedded engine: the `lite` profile runs Pipelines alone — no Triggers, no Dashboard, no Chains — for platforms that create PipelineRuns programmatically and present their own UI. Every..."
---

# Tekton Presets

Ready-to-deploy configuration presets for Tekton. Each preset is a complete manifest you can copy, customize, and deploy.
