# Kubernetes Kube Prometheus Stack

Deploys the kube-prometheus-stack — the industry-standard cluster
monitoring bundle — from the official prometheus-community Helm
chart. One install gives the Prometheus Operator, a Prometheus server
declared through it, an Alertmanager, kube-state-metrics and
node-exporter, the curated Kubernetes alerting/recording rule set,
and a bundled Grafana pre-loaded with the matching dashboards. By
default this component discovers every ServiceMonitor, PodMonitor,
PrometheusRule, Probe and ScrapeConfig in the cluster — which is what
makes every other catalog component's `service_monitor_enabled`
toggle light up with zero extra wiring — and provisions persistent
volumes for Prometheus and Alertmanager (the chart's own default is
an emptyDir that loses everything on restart).

> **Why a monitoring stack**: metrics are how a cluster tells you it
> is unwell before users do. Prometheus scrapes and stores them, the
> rule set turns them into alerts, Alertmanager routes those alerts
> to humans, and Grafana renders the dashboards — one composed
> install instead of four hand-wired ones.

## What Gets Created

- **Namespace** (optional) — created and owned when
  `create_namespace` is set
- **Helm release** (official `kube-prometheus-stack` chart, pinned
  87.19.1, named `metadata.name`; the modules pin the chart fullname
  to the resource name, so keep it at 26 characters or fewer — both
  engines fail loudly beyond that):
  - the **Prometheus Operator** Deployment (`<name>-operator`) with
    its admission webhooks (certgen hook Job by default; cert-manager
    arm available)
  - a **Prometheus** StatefulSet (`prometheus-<name>-prometheus`)
    with a persistent TSDB volume per replica
  - an **Alertmanager** StatefulSet
    (`alertmanager-<name>-alertmanager`) with a persistent state
    volume
  - the bundled **Grafana** Deployment (`<name>-grafana`) with the
    stack's dashboard set and a datasource pre-wired to this
    Prometheus
  - **kube-state-metrics** (Deployment) and **node-exporter**
    (DaemonSet)
  - the curated **PrometheusRule** set
- **monitoring.coreos.com CRDs** — installed once via the chart's
  crds subchart; never upgraded by chart bumps and KEPT on uninstall
  (monitors and rules across the cluster survive removal of the
  stack)
- **Remote-write auth Secret** (`<name>-remote-write-auth`, when any
  remote-write destination declares basic auth) — module-owned home
  of the usernames; passwords stay in the Secrets you reference

## Prerequisites

- A Kubernetes namespace that already exists, or set
  `create_namespace`
- A StorageClass for the volumes (most managed clusters provide a
  default)
- For `remote_write` with basic auth / bearer token / static SigV4
  keys: the referenced password/token/key Secrets must exist in the
  namespace
- For `operator.admission_webhooks.cert_manager`: cert-manager on the
  cluster
- For a keyless SigV4 arm (Amazon Managed Prometheus): pod identity
  (IRSA) on the cluster's service account

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesKubePrometheusStack
metadata:
  name: metrics
spec:
  namespace:
    value: monitoring
  create_namespace: true
```

That single manifest deploys the full stack: persistent Prometheus
(50Gi, 10-day retention), Alertmanager, bundled Grafana (credentials
in the `metrics-grafana` Secret, name exported in the outputs), both
exporters, the curated rules, and cluster-wide monitor discovery.
Port-forward commands for Prometheus and Grafana land in the stack
outputs.

## Configuration

### Prometheus

`prometheus.replicas` is HA duplication (each replica scrapes the
full target set independently), not sharding — 2 is the HA pair;
Thanos-based deduplicated long-term storage rides `helm_values`.
Retention ends when either `retention` (time) or `retention_size`
(bytes — keep it below `disk_size`) is reached. Memory scales with
active series; undersized limits are the most common cause of a
crash-looping Prometheus. Set `external_labels.cluster` so
multi-cluster backends can tell series apart.

### Discovery

The component default is cluster-wide: every monitor/rule object in
the cluster is discovered, whoever created it. That is deliberately
wider than the chart's own release-fenced default — upstream's
most-tripped-over behavior ("my ServiceMonitor is ignored").
`discovery: release_managed_only` restores the fence for multi-tenant
clusters running several Prometheus servers.

### Alerting

Alerts flow without any configuration — the curated rules evaluate
and Alertmanager receives them, routing everything to a "null"
receiver until `alertmanager.config_yaml` declares real destinations
(Slack, PagerDuty, email, webhooks). Webhook URLs and API keys in
that document are credentials: reference them with Alertmanager's
`_file` fields and mount the Secret via `helm_values`, never inline.

### Managed clouds

On EKS/GKE/AKS the controller-manager, scheduler and etcd are
provider-internal — their scrape targets can never be reached.
Disable them via `control_plane_scrapers` and silence the matching
rule groups via `default_rules.disabled_groups` (the presets carry
the exact set). The API server and kubelets stay scraped everywhere.

### Remote write

Ship samples to long-term/managed backends: basic auth (Grafana
Cloud/Mimir — username materializes in the module-owned
`<name>-remote-write-auth` Secret, password rides your Secret),
bearer token, AWS SigV4 (keyless via IRSA, or static keys), Azure AD
managed identity. `enable_remote_write_receiver` turns this
Prometheus into an aggregation point for other agents.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the stack runs in |
| `release_name` | Helm release name (= `metadata.name`) |
| `prometheus_service` | Prometheus Service name (port 9090) |
| `prometheus_endpoint` | In-cluster Prometheus URL — the datasource/API handle |
| `alertmanager_service` / `alertmanager_endpoint` | Alertmanager handles (empty when disabled) |
| `grafana_service` / `grafana_endpoint` | Bundled-Grafana handles (empty when disabled) |
| `grafana_admin_secret_name` | Bundled Grafana's credentials Secret (keys `admin-user`/`admin-password`) |
| `prometheus_port_forward_command` / `grafana_port_forward_command` | Workstation access without composed exposure |

## Related Components

- [KubernetesGrafana](/docs/catalog/kubernetes/grafana) — the
  standalone composition hub; references this stack's
  `prometheus_endpoint` as a datasource
- [KubernetesNamespace](/docs/catalog/kubernetes/kubernetesnamespace)
  — provides the target namespace via reference
- [KubernetesStorageClass](/docs/catalog/kubernetes/kubernetesstorageclass)
  — backs the Prometheus/Alertmanager/Grafana volumes via reference
- [KubernetesCertManager](/docs/catalog/kubernetes/kubernetescertmanager)
  — issues the admission-webhook certificate in the cert-manager arm
- [KubernetesIngress](/docs/catalog/kubernetes/kubernetesingress) —
  composes exposure over the Grafana/Prometheus service handles

## Next Steps

Deploy the stack once per cluster, then let composition do the work:
flip `service_monitor_enabled` on the components you already run and
their metrics appear without further wiring. Declare
`alertmanager.config_yaml` the moment alerts should reach a human.
On managed clouds start from the managed-cloud preset so the target
list is honest from day one, and set
`prometheus.external_labels.cluster` before the first remote-write
destination is added.
