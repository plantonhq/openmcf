---
title: "Tempo"
description: "Tempo deployment documentation"
icon: "package"
order: 100
componentName: "kubernetestempo"
---

# Tempo

Deploy [Grafana Tempo](https://grafana.com/oss/tempo/) — cost-efficient distributed tracing — from the official `tempo` Helm chart. Tempo stores whole traces in object storage and needs no expensive index (you retrieve by trace ID or TraceQL), so it scales to high span volumes cheaply. It is the traces half of a complete observability stack alongside metrics (Prometheus) and logs (Loki).

The grain is deliberate: **this kind is the single-binary Tempo** — one StatefulSet, production-capable with an object-storage backend; the `tempo-distributed` microservices chart is deliberately not modeled. Spans arrive over OTLP (from applications or an OpenTelemetry collector at the exported ingest endpoints) and Grafana reads them back (a datasource of type `tempo` at the exported HTTP endpoint). Everything stays ClusterIP; external reachability composes from ingress and Gateway API kinds over the exported Service handle.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Helm release** (official grafana-community `tempo` chart, default pin `2.2.3` pairing with Tempo `2.10.7`, named `metadata.name`) — one Tempo StatefulSet with OTLP receivers on gRPC 4317 and HTTP 4318 always on
- **A PersistentVolumeClaim** (default 10Gi on the cluster's default StorageClass) — the chart's own default is an emptyDir that loses every trace on pod restart; this component deliberately inverts that. With local storage the volume holds ALL trace blocks; with an object-storage backend it holds only the write-ahead log. `ephemeral: true` restores the chart's throwaway posture
- **Object-storage wiring** when a backend is declared — S3/S3-compatible (bucket AND endpoint required; Tempo never derives the endpoint from the region), GCS, or Azure Blob. Empty credentials mean the pod's ambient identity (IRSA on EKS, workload identity on GKE, federated identity on AKS — the recommended keyless postures); declared credentials are references to existing Secrets, injected as environment variables and never rendered into config
- **Kubernetes Namespace** — created only when `create_namespace` is true; otherwise the namespace must already exist

The modules pin the chart's fullname to `metadata.name` so the exported Service and endpoints stay predictable. **Keep the resource name at 45 characters or fewer** — the chart truncates composed child names at 63 characters, so both modules fail loudly on a longer name instead of letting the naming contract break.

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** — an active connection in the Connect module with credentials for the target cluster.

### Cluster Side

- **A StorageClass** for the data volume (unless `ephemeral`) — most managed clusters provide a default; reference a **Kubernetes Storage Class** for explicit (SSD) placement. Block compaction rewards SSD-backed classes.
- **An object-storage bucket** when a backend is declared — the bucket must exist; Tempo does not create it. An in-cluster **Kubernetes SeaweedFS** S3 endpoint works with `force_path_style: true`.
- **kube-prometheus-stack** — only if you enable `service_monitor_enabled` (the monitoring.coreos.com CRDs) or point the metrics generator at its Prometheus (which must set `prometheus.enable_remote_write_receiver: true`).

## Deploy

### Console

Open the deployment store, find **Tempo**, and click **Deploy**. The creation wizard walks you through namespace placement (with the live naming-budget warning), the chart pin, topology and its volume, the trace-storage backend, retention, the ingest surface, the metrics generator, the query and self-observability dials, sizing, the air-gap path, placement, and the Helm-values escape hatch. Start from the **Production S3** preset in the [Presets](#presets) tab.

### CLI

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesTempo
metadata:
  name: traces
  org: acme-corp
  env: prod
spec:
  namespace:
    value: tracing
  create_namespace: true
```

```shell
planton apply -f tempo.yaml
```

This empty-spec install is a complete trace store: one replica on a 10Gi persistent volume, OTLP receivers on 4317/4318, a 24-hour retention window, single-tenant, and no data leaving the cluster.

### InfraChart

Compose Tempo behind its namespace with a reference, and the InfraPipeline orders the deploys:

```yaml
spec:
  namespace:
    valueFrom:
      kind: KubernetesNamespace
      name: tracing-namespace
      fieldPath: spec.name
  create_namespace: false
```

## Key Configuration

**Replicas above 1 require an object-storage backend** — the spec enforces it: replicas cannot share local trace storage. With S3, GCS, or Azure declared, replicas share the backend and scale ingest and query; without one, a single replica is the honest ceiling.

**Persistent by default, ephemeral by choice** — the chart's own emptyDir default loses every trace on pod restart, so this component provisions a PVC unless you opt out. `ephemeral: true` excludes a custom `disk_size` and a `storage_class` (the platform-stamped 10Gi default is tolerated) — there is no volume for them to configure.

**Retention speaks Go durations — there is no day unit** — `retention` accepts minutes or hours only (`30m`, `24h`; a week is `168h`, never `7d`). The chart default of 24h suits a dev loop; raise it for anything users depend on. Longer retention costs volume capacity with local storage and object-store bytes with a backend.

**OTLP-first ingest** — gRPC 4317 and HTTP 4318 are always on; they are the 2026 wire standard. `jaeger_receivers_enabled` opens the four legacy Jaeger protocols (gRPC 14250, thrift-binary 6832, thrift-compact 6831, thrift-http 14268) for fleets still migrating — the component deliberately narrows the chart's all-receivers default, because every closed port is one less ingest surface.

**Multi-tenancy is a header contract** — with `multi_tenancy_enabled`, an `X-Scope-OrgID` tenant header is required on every push AND every query; senders and the Grafana datasource must both carry it, or queries return empty results rather than errors.

**The metrics generator lights up Grafana's service map** — it derives service-graph and span metrics from the trace stream and remote-writes them to a Prometheus. The URL accepts a literal or a reference to a **kube-prometheus-stack** (its `prometheus_endpoint` output); the target Prometheus must accept pushes (`prometheus.enable_remote_write_receiver: true` on the stack), and when the URL carries no path the modules append the standard `/api/v1/write`. An empty `processors` list runs both `service_graphs` and `span_metrics` — Tempo's own default set.

**Grafana is the query surface** — `tempo_query_enabled` adds the Jaeger-UI-compatible query sidecar on 16686 only for tooling that speaks the Jaeger API. `usage_reporting` is the component's privacy-first divergence from Tempo's report-by-default: no anonymous statistics leave the cluster without an explicit opt-in.

**`helm_values` merges last** — the escape hatch for chart surface beyond the typed fields (per-receiver tuning, tenant overrides, search concurrency). Anything here silently overrides the typed fields on every deploy; never put secrets in it (object-storage credentials belong in the typed Secret-reference fields), and leave `fullnameOverride` alone — the naming contract the outputs derive from depends on it.

## Outputs and Dependencies

### What This Component Consumes

| Field | References | Purpose |
|-------|-----------|---------|
| `spec.namespace` | KubernetesNamespace (`spec.name`) | Where Tempo runs |
| `spec.storage_class` | KubernetesStorageClass (`metadata.name`) | The data-volume class |
| `spec.storage.s3.credentials` | Existing Secrets (name + key) | Static S3 keys — only when ambient IRSA identity is unavailable |
| `spec.storage.gcs.service_account_key_secret` | Existing Secret (name + key) | GCS service-account key — only when workload identity is unavailable |
| `spec.storage.azure.account_key_secret` | Existing Secret (name + key) | Azure account key — only when federated identity is unavailable |
| `spec.metrics_generator.remote_write_url` | KubernetesKubePrometheusStack (`status.outputs.prometheus_endpoint`) | Where derived metrics are remote-written |
| `spec.image_pull_secrets` | Existing Secrets | Air-gap/private-mirror image pulls |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `namespace` | Namespace Tempo runs in | Application deployment manifests |
| `release_name` | Helm release name (= metadata.name); the Service name derives from it | Operational tooling |
| `service` | The Tempo Service (all ports) | Ingress/Gateway exposure |
| `http_endpoint` | In-cluster HTTP endpoint (port 3200) | Grafana `tempo` datasources, TraceQL clients |
| `otlp_grpc_endpoint` | In-cluster OTLP gRPC trace-ingest endpoint (port 4317) | Application SDKs, OpenTelemetry collector `otlp` exporters |
| `otlp_http_endpoint` | In-cluster OTLP HTTP trace-ingest endpoint (port 4318) | OTLP/HTTP senders |
| `port_forward_command` | Copy-paste `kubectl port-forward` for the Tempo API | Local development access |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

| Preset | Posture | What it carries |
|--------|---------|-----------------|
| **Dev Single** | local development, demos, single-node clusters | One monolithic replica on a persistent volume with OTLP receivers — the smallest honest trace store; the component defaults carry the whole posture |
| **Production S3** | production tracing at scale | Two replicas against an S3-compatible backend (an in-cluster SeaweedFS; AWS S3, GCS or Azure by swapping the block), a two-week retention window, and the metrics generator remote-writing to the cluster's Prometheus |
| **Jaeger Compat** | migration from Jaeger instrumentation | The four legacy Jaeger receiver protocols opened alongside OTLP and the Jaeger-UI-compatible query sidecar — a bridge posture; disable the receivers once emitters are on OTLP |

## Works With

- **Kubernetes Namespace** — referenced placement; the InfraPipeline orders namespace-first.
- **Kubernetes Storage Class** — SSD-backed classes for the data volume.
- **Kubernetes Secret** — object-storage credentials and image-pull Secrets, always by reference.
- **Kubernetes SeaweedFS** — the in-cluster S3-compatible backend for the storage block.
- **Kubernetes Otel Collector** — sends spans to the exported OTLP endpoints from a traces pipeline.
- **Kubernetes kube-prometheus-stack** — receives the metrics generator's remote-write (the service map's data) and provides the ServiceMonitor CRDs.
- **Kubernetes Grafana** — reads traces back through a `tempo` datasource at the exported HTTP endpoint.
- **Kubernetes Ingress / Gateway API kinds** — HTTP exposure over the exported Service handle.
