---
title: "Argo Workflows"
description: "Argo Workflows deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesargoworkflows"
---

# Argo Workflows

The Kubernetes-native workflow engine. Every step is a container, every
pipeline a DAG declared as a custom resource — CI jobs, data pipelines,
ML training runs and batch orchestration all run on the same engine,
with a UI that shows every step's logs and artifacts.

## Highlights

- **Engine, server, identity — typed** — controller knobs (parallelism
  caps, instance claims, retention), the Argo server with its auth
  modes, and the runner ServiceAccount pipelines execute under.
- **Artifacts anywhere** — S3-compatible (in-cluster SeaweedFS pairs by
  reference), GCS or Azure Blob; keyless cloud-identity arms on all
  three; credentials only ever referenced, never rendered.
- **History that outlives the cluster's etcd** — the workflow archive
  writes completed runs to Postgres/MySQL by reference, with retention
  keeping the live CR set lean.
- **Multi-team by instance ID** — several engines coexist in one
  cluster, each reconciling only its own labeled workflows.
- **Safe by default** — the server's `client` auth mode makes every
  caller act with their own Kubernetes permissions; CRDs are kept on
  uninstall.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
