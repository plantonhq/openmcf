# kube-prometheus-stack

Deploy [kube-prometheus-stack](https://prometheus-community.github.io/helm-charts) — the industry-standard cluster monitoring bundle — from the official `kube-prometheus-stack` Helm chart. One resource is the whole observability plane: the **Prometheus Operator** (the controller), an HA-capable **Prometheus** server on persistent volumes, an **Alertmanager**, the **kube-state-metrics** and **node-exporter** metric sources, the chart's curated Kubernetes alerting/recording rules, and a **bundled Grafana** pre-loaded with the matching dashboards.

The grain is deliberate: **one stack per cluster is the norm**. The chart installs the `monitoring.coreos.com` CRDs (ServiceMonitor, PodMonitor, PrometheusRule, Prometheus, Alertmanager, ...), which are cluster-scoped singletons — a second stack on the same cluster must set `skip_crds` and fence its discovery, an advanced posture, never the default. Every UI stays ClusterIP; external reachability composes from ingress and Gateway API kinds over the exported service handles.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Helm release** (official `kube-prometheus-stack` chart, default pin `87.19.1`, named `metadata.name`) — the operator Deployment, the Prometheus StatefulSet on a 50Gi PVC per replica (the chart's own emptyDir default is deliberately overridden), the Alertmanager StatefulSet on a 2Gi PVC, kube-state-metrics, the node-exporter DaemonSet, the curated PrometheusRule set, and the bundled Grafana with its pre-wired datasource and dashboards
- **The monitoring.coreos.com CRDs** — installed ONCE by Helm and never touched again: chart upgrades do NOT upgrade CRDs (pair version bumps across operator minors with `crd_upgrade_job`, a pre-upgrade hook that server-side-applies the new bundle), and uninstall KEEPS them, so ServiceMonitors and rules across the cluster survive removal of the stack
- **Grafana admin Secret** (`<name>-grafana`, keys `admin-user`/`admin-password`) — generated once at first install when `grafana.admin_secret` is empty, stable across upgrades, surfaced in the outputs
- **Kubernetes Namespace** — created only when `create_namespace` is true; otherwise the namespace must already exist

The modules pin the chart's fullname to `metadata.name` so child Services stay predictable (`<name>-prometheus`, `<name>-grafana`). **Keep the resource name at 26 characters or fewer** — the chart silently truncates its fullname past that, so both modules fail loudly on a longer name instead of letting the naming contract break.

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** — an active connection in the Connect module with credentials for the target cluster.

### Cluster Side

- **A StorageClass** for the Prometheus, Alertmanager, and (optional) Grafana volumes — most managed clusters provide a default; reference a **Kubernetes Storage Class** for explicit (SSD) placement. The TSDB rewards SSD-backed classes.
- **cert-manager** — only if you set `operator.admission_webhooks.cert_manager`; the chart then renders an Issuer and Certificates for the webhook instead of its self-contained certgen Job.

## Deploy

### Console

Open the deployment store, find **kube-prometheus-stack**, and click **Deploy**. The creation wizard walks you through namespace placement (with the live naming-budget warning), the chart pin and the CRD lifecycle, the Prometheus server and its storage, the discovery fence, remote-write destinations, Alertmanager, the bundled Grafana, the operator and its admission webhooks, the exporters, the managed-cloud scraper set, the air-gap path, placement, and the Helm-values escape hatch. Start from the **Managed Cloud** preset in the [Presets](#presets) tab.

### CLI

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesKubePrometheusStack
metadata:
  name: cluster-metrics
  org: acme-corp
  env: prod
spec:
  namespace:
    value: monitoring
  create_namespace: true
```

```shell
planton apply -f kube-prometheus-stack.yaml
```

This empty-spec install is a complete monitoring plane: a single Prometheus with 10-day retention on a 50Gi volume, Alertmanager, both exporters, the curated rules, the bundled Grafana with generated credentials, and cluster-wide monitor discovery.

### InfraChart

Compose the stack behind its namespace with a reference, and the InfraPipeline orders the deploys:

```yaml
spec:
  namespace:
    valueFrom:
      kind: KubernetesNamespace
      name: monitoring-namespace
      fieldPath: spec.name
  create_namespace: false
```

## Key Configuration

**Discovery is deliberately wider than the chart's own default** — unset, this component discovers EVERY ServiceMonitor, PodMonitor, PrometheusRule, Probe, and ScrapeConfig in the cluster, whoever created it. The chart's own default only discovers objects labeled by its release (upstream's most-tripped-over behavior); cluster-wide discovery is what makes every catalog component's `service_monitor_enabled` toggle and any hand-authored monitor light up with zero wiring. Set `discovery: release_managed_only` to restore the chart's fence for deliberate multi-tenant ownership boundaries.

**On managed clouds, disable the scrapers that can never succeed** — on EKS/GKE/AKS the controller-manager, etcd, and scheduler run on provider machines the cluster network can never reach; leaving their scrapers on produces permanently-down targets and alerts that can never recover (kube-proxy's metrics port is often localhost-bound too). Disable them via `control_plane_scrapers` and pair each with its rule group in `default_rules.disabled_groups`, so the alert set stays truthful. The **Managed Cloud** preset carries the exact set.

**Alertmanager's default notifies NOBODY — by design, named honestly** — an empty `alertmanager.config_yaml` routes everything to a `"null"` receiver: alerts are visible in the UIs and APIs but page no one, with the always-firing Watchdog alert routed separately as the dead-man's-switch hook. Wire real receivers (Slack, PagerDuty, email, webhooks) before trusting the stack to page anyone — and keep webhook URLs and API keys in Alertmanager's `_file` references or an AlertmanagerConfig object, never inlined in the document.

**Remote write ships samples as they are scraped** — each destination carries one auth arm: basic auth (password in an existing Secret), a bearer-token Secret, **AWS SigV4 where NO keys means the pod's ambient IRSA identity** (the recommended keyless posture for Amazon Managed Prometheus; static keys are the both-or-neither fallback), or Azure AD managed identity. Name each destination when declaring more than one — the queue-metrics identity. `enable_remote_write_receiver` is the other direction: accept pushes FROM other Prometheus servers, turning this one into an aggregation point.

**Prometheus replicas duplicate, they never shard** — each replica scrapes and stores the FULL target set independently; 2 is the HA pair, and queries against the Service may see either replica. Deduplicated long-term storage is the Thanos surface, reached through `helm_values`. Retention is dual-bounded: `retention` (default 10d) by age and `retention_size` by disk — size the latter BELOW the volume, because a full volume crash-loops instead of trimming.

**The bundled Grafana is the stack's dashboard console, not your fleet Grafana** — it arrives pre-wired with a datasource for this Prometheus and the curated dashboard set, with admin credentials generated once and exported. For a standalone, multi-datasource Grafana, disable it here and deploy **KubernetesGrafana** instead, pointed at this stack's exported Prometheus endpoint. Give the bundled one a `storage` block if people build dashboards by hand — the ephemeral default loses them with the pod.

**Nothing is exposed by default** — Prometheus, Alertmanager, and Grafana stay ClusterIP. Expose them by composing first-class kinds (KubernetesIngress, the Gateway API kinds) over the exported service handles; the stack never opens its own doors.

**`helm_values` merges last** — the escape hatch for chart surface beyond the typed fields (Thanos sidecar/ruler, windows monitoring, scrape classes, per-component securityContexts). Anything here silently overrides the typed fields on every deploy; never put secrets in it, and leave `fullnameOverride` alone — the naming contract the outputs derive from depends on it.

## Outputs and Dependencies

### What This Component Consumes

| Field | References | Purpose |
|-------|-----------|---------|
| `spec.namespace` | KubernetesNamespace (`spec.name`) | Where the stack runs |
| `spec.prometheus.storage_class` | KubernetesStorageClass (`metadata.name`) | Prometheus TSDB volume class |
| `spec.alertmanager.storage_class` | KubernetesStorageClass (`metadata.name`) | Alertmanager state volume class |
| `spec.grafana.storage.storage_class` | KubernetesStorageClass (`metadata.name`) | Grafana state volume class |
| `spec.grafana.admin_secret` | Existing Secret (username/password keys) | Bring-your-own Grafana credentials |
| `spec.prometheus.remote_write[].basic_auth.password_secret` | Existing Secret (name + key) | Remote-write basic-auth password |
| `spec.prometheus.remote_write[].bearer_token_secret` | Existing Secret (name + key) | Remote-write bearer token |
| `spec.prometheus.remote_write[].sigv4.access_key_secret` / `secret_key_secret` | Existing Secrets (name + key) | Static SigV4 keys — only when ambient IRSA identity is unavailable |
| `spec.image_pull_secrets` | Existing Secrets | Air-gap/private-mirror image pulls |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `namespace` | Namespace the stack runs in | Application deployment manifests |
| `release_name` | Helm release name (= metadata.name); every child name derives from it | Operational tooling |
| `prometheus_service` | The Prometheus Service (`<name>-prometheus`, port 9090) | ServiceMonitor-free scrape targets, diagnostics |
| `prometheus_endpoint` | In-cluster Prometheus URL | Grafana datasources, remote readers, API clients |
| `alertmanager_service` | The Alertmanager Service (`<name>-alertmanager`, port 9093). Empty when Alertmanager is disabled | External alert routing |
| `alertmanager_endpoint` | In-cluster Alertmanager URL. Empty when Alertmanager is disabled | Alert-source integrations |
| `grafana_service` | The bundled Grafana Service (`<name>-grafana`, port 80). Empty when the bundled Grafana is disabled | Ingress/Gateway exposure |
| `grafana_endpoint` | In-cluster Grafana URL. Empty when the bundled Grafana is disabled | Operator access, links |
| `grafana_admin_secret_name` | The Grafana admin-credentials Secret (`<name>-grafana`, keys `admin-user`/`admin-password`; echoes your own Secret's name when `admin_secret` is set). Empty when the bundled Grafana is disabled | Operator authentication |
| `prometheus_port_forward_command` | Copy-paste `kubectl port-forward` for the Prometheus UI | Local development access |
| `grafana_port_forward_command` | Copy-paste `kubectl port-forward` for Grafana. Empty when the bundled Grafana is disabled | Local development access |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

| Preset | Posture | What it carries |
|--------|---------|-----------------|
| **Dev Single Cluster** | kind / laptop lab / private single-node | Small volumes (10Gi/2Gi), 2-day retention, and the scrapers that cannot succeed on those platforms turned off with their rule groups paired |
| **Managed Cloud** | EKS / GKE / AKS production | 50Gi Prometheus volume with 15-day retention, durable Alertmanager and Grafana state, and every provider-internal scraper off with its rule group disabled — the target list and alert set stay honest |
| **Production Sized** | Full HA | An HA Prometheus pair, a quorum-safe 3-replica Alertmanager gossip cluster, 30-day/80GiB dual-bounded retention on 100Gi volumes, and explicit resources on every component so the stack cannot be evicted by the workloads it watches |

## Works With

- **Kubernetes Namespace** — referenced placement; the InfraPipeline orders namespace-first.
- **Kubernetes Storage Class** — SSD-backed classes for the TSDB and state volumes.
- **Kubernetes Secret** — bring-your-own Grafana credentials, remote-write credentials, image-pull Secrets.
- **Kubernetes Cert Manager** — issues the operator's admission-webhook certificate when `cert_manager` is chosen.
- **Kubernetes Ingress / Gateway API kinds** — HTTP exposure over the exported Grafana and Prometheus handles.
- **Kubernetes Grafana** — the standalone, multi-datasource alternative to the bundled Grafana, consuming the exported `prometheus_endpoint`.
- **Every catalog component with a `service_monitor_enabled` toggle** — discovered automatically under the default cluster-wide discovery.
