---
title: "Grafana Loki"
description: "Grafana Loki deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesloki"
---

# Grafana Loki

Horizontally-scalable, cost-efficient log aggregation. Loki indexes only
the *labels* of your logs and stores the compressed content in object
storage, so it stays cheap at high volume — the logs half of a complete
observability stack alongside metrics (Prometheus) and traces (Tempo).

## Highlights

- **Two honest topologies** — a single-node monolithic install for dev and
  small clusters, or independently-scaling write/read/backend tiers for
  production.
- **Object storage, keyless** — S3 (including in-cluster SeaweedFS), GCS
  and Azure, using ambient cloud identity (IRSA / workload identity) or
  referenced Secrets; credentials never touch rendered config.
- **Zero-config schema** — the index schema is derived for you; no
  hand-authored `schema_config`.
- **Composes end to end** — an OpenTelemetry collector ships logs in, a
  Grafana datasource reads them back, and the ruler alerts through
  kube-prometheus-stack's Alertmanager, all by reference.
- **Multi-tenant** — one shared Loki serving many teams with gateway-
  enforced isolation.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
