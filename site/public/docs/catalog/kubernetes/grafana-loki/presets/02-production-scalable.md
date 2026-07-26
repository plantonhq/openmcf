---
title: "Production scalable Loki"
description: "The write/read/backend topology on object storage — how Loki runs in production. Ingest, query and compaction scale independently, and chunks live in an S3-compatible bucket (an in-cluster..."
type: "preset"
rank: "02"
presetSlug: "02-production-scalable"
componentSlug: "grafana-loki"
componentTitle: "Grafana Loki"
provider: "kubernetes"
icon: "package"
order: 2
---

# Production scalable Loki

The write/read/backend topology on object storage — how Loki runs in
production. Ingest, query and compaction scale independently, and chunks
live in an S3-compatible bucket (an in-cluster KubernetesSeaweedFs here;
AWS S3, GCS or Azure by swapping the storage block). Retention is enforced
by the compactor.

**When to use:** any real log volume, multi-node clusters, retention that
matters.

**Cross-cloud:** on EKS drop the endpoint, set a region, and leave
credentials empty to use keyless IRSA; the same shape works with GCS
(workload identity) or Azure Blob (federated token).
