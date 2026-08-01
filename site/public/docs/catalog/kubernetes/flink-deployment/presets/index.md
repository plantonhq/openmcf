---
title: "Presets"
description: "Ready-to-deploy configuration presets for Flink Deployment"
type: "preset-list"
componentSlug: "flink-deployment"
componentTitle: "Flink Deployment"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-application-stateless"
    rank: "01"
    title: "Application (stateless) preset"
    excerpt: "The recommended production grain: one APPLICATION cluster per pipeline — the cluster exists to run this one job and follows its lifecycle, with a custom image that BAKES the job jar (`local:///`..."
  - slug: "02-stateful-ha-s3"
    rank: "02"
    title: "Stateful HA-on-S3 preset"
    excerpt: "The full durable-state posture, every requirement the operator enforces made explicit:"
---

# Flink Deployment Presets

Ready-to-deploy configuration presets for Flink Deployment. Each preset is a complete manifest you can copy, customize, and deploy.
