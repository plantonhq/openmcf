---
title: "Presets"
description: "Ready-to-deploy configuration presets for SageMaker Model Registry"
type: "preset-list"
componentSlug: "sagemaker-model-registry"
componentTitle: "SageMaker Model Registry"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-team-model-registry"
    rank: "01"
    title: "Team Model Registry"
    excerpt: "This preset is a single team's model package group: the named shell a training pipeline registers versioned model packages into for review, approval, and deployment."
  - slug: "02-shared-model-registry"
    rank: "02"
    title: "Shared Model Registry"
    excerpt: "This preset is a cross-account registry group: a resource policy on the group lets another AWS account discover and list the model versions registered into it."
---

# SageMaker Model Registry Presets

Ready-to-deploy configuration presets for SageMaker Model Registry. Each preset is a complete manifest you can copy, customize, and deploy.
