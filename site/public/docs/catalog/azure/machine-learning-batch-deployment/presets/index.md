---
title: "Presets"
description: "Ready-to-deploy configuration presets for Machine Learning Batch Deployment"
type: "preset-list"
componentSlug: "machine-learning-batch-deployment"
componentTitle: "Machine Learning Batch Deployment"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-model-scoring-recipe"
    rank: "01"
    title: "Model Scoring Recipe"
    excerpt: "This preset creates the classic batch inference recipe: a registered model scored on a compute-cluster pool, with the batching dials tuned explicitly and an appended predictions file as output."
  - slug: "02-pipeline-component-recipe"
    rank: "02"
    title: "Pipeline Component Recipe"
    excerpt: "This preset creates a pipeline-component batch deployment: each endpoint invocation runs a registered pipeline component as a pipeline job -- multi-step batch work (prepare, score, post-process)..."
---

# Machine Learning Batch Deployment Presets

Ready-to-deploy configuration presets for Machine Learning Batch Deployment. Each preset is a complete manifest you can copy, customize, and deploy.
