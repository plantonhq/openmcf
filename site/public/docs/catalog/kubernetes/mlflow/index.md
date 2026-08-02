---
title: "MLflow"
description: "MLflow deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesmlflow"
---

# MLflow

The experiment tracker and model registry. Every training run's
parameters, metrics, artifacts and models — logged, compared, versioned
and staged from one server your whole team points MLFLOW_TRACKING_URI
at.

## Highlights

- **The official image, a clean deployment** — MLflow ships no Helm
  chart, so this component renders its own manifests around
  ghcr.io/mlflow/mlflow (the `-full` variant, which carries the
  database drivers, object-store clients and auth dependencies the
  bare image omits): server, credentials, volumes, optional garbage
  collection and metrics.
- **Secured by default** — upstream's open server and its
  admin/password1234 example never ship; basic auth is on with a
  generated admin password, and experiments can be private-by-default.
- **Composable state** — tracking data in a composed
  KubernetesPostgres, artifacts in a composed KubernetesSeaweedFs
  bucket (or AWS S3/GCS/Azure Blob) — or both on PVCs for the
  zero-dependency start.
- **Credential-free clients** — the server proxies artifact traffic,
  so notebooks and pipelines need only the tracking URI and their
  login; storage keys never leave the server.

## Operational notes

Deleted runs stay restorable until the garbage-collection CronJob
(off by default) hard-deletes them past the retention window. Scale
beyond one replica only on the postgres + object-store shape — the
sqlite and PVC arms are single-writer by nature (enforced at
validation).
