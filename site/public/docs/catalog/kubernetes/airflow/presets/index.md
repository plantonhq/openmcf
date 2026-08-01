---
title: "Presets"
description: "Ready-to-deploy configuration presets for Airflow"
type: "preset-list"
componentSlug: "airflow"
componentTitle: "Airflow"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-dev-kubernetes-executor"
    rank: "01"
    title: "Dev preset — KubernetesExecutor"
    excerpt: "The smallest useful Airflow 3: API server (UI + REST API), scheduler, DAG processor and triggerer against a composed KubernetesPostgres named `airflow-db` in the same Kubernetes namespace, running..."
  - slug: "02-production-celery"
    rank: "02"
    title: "Production preset — Celery + git-sync + KEDA"
    excerpt: "The production shape: a Celery worker fleet that KEDA scales on real queue depth (polling the metadata database for queued tasks — and back to ZERO between runs), DAGs synced from your Git repository..."
---

# Airflow Presets

Ready-to-deploy configuration presets for Airflow. Each preset is a complete manifest you can copy, customize, and deploy.
