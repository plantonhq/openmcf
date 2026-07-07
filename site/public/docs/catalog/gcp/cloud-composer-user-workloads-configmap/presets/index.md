---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cloud Composer User Workloads ConfigMap"
type: "preset-list"
componentSlug: "cloud-composer-user-workloads-configmap"
componentTitle: "Cloud Composer User Workloads ConfigMap"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-dag-configuration"
    rank: "01"
    title: "DAG Configuration"
    excerpt: "Runtime configuration for a pipeline — endpoints, batch sizing, retry tuning — delivered as a Kubernetes ConfigMap into a Composer environment. The DAG reads values by key instead of hard-coding them."
  - slug: "02-feature-flags"
    rank: "02"
    title: "Feature Flags"
    excerpt: "Boolean flags for gating pipeline behavior — delivered as a Kubernetes ConfigMap so a flag flip is a config apply, not a DAG redeploy."
---

# Cloud Composer User Workloads ConfigMap Presets

Ready-to-deploy configuration presets for Cloud Composer User Workloads ConfigMap. Each preset is a complete manifest you can copy, customize, and deploy.
