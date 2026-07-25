# Kubernetes SigNoz

## When NOT to Use This

**One resource is ONE SigNoz install** — the all-in-one open-source
observability platform: traces, metrics and logs in one UI, stored in
ClickHouse, ingested over OpenTelemetry.

Not the right component when:

- **You want the composed best-of-breed stack** — that is
  `KubernetesKubePrometheusStack` + `KubernetesGrafana` +
  `KubernetesLoki` + `KubernetesTempo`. SigNoz is the "one product, one
  UI" alternative; both paths are first-class.
- **You need something to SHIP cluster telemetry** — SigNoz receives
  OTLP; it does not scrape nodes or tail container logs by itself.
  Deploy a `KubernetesOtelCollector` pointed at this component's
  exported `otlp_grpc_endpoint` / `otlp_http_endpoint`.
- **You expect a public endpoint out of the box** — everything is
  ClusterIP; exposure composes from first-class kinds
  (`KubernetesIngress`, the Gateway API kinds) over the exported
  service handles.

## Database arms

- **Bundled** (default) — the chart deploys and owns a ClickHouse stack:
  a namespace-fenced Altinity operator, the installation, and ZooKeeper.
  The appliance posture: capacity and topology knobs only. The admin
  password is module-GENERATED per install (the chart's
  publicly-documented default never ships) and exported through the
  `<name>-clickhouse-auth` Secret.
- **External** — bring your own ClickHouse; nothing database-related
  installs. The fields default-reference a `KubernetesClickHouse`
  resource's outputs (client Service, cluster name, auth Secret), so the
  composition is one `valueFrom` per field. Declare the SigNoz user
  there with `access_management` and cluster-wide DDL grants — SigNoz
  creates and migrates its own databases.

**Co-existence warning (bundled arm):** the bundled Altinity operator is
namespace-fenced, and the clickhouse.altinity.com CRDs install only when
absent and are KEPT on uninstall. On a cluster that also runs
`KubernetesAltinityOperator`, keep this component in its own namespace
and that operator's watch fenced away from it.

## Single-instance truth

The community SigNoz server keeps users, dashboards and alert rules in
SQLite on a PersistentVolume — a single-writer store, so the server runs
exactly one replica (the Postgres-backed HA store is enterprise-only
upstream). Telemetry lives in ClickHouse and scales independently; the
ingestion collector scales horizontally (replicas or the HPA arm).

## Retention

Telemetry retention is managed inside SigNoz (UI → Settings → General),
not in this spec — size the bundled arm's `disk_size` for the retention
you set there.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
