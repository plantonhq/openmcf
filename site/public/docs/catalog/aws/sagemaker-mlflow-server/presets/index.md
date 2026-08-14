---
title: "Presets"
description: "Ready-to-deploy configuration presets for SageMaker MLflow Server"
type: "preset-list"
componentSlug: "sagemaker-mlflow-server"
componentTitle: "SageMaker MLflow Server"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-team-experiment-tracking"
    rank: "01"
    title: "Team Experiment Tracking"
    excerpt: "This preset stands up a `Small` tracking server for a single ML team — experiments, runs, and model tracking with artifacts in your S3 bucket, and automatic model registration left off (the safe..."
  - slug: "02-regulated-experiment-tracking"
    rank: "02"
    title: "Regulated Experiment Tracking"
    excerpt: "This preset pins down every variable a governed ML platform cares about: a `Medium` server on a pinned MLflow version, maintenance in a declared quiet window, and every logged model auto-registered..."
---

# SageMaker MLflow Server Presets

Ready-to-deploy configuration presets for SageMaker MLflow Server. Each preset is a complete manifest you can copy, customize, and deploy.
