---
title: "Presets"
description: "Ready-to-deploy configuration presets for SageMaker Model"
type: "preset-list"
componentSlug: "sagemaker-model"
componentTitle: "SageMaker Model"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-prebuilt-framework-model"
    rank: "01"
    title: "Prebuilt Framework Model"
    excerpt: "This preset serves a trained scikit-learn artifact on AWS's own prebuilt framework image — the fastest path from a `model.tar.gz` in S3 to a deployable model, with no container of your own to build."
  - slug: "02-inference-pipeline"
    rank: "02"
    title: "Inference Pipeline"
    excerpt: "This preset chains two containers into a Serial inference pipeline — raw requests hit the preprocessing container, its output feeds the scoring container, and one endpoint serves the whole chain."
---

# SageMaker Model Presets

Ready-to-deploy configuration presets for SageMaker Model. Each preset is a complete manifest you can copy, customize, and deploy.
