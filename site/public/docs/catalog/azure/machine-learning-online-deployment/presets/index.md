---
title: "Presets"
description: "Ready-to-deploy configuration presets for Machine Learning Online Deployment"
type: "preset-list"
componentSlug: "machine-learning-online-deployment"
componentTitle: "Machine Learning Online Deployment"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-registered-model-deployment"
    rank: "01"
    title: "Registered-Model Deployment"
    excerpt: "This preset serves a registered model version from the workspace on one managed instance -- the everyday deployment an MLflow model needs, since the service infers scoring code and environment for..."
  - slug: "02-hardened-monitored-deployment"
    rank: "02"
    title: "Hardened Monitored Deployment"
    excerpt: "This preset is the production posture: secure egress through the workspace's managed network, Application Insights on, a patient startup probe for slow-loading models, honest request limits, and..."
---

# Machine Learning Online Deployment Presets

Ready-to-deploy configuration presets for Machine Learning Online Deployment. Each preset is a complete manifest you can copy, customize, and deploy.
