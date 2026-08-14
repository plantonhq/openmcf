---
title: "Presets"
description: "Ready-to-deploy configuration presets for Machine Learning Batch Endpoint"
type: "preset-list"
componentSlug: "machine-learning-batch-endpoint"
componentTitle: "Machine Learning Batch Endpoint"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-batch-scoring-endpoint"
    rank: "01"
    title: "Batch Scoring Endpoint"
    excerpt: "This preset creates the everyday batch endpoint: the minimal routing object batch deployments attach to, with Microsoft Entra authentication applied by default (the only mode the batch service..."
  - slug: "02-routed-endpoint-with-identity"
    rank: "02"
    title: "Routed Endpoint with Identity"
    excerpt: "This preset creates a batch endpoint with a system-assigned identity and an explicit default-deployment pointer -- the shape for teams whose role-grant conventions target the endpoint object and..."
---

# Machine Learning Batch Endpoint Presets

Ready-to-deploy configuration presets for Machine Learning Batch Endpoint. Each preset is a complete manifest you can copy, customize, and deploy.
