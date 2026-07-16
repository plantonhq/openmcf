---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cloud Composer User Workloads Secret"
type: "preset-list"
componentSlug: "cloud-composer-user-workloads-secret"
componentTitle: "Cloud Composer User Workloads Secret"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-airflow-connection"
    rank: "01"
    title: "Airflow Connection"
    excerpt: "A single connection URI delivered as a Kubernetes Secret into a Composer environment — the standard way DAGs reach a database without the URI living in DAG code or environment variables."
  - slug: "02-api-credentials"
    rank: "02"
    title: "API Credentials"
    excerpt: "Two token entries — an API key and a webhook signing secret — delivered as one Kubernetes Secret. DAG tasks that call the external service mount both from a single named Secret."
---

# Cloud Composer User Workloads Secret Presets

Ready-to-deploy configuration presets for Cloud Composer User Workloads Secret. Each preset is a complete manifest you can copy, customize, and deploy.
