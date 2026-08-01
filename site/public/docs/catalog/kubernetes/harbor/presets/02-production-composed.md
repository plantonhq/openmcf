---
title: "Production — composed data plane, object storage, HA components"
description: "The production posture: every stateful concern leaves the chart's evaluation-grade internals and composes the catalog's own kinds — PostgreSQL from a KubernetesPostgres (the operator-maintained..."
type: "preset"
rank: "02"
presetSlug: "02-production-composed"
componentSlug: "harbor"
componentTitle: "Harbor"
provider: "kubernetes"
icon: "package"
order: 2
---

# Production — composed data plane, object storage, HA components

The production posture: every stateful concern leaves the chart's
evaluation-grade internals and composes the catalog's own kinds —
PostgreSQL from a KubernetesPostgres (the operator-maintained
credential Secret plugs into Harbor's contract as-is), Redis from a
KubernetesValkey, and artifact blobs in S3-compatible object storage
from a KubernetesSeaweedFs (for AWS S3, drop the endpoint and omit
credentials — the registry resolves IRSA/ambient credentials from the
pod environment). Two replicas of every stateless component; an
object-storage backend is exactly what makes multi-replica registries
safe.

Exposure: a cloud LoadBalancer with TLS from a cert-manager-issued
Secret (the KubernetesCertificate reference) — `externalUrl` names
that public address, and OCI clients authenticate against exactly it.

The two bridge Secrets this preset references (`harbor-redis-auth`
with the chart's `REDIS_PASSWORD` key; `harbor-s3-auth` with
`REGISTRY_STORAGE_S3_ACCESSKEY`/`REGISTRY_STORAGE_S3_SECRETKEY`) carry
credentials the composed kinds generate under their own key shapes —
one KubernetesSecret each, or declare the values inline and let the
module materialize them.

Metrics are on with a ServiceMonitor — a KubernetesKubePrometheusStack
picks Harbor up without further wiring.
