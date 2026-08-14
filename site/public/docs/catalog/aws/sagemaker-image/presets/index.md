---
title: "Presets"
description: "Ready-to-deploy configuration presets for SageMaker Image"
type: "preset-list"
componentSlug: "sagemaker-image"
componentTitle: "SageMaker Image"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-custom-kernel-image"
    rank: "01"
    title: "Custom Kernel Image"
    excerpt: "This preset creates the registry entry for a team's custom Studio kernels — the named shelf with display metadata, ready before any container image exists. Versions land later, appended as images are..."
  - slug: "02-versioned-training-image"
    rank: "02"
    title: "Versioned Training Image"
    excerpt: "This preset registers a GPU training environment as a fully-annotated version — ECR image, compatibility metadata, and the `latest` and `stable` aliases — so training jobs select it by name instead..."
---

# SageMaker Image Presets

Ready-to-deploy configuration presets for SageMaker Image. Each preset is a complete manifest you can copy, customize, and deploy.
