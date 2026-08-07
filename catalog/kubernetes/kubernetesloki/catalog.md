# Loki

Deploy [Grafana Loki](https://grafana-community.github.io/helm-charts) — horizontally-scalable, cost-efficient log aggregation — from the official `loki` Helm chart. Loki indexes only the *labels* of your logs and stores the compressed content in object storage, so it stays cheap at high volume: the logs half of a complete observability stack alongside metrics (Prometheus) and traces (Tempo).

The grain is deliberate: **Loki stores logs; something must ship them.** Deploy a KubernetesOtelCollector in daemonset mode with the cluster-logs pipeline (its presets carry it) pointed at this component's exported `gateway_endpoint` — and read them back by pointing a KubernetesGrafana datasource of type `loki` at the same endpoint. Everything stays ClusterIP; the nginx gateway is the one front door.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Helm release** (official `loki` chart, default pin `18.5.4` — pairs with Loki 3.7.4, named `metadata.name`) — a single-replica monolithic StatefulSet on a persistent volume by default, or the topology the spec declares; plus the nginx **gateway** (the single front door routing pushes and queries in every mode), the **chunks and query-results memcached caches**, and the **canary** DaemonSet that writes and reads test log lines through the full pipeline, turning silent log loss into a visible metric
- **A derived index schema** — Loki normally requires a hand-authored `schema_config`; the modules derive it (TSDB, schema v13, the object store matching your storage backend) so a new install never writes one. The `schema_from_date` override exists solely for importing clusters whose existing schema began on a real date
- **Kubernetes Namespace** — created only when `create_namespace` is true; otherwise the namespace must already exist

The modules pin the chart's fullname to `metadata.name` so child names stay predictable (`<name>-gateway`, `<name>-backend-headless`). **Keep the resource name at 40 characters or fewer** — the chart truncates composed child names at 63 characters, so the modules fail loudly over the budget instead of letting the naming contract break.

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** — an active connection in the Connect module with credentials for the target cluster.

### Cluster Side

- **An object-storage bucket** for anything beyond a single filesystem replica — S3 (or an in-cluster **Kubernetes SeaweedFS** S3 endpoint), GCS, or Azure Blob. The bucket must exist; Loki does not create it.
- **A StorageClass** for the persistent volumes — most managed clusters provide a default.
- **kube-prometheus-stack** — only if you set `service_monitor_enabled` (the monitoring.coreos.com CRDs) or wire the ruler's `alertmanager_url` by reference.

## Deploy

### Console

Open the deployment store, find **Loki**, and click **Deploy**. The creation wizard walks you through namespace placement (with the live naming-budget warning), the chart pin, the deployment mode, object storage, retention and the import-only schema override, the ingestion/query limits, multi-tenancy, the gateway, the caches (with the live sizing warning), the ruler, self-observability, the air-gap path, placement, and the Helm-values escape hatch. Start from the **Dev Single** preset in the [Presets](#presets) tab.

### CLI

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesLoki
metadata:
  name: central-logs
  org: acme-corp
  env: prod
spec:
  namespace:
    value: logging
  create_namespace: true
  caching:
    chunks_cache_memory_mb: 256
    results_cache_memory_mb: 128
```

```shell
planton apply -f loki.yaml
```

This near-empty spec is a complete log store: a single-replica monolithic Loki on a 10Gi persistent volume, the gateway, both caches (sized down here for a small cluster — see below), the canary, and single-tenant access with no tenant header needed.

### InfraChart

Compose Loki behind its namespace with a reference, and the InfraPipeline orders the deploys:

```yaml
spec:
  namespace:
    valueFrom:
      kind: KubernetesNamespace
      name: logging-namespace
      fieldPath: spec.name
  create_namespace: false
```

## Key Configuration

**Two honest topologies — and the storage doctrine that binds them.** `monolithic` (the default when nothing is declared) runs every Loki target in one StatefulSet: right for single-node clusters, dev environments and small production volumes. `simple_scalable` splits into write/read/backend tiers that scale independently and REQUIRES object storage — the tiers rendezvous in the object store, not on a shared disk. Filesystem storage is honest ONLY for a single monolithic replica; more than one replica, or any scalable tier, requires s3/gcs/azure (mirroring the chart's own validation). The chart's microservices mode is deliberately not modeled.

**Object storage is keyless first.** On EKS leave the S3 credentials empty for ambient IRSA identity; on GKE the missing key Secret means workload identity; on AKS the missing account key means federated identity. Declared credentials are name+key REFERENCES to existing Secrets — the modules inject them as environment variables (S3) or a mounted key file (GCS); they never land in the rendered Loki config. The s3-compatible arm (`endpoint` + `force_path_style`) composes with an in-cluster KubernetesSeaweedFs. The chart's bundled MinIO subchart is deprecated by the chart itself and is never enabled by this component.

**THE CACHE-SIZING TRAP (verified live):** the caching defaults are production-scale. Unset, the chunks cache allocates 8192MB and the chart requests container memory at 1.2× — a **9830Mi request that never schedules** on a node with less than ~10Gi allocatable. The pod stays Pending and, because the install is atomic, the WHOLE release rolls back after its timeout. Set `chunks_cache_memory_mb` explicitly on any small or dev cluster (128–1024 is plenty for light query loads; the dev preset carries 256).

**Retention is off by default — deliberately visible.** Empty `retention_period` keeps everything forever (Loki's own default) and object-storage costs grow unbounded; production installs should set it (`744h` / `31d` — hours or days only). Deletion is asynchronous: the compactor marks and later sweeps, so any bucket lifecycle policy must expire LATER than the period, never earlier.

**Single-tenant by default — one line of wiring.** This component diverges from the chart's multi-tenant-on default so pushes and queries need no `X-Scope-OrgID` header. Enable `multi_tenancy` for isolation: every client then sends its tenant header, and the gateway enforces HTTP basic auth for the declared tenants (name + bcrypt htpasswd hash — one-way material, safe in a manifest; generate with `htpasswd -nbBC10`) or an existing htpasswd Secret. Basic auth AUTHENTICATES clients; each client still declares its own tenant header.

**The gateway is the one front door.** It routes pushes and queries to the right internal target in every mode, and the exported endpoints assume it — disable it only when clients address the internal services directly (single-tenant monolithic only). Expose it via KubernetesIngress or Gateway API kinds over the exported handles; Loki never opens its own doors.

**Alerting on logs** — the ruler evaluates LogQL alerting/recording rules discovered from ConfigMaps labeled `loki_rule: "1"` (the same sidecar contract Grafana dashboards use) and fires at `alertmanager_url` — a literal URL or a KubernetesKubePrometheusStack reference (its Alertmanager endpoint), the one-line wiring into the cluster's alerting.

**`helm_values` merges last** — the escape hatch for chart surface beyond the typed fields (bloom filters, the pattern ingester, zone-aware rollouts, per-component overrides, ruler storage tuning). Anything here silently overrides the typed fields on every deploy; never put secrets in it, and leave `fullnameOverride` alone — the naming contract the outputs derive from depends on it.

## Outputs and Dependencies

### What This Component Consumes

| Field | References | Purpose |
|-------|-----------|---------|
| `spec.namespace` | KubernetesNamespace (`spec.name`) | Where Loki runs |
| `spec.monolithic.storage_class` / `spec.simple_scalable.storage_class` | KubernetesStorageClass (`metadata.name`) | Volume class for the data/WAL volumes |
| `spec.storage.s3.credentials.*` | Existing Secrets (name + key) | Static S3 keys — only when ambient IRSA identity is unavailable |
| `spec.storage.gcs.service_account_key_secret` | Existing Secret (name + key) | GCS key — only when GKE workload identity is unavailable |
| `spec.storage.azure.account_key_secret` | Existing Secret (name + key) | Azure account key — only when AKS federated identity is unavailable |
| `spec.multi_tenancy.existing_htpasswd_secret` | Existing Secret (`.htpasswd` key) | Bring-your-own gateway credentials |
| `spec.ruler.alertmanager_url` | KubernetesKubePrometheusStack (`status.outputs.alertmanager_endpoint`) | Where log alerts fire |
| `spec.image_pull_secrets` | Existing Secrets | Air-gap/private-mirror image pulls |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `namespace` | Namespace Loki runs in | Application deployment manifests |
| `release_name` | Helm release name (= metadata.name); every child name derives from it | Operational tooling |
| `gateway_service` | The gateway Service (`<name>-gateway`, port 80). Empty when the gateway is disabled | Ingress/Gateway exposure |
| `gateway_endpoint` | In-cluster gateway URL — log shippers push here and Grafana `loki` datasources read here | Collector pipelines, Grafana datasources |
| `otlp_push_endpoint` | The gateway's `/otlp` route — point a KubernetesOtelCollector otlphttp exporter here | OTLP log ingest |
| `loki_service` | The Loki HTTP Service (`<name>`, port 3100) — the direct internal API behind the gateway | Diagnostics, direct API clients |
| `port_forward_command` | Copy-paste `kubectl port-forward` for the gateway | Local development access |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

| Preset | Posture | What it carries |
|--------|---------|-----------------|
| **Dev Single** | dev loop / small single-node cluster | The smallest honest Loki: one monolithic replica on a filesystem volume with the gateway — and the caches sized down explicitly (256/128MB), the one knob a small cluster cannot leave to defaults |
| **Production Scalable** | production log volume | Write/read/backend tiers (3/3/3) on an in-cluster SeaweedFS S3 bucket with credential references, 31-day compactor-enforced retention, and raised ingest limits |
| **Multitenant** | a shared logging service | The scalable shape plus tenant isolation: every push/query carries `X-Scope-OrgID`, and the gateway enforces per-team basic auth from bcrypt htpasswd hashes — no Secret fixture needed |

## Works With

- **Kubernetes OTel Collector** — ships logs in: daemonset mode with the cluster-logs pipeline pointed at the exported `gateway_endpoint` / `otlp_push_endpoint`.
- **Kubernetes Grafana** — reads logs back: a datasource of type `loki` at the gateway endpoint.
- **Kubernetes kube-prometheus-stack** — scrapes Loki's ServiceMonitors and receives the ruler's log alerts by reference.
- **Kubernetes SeaweedFS** — in-cluster S3-compatible object storage for the s3 arm (`endpoint` + `force_path_style`).
- **Kubernetes Namespace** — referenced placement; the InfraPipeline orders namespace-first.
- **Kubernetes Storage Class** — explicit volume classes for the data/WAL volumes.
- **Kubernetes Ingress / Gateway API kinds** — HTTP exposure over the exported gateway Service handle.
