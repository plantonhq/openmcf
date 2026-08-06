---
title: "Grafana"
description: "Grafana deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesgrafana"
---

# Kubernetes Grafana

Deploys a standalone Grafana — the observability dashboard — from the
official Grafana Helm chart. This kind is the composition hub: point
it at any mix of datasources (a KubernetesKubePrometheusStack's
Prometheus, Loki, Tempo, ClickHouse, Postgres) and run it
independently of all of them. Admin credentials are chart-generated
once into a Kubernetes Secret (or read from a Secret you own) and
never appear in rendered values. State is a deliberate choice:
ephemeral by default, a persistent volume for a single stateful
instance, or an external database — which is also the requirement for
running more than one replica. The Service stays ClusterIP; exposure
composes from ingress and gateway kinds.

> **When the bundled Grafana is enough**: if all you need is
> dashboards for one kube-prometheus-stack, its built-in Grafana
> (on by default there, pre-wired with datasource and dashboards) is
> the simpler path. Deploy this kind when Grafana outlives any single
> stack.

## What Gets Created

- **Namespace** (optional) — created and owned when
  `create_namespace` is set
- **Helm release** (official `grafana` chart from
  grafana-community/helm-charts, pinned 12.8.0 = Grafana 13.1.1,
  named `metadata.name`): the Grafana Deployment, the ClusterIP
  Service (port 80 → container 3000), the provisioning ConfigMaps for
  declared datasources and dashboards, the dashboard-discovery
  sidecar (on by default), and a PersistentVolumeClaim when `storage`
  is declared
- **Admin Secret** (`<name>`, keys `admin-user` / `admin-password`) —
  chart-owned; the password generates ONCE at first install and stays
  stable across upgrades. When `admin_secret` references a Secret you
  own, that Secret is used instead and its name is echoed in the
  outputs

## Prerequisites

- A Kubernetes namespace that already exists, or set
  `create_namespace`
- For `storage`: a StorageClass (most managed clusters provide a
  default; or reference a KubernetesStorageClass)
- For `database`: a reachable Postgres/MySQL with the database
  already created, and its password in an existing Secret
- For `admin_secret` / datasource `basic_auth` / `smtp` credentials:
  the referenced Secrets must exist before the install
- For `service_monitor_enabled`: the Prometheus Operator CRDs (deploy
  KubernetesKubePrometheusStack first)
- For `plugins` and `community_dashboards`: internet egress at pod
  start (both download at startup)

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesGrafana
metadata:
  name: dashboards
spec:
  namespace:
    value: observability
  create_namespace: true
  datasources:
    - name: Prometheus
      url:
        value: http://cluster-metrics-prometheus.monitoring.svc.cluster.local:9090
      is_default: true
```

Sign in with the credentials from the `dashboards` Secret (name in
the stack outputs, along with a port-forward command). The datasource
`url` also accepts a reference to a KubernetesKubePrometheusStack
resource — the FK resolves to its Prometheus endpoint.

## Configuration

### State: the decision that matters most

Grafana keeps UI-authored dashboards, users and preferences in an
embedded SQLite database on local disk, and the chart default is
ephemeral — everything hand-made vanishes on pod restart. Three
postures:

- **Ephemeral** (default) — right when everything is provisioned as
  code (datasources, dashboards from ConfigMaps).
- **`storage`** — a ReadWriteOnce PVC for one stateful instance.
- **`database`** — external Postgres/MySQL; durable, and the
  REQUIREMENT for `replicas > 1` (SQLite cannot be shared — the spec
  enforces it). The password rides an existing Secret through
  environment expansion, never the rendered config.

### Provisioning as code

`datasources` render Grafana's provisioning files — present from
first boot, no clicking. The dashboard sidecar (on by default)
discovers any ConfigMap labeled `grafana_dashboard: "1"` cluster-wide
— the contract by which other components and teams ship dashboards to
this Grafana without touching its spec. `community_dashboards`
imports grafana.com dashboards by ID; pin `revision` for
reproducible installs.

### Plugins on Grafana 13

Several once-core datasource plugins (elasticsearch, cloudwatch)
moved out of the core image — list them in `plugins` to keep using
them. The module enables the chart's bundled-plugin shadowing so
installs succeed on the read-only image directory.

### Exposure and auth

The Service is ClusterIP. Compose a KubernetesIngress or Gateway API
route over the service handle and set `server.root_url` to the public
URL — OAuth redirects and rendered links embed it. `auth` covers
anonymous viewing (Viewer role by default) and login-form hiding for
pure-SSO fronts; LDAP/OAuth providers ride `helm_values`.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace Grafana runs in |
| `release_name` | Helm release name (= `metadata.name`) |
| `service` | The Grafana Service (port 80 → container 3000) |
| `endpoint` | In-cluster URL (`http://<name>.<ns>.svc.cluster.local`) |
| `admin_secret_name` | The credentials Secret — chart-owned `<name>`, or the referenced Secret echoed |
| `port_forward_command` | Workstation access when no exposure is composed |

## Related Components

- [KubernetesKubePrometheusStack](/docs/catalog/kubernetes/kube-prometheus-stack)
  — the cluster's metrics engine; its Prometheus endpoint is the
  default-kind reference for datasource URLs
- [KubernetesPostgres](/docs/catalog/kubernetes/postgres) — backs the
  external-database HA path via reference
- [KubernetesNamespace](/docs/catalog/kubernetes/namespace)
  — provides the target namespace via reference
- [KubernetesStorageClass](/docs/catalog/kubernetes/storage-class)
  — backs the storage volume via reference
- [KubernetesIngress](/docs/catalog/kubernetes/ingress) —
  composes exposure over the service handle

## Next Steps

Decide the state posture before anyone builds a dashboard by hand —
moving from ephemeral to persistent later cannot recover what a
restart already erased. Wire the first datasource at deploy time (the
stack reference is one line), let teams ship dashboards through
labeled ConfigMaps, and when Grafana becomes the pane of glass for
several backends, move state to an external database and scale
replicas behind the one Service.
