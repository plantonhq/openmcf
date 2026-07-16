---
title: "Dev Jupyter"
description: "A lightweight development cluster with JupyterLab for interactive data exploration and prototyping Spark jobs, wired to delete itself after 30 minutes of inactivity."
type: "preset"
rank: "01"
presetSlug: "01-dev-jupyter"
componentSlug: "dataproc-cluster"
componentTitle: "Dataproc Cluster"
provider: "gcp"
icon: "package"
order: 1
---

# Dev Jupyter

A lightweight development cluster with JupyterLab for interactive data
exploration and prototyping Spark jobs, wired to delete itself after 30
minutes of inactivity.

## When to use

- Interactive Spark development and debugging
- Data exploration with Jupyter notebooks
- Prototyping ML pipelines before production deployment
- Learning and experimentation with Spark/Hadoop

## What to customize

- `projectId` — your GCP project ID.
- `region` — where the cluster's nodes live.
- `lifecycleConfig.idleDeleteTtl` — extend beyond `1800s` if notebook
  sessions routinely idle longer than 30 minutes between cells.
- Machine types — `e2-standard-4` keeps dev costs low; move to `n2`
  machines when notebook workloads grow.

## Key configuration

- **1 master + 2 workers** — the smallest useful Spark topology
- **JUPYTER component** — JupyterLab reachable through the Component
  Gateway (no SSH tunnels)
- **Component Gateway enabled** — authenticated web access to the Spark
  UI, YARN ResourceManager, and HDFS NameNode
- **30-minute idle auto-delete** — the cost-control lever for dev
  clusters; both lifecycle TTLs update in place, so tuning it later
  never recreates the cluster

## Related presets

- **02-ha-production** — high-availability cluster for production workloads
- **03-cost-optimized-batch** — Spot secondaries and autoscaling for batch jobs
- **04-spark-on-gke** — run Spark as pods on an existing GKE cluster
