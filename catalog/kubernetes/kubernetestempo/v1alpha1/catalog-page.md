# Grafana Tempo

Cost-efficient distributed tracing. Tempo stores whole traces in object
storage and needs no expensive index — you retrieve by trace ID or TraceQL
— so it scales to high span volumes cheaply. The traces half of a complete
observability stack alongside metrics (Prometheus) and logs (Loki).

## Highlights

- **Object storage, keyless** — S3 (including in-cluster SeaweedFS), GCS
  and Azure, using ambient cloud identity (IRSA / workload identity) or
  referenced Secrets; credentials never touch rendered config.
- **Persistent by default** — traces survive pod restarts (the upstream
  chart's emptyDir default would lose them).
- **OTLP-first** — the modern wire standard on by default; the legacy
  Jaeger protocols available for fleets still migrating.
- **Service map from traces** — the metrics generator derives
  service-graph and span metrics and remote-writes them to
  kube-prometheus-stack, lighting up Grafana's service map.
- **Composes end to end** — an OpenTelemetry collector sends spans in, a
  Grafana datasource reads them back, all by reference.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
