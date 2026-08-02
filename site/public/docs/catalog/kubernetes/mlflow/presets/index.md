---
title: "Presets"
description: "Ready-to-deploy configuration presets for MLflow"
type: "preset-list"
componentSlug: "mlflow"
componentTitle: "MLflow"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-team-tracking"
    rank: "01"
    title: "Team tracking server — zero dependencies"
    excerpt: "The one-manifest MLflow: experiments, runs, metrics and the model registry on a sqlite volume, artifacts on a second volume served through the tracking server itself — nothing else to operate, and..."
  - slug: "02-production-postgres-s3"
    rank: "02"
    title: "Production — PostgreSQL backend, S3 artifacts"
    excerpt: "The durable, team-scale shape: experiments, runs, metrics and the model registry in a composed KubernetesPostgres; models and datasets in a composed KubernetesSeaweedFs bucket; two stateless server..."
---

# MLflow Presets

Ready-to-deploy configuration presets for MLflow. Each preset is a complete manifest you can copy, customize, and deploy.
