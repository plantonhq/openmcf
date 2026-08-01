---
title: "OTLP fan-in preset"
description: "One front door for everything: applications push OTLP (gRPC on 4317, HTTP on 4318 — the exported endpoints), and the collector fans out by signal — logs to an `otlphttp` backend, traces to an `otlp`..."
type: "preset"
rank: "03"
presetSlug: "03-otlp-fanin"
componentSlug: "opentelemetry-collector"
componentTitle: "OpenTelemetry Collector"
provider: "kubernetes"
icon: "package"
order: 3
---

# OTLP fan-in preset

One front door for everything: applications push OTLP (gRPC on 4317,
HTTP on 4318 — the exported endpoints), and the collector fans out by
signal — logs to an `otlphttp` backend, traces to an `otlp` gRPC
backend — through the `memory_limiter` and `batch` processors. Apps
configure ONE endpoint; where telemetry lands is this resource's
decision, changed here without touching a single application.

PREREQUISITE: a `KubernetesOtelOperator` on the cluster. No custom
ServiceAccount is needed — this pipeline reads no cluster state.

Instead of a fixed replica count, the operator manages an HPA: two to
five replicas targeting 75% average CPU. The autoscaler owns the
count — the spec's validation rejects a non-default `replicas`
alongside it.

Change first: both exporter hosts — the `otlphttp` endpoint from your
`KubernetesLoki` gateway (Loki ingests OTLP at its `/otlp` route), the
`otlp` endpoint from your Tempo install's OTLP gRPC endpoint. Then add
`resources` so the CPU target has real requests to measure against,
keeping the memory limit in step with the `memory_limiter`'s 400 MiB.

See [03-otlp-fanin.yaml](./03-otlp-fanin.yaml) for the manifest.
