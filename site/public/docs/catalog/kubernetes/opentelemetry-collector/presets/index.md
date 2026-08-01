---
title: "Presets"
description: "Ready-to-deploy configuration presets for OpenTelemetry Collector"
type: "preset-list"
componentSlug: "opentelemetry-collector"
componentTitle: "OpenTelemetry Collector"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-cluster-logs-to-loki"
    rank: "01"
    title: "Cluster-logs-to-Loki preset"
    excerpt: "The per-node log pipeline: daemonset mode puts one collector on every node, the filelog receiver tails every container's log files under `/var/log/pods` (the standard `container` operator parses the..."
  - slug: "02-traces-gateway-to-tempo"
    rank: "02"
    title: "Traces-gateway-to-Tempo preset"
    excerpt: "The gateway shape for traces: a two-replica deployment applications push OTLP spans to (gRPC on 4317, HTTP on 4318 — the exported `otlp_grpc_endpoint`/`otlp_http_endpoint` outputs are what they point..."
  - slug: "03-otlp-fanin"
    rank: "03"
    title: "OTLP fan-in preset"
    excerpt: "One front door for everything: applications push OTLP (gRPC on 4317, HTTP on 4318 — the exported endpoints), and the collector fans out by signal — logs to an `otlphttp` backend, traces to an `otlp`..."
---

# OpenTelemetry Collector Presets

Ready-to-deploy configuration presets for OpenTelemetry Collector. Each preset is a complete manifest you can copy, customize, and deploy.
