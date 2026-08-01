---
title: "OpenTelemetry Collector"
description: "OpenTelemetry Collector deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesotelcollector"
---

# OpenTelemetry Collector

The vendor-neutral telemetry pipeline: receives logs, metrics and
traces in dozens of protocols, processes them, and exports them to any
backend. This component declares the `OpenTelemetryCollector` CR the
OpenTelemetry Operator reconciles into the collector workload, its
Services and its rendered config — the pipeline document itself is the
product (a `KubernetesOtelOperator` on the cluster is the
prerequisite).

## Highlights

- **The pipeline is the product** — `config_yaml` carries the
  collector's own open configuration contract (receivers → processors
  → exporters); the OpenTelemetry component registry is unbounded by
  design, so the upstream grain is the right grain.
- **Four modes, one field** — a scalable gateway (deployment), per-node
  collection for logs and host metrics (daemonset), stable identities
  (statefulset), or operator-injected sidecars; mode-impossible fields
  are validation errors at apply time.
- **Credentials never touch the config** — Secrets load as environment
  variables (`env_from_secrets`) and the config references them as
  `${env:VAR_NAME}`; nothing secret-bearing lands in the rendered
  ConfigMap.
- **The operator fills the gaps** — it injects the default collector
  image when none is declared and derives Service ports from the
  declared receivers; extra ports exist only for receivers it cannot
  infer.
- **Composes end to end** — presets ship cluster logs to a
  `KubernetesLoki`, traces to Tempo, and OTLP fan-in to both; the
  exported OTLP endpoints are what applications and sibling kinds
  point at.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
