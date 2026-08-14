---
title: "Presets"
description: "Ready-to-deploy configuration presets for SageMaker Pipeline"
type: "preset-list"
componentSlug: "sagemaker-pipeline"
componentTitle: "SageMaker Pipeline"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-inline-training-pipeline"
    rank: "01"
    title: "Inline Training Pipeline"
    excerpt: "This preset carries the pipeline definition inline in the manifest — the definition lives next to the rest of the spec as diffable YAML, and drift on it is visible like any other field."
  - slug: "02-s3-definition-pipeline"
    rank: "02"
    title: "S3 Definition Pipeline"
    excerpt: "This preset reads the pipeline definition from an S3 object — for definitions too large to inline, or published to S3 by the pipeline's own build tooling — with a parallelism cap on executions."
---

# SageMaker Pipeline Presets

Ready-to-deploy configuration presets for SageMaker Pipeline. Each preset is a complete manifest you can copy, customize, and deploy.
