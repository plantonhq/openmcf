---
title: "Presets"
description: "Ready-to-deploy configuration presets for SageMaker MLflow App"
type: "preset-list"
componentSlug: "sagemaker-mlflow-app"
componentTitle: "SageMaker MLflow App"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-serverless-experiment-tracking"
    rank: "01"
    title: "Serverless Experiment Tracking"
    excerpt: "This preset is MLflow with no meter running: a serverless MLflow 3.x app storing artifacts in your S3 bucket, billed per use and $0 when idle — the successor to the hourly-billed tracking server."
  - slug: "02-studio-default-mlflow"
    rank: "02"
    title: "Studio Default MLflow"
    excerpt: "This preset makes the app the default MLflow for a SageMaker domain — Studio users in the domain track to it automatically, every logged model lands in the SageMaker Model Registry, and maintenance..."
---

# SageMaker MLflow App Presets

Ready-to-deploy configuration presets for SageMaker MLflow App. Each preset is a complete manifest you can copy, customize, and deploy.
