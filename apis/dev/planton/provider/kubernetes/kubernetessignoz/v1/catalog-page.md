# SigNoz

The open-source, OpenTelemetry-native observability platform: traces,
metrics and logs in ONE application with one UI, stored in ClickHouse —
the one-component alternative to composing a metrics stack, a log store,
a trace store and a dashboard tool separately.

## Highlights

- **The whole platform from one manifest** — UI, API, alerting, the
  ingestion collector and (by default) a bundled ClickHouse stack, with
  a module-generated database credential exported as a Secret handle.
- **Bring your own ClickHouse** — point it at a `KubernetesClickHouse`
  by reference (Service, cluster name, auth Secret) when the database
  deserves its own lifecycle.
- **OTLP-first ingestion** — gRPC and HTTP endpoints exported as
  composition handles; Jaeger and Zipkin protocols for legacy emitters;
  the collector autoscales with ingest volume.
- **Alerting built in** — rule evaluation and notification in the server
  itself, with SMTP over secret-safe wiring for email delivery.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
