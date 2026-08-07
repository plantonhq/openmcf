---
title: "Altinity ClickHouse Operator"
description: "Altinity ClickHouse Operator deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesaltinityoperator"
---

# Altinity ClickHouse Operator

Install the [Altinity ClickHouse Operator](https://docs.altinity.com/clickhouse-operator/) — the Apache-2.0 operator for running ClickHouse (the columnar OLAP database) on Kubernetes — from the official `altinity-clickhouse-operator` Helm chart. The operator reconciles `ClickHouseInstallation` and `ClickHouseKeeperInstallation` custom resources (declared with **Kubernetes ClickHouse**) into running clusters with generated server configuration, rolling restarts, and per-host StatefulSets.

This component installs and configures the **engine**. ClickHouse clusters themselves are declared with Kubernetes ClickHouse resources — one per cluster — which this operator reconciles.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Helm release (altinity-clickhouse-operator)** — the operator Deployment with its metrics-exporter sidecar, its RBAC (cluster-wide by default, namespace-scoped Roles when `namespace_scoped_rbac` is set), the chart-managed credentials Secret, and the `<name>-metrics` Service
- **The four ClickHouse CRDs** (`clickhouse.altinity.com` and `clickhouse-keeper.altinity.com` API groups) — shipped in the chart's `crds/` directory: Helm installs them on first install and NEVER deletes them on uninstall, so destroying the operator never cascade-deletes ClickHouse clusters or their data; the chart's pre-install/pre-upgrade hook job server-side-applies them on every install and upgrade
- **Kubernetes Namespace** — created only when `create_namespace` is `true`; otherwise the operator installs into an existing namespace

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** — an active connection in the Connect module with credentials for the target cluster.

### Watch Scope

Unlike most operators, the chart default watches **only the operator's own namespace**. Decide up front where Kubernetes ClickHouse resources will live: cover every such namespace in `watch_namespaces` (entries are regular expressions), or use `[".*"]` for the whole cluster.

### Credentials

`operator_credentials` unset falls back to the chart's **publicly documented** username/password pair — fine for a throwaway environment, never for production.

## Deploy

### Console

Open the deployment store, find **Altinity ClickHouse Operator**, and click **Deploy**. The creation wizard walks you through namespace placement (with the 39-character naming budget), the chart pin, the three-posture watch scope, RBAC reach, the operator's own credentials, the CRD hook and metrics exporter, resources and scheduling, image sourcing, and the Helm-values escape hatch. Start from the **Standard** preset in the [Presets](#presets) tab.

### CLI

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesAltinityOperator
metadata:
  name: clickhouse-operator
  org: acme-corp
  env: prod
spec:
  namespace:
    value: clickhouse-operator
  create_namespace: true
  watch_namespaces:
    - ".*"
  operator_credentials:
    username: clickhouse_operator
    password:
      value: change-me-operator-password
```

```shell
planton apply -f clickhouse-operator.yaml
```

The install waits for the operator to become Available. From that point, Kubernetes ClickHouse resources in any namespace reconcile into running clusters — because this manifest widens the watch scope to `[".*"]`; the chart default watches only the operator's own namespace.

### InfraChart

Compose the operator with a managed namespace:

```yaml
spec:
  namespace:
    valueFrom:
      kind: KubernetesNamespace
      name: data-platform
      fieldPath: spec.name
  create_namespace: false
```

## Key Configuration

**Chart and image pinning** — `chart_version` defaults to **0.27.2**. Chart versions track operator releases one-to-one (chart 0.27.2 runs operator image 0.27.2), so the chart pin IS the operator pin; an overridden image tag should move with it. Versions must exist as SERVED charts in the repository index (docs.altinity.com/clickhouse-operator/index.yaml). Editing the version later IS the upgrade — the chart's built-in hook carries the CRD schema changes with it.

**Watch scope** — empty `watch_namespaces` watches ONLY the operator's own namespace (the chart default — the narrowest posture in the catalog). Entries are regular expressions: cover every namespace that will hold Kubernetes ClickHouse resources, or use `[".*"]` for cluster-wide. The fence is silent on the outside — a ClickHouse cluster declared beyond it is never reconciled, with no event pointing at the fence.

**RBAC reach** — `namespace_scoped_rbac` swaps cluster-wide RBAC for namespace-scoped Roles: the tenancy posture for shared clusters. Only sound when the operator watches its own namespace alone — a namespace-scoped operator cannot manage clusters it can watch but not touch.

**Operator credentials** — `operator_credentials` is the login the operator itself uses on every managed ClickHouse cluster (host management, schema propagation, metrics scraping), provisioned as a chart-managed Secret and injected as a network-restricted user into every managed cluster. The password accepts a literal or a reference to another resource's output. UNSET IS UNSAFE FOR PRODUCTION: the fallback pair is publicly documented upstream.

**CRD hook and metrics** — both default ON upstream. The hook is how chart upgrades carry CRD schema changes (disable only if CRD lifecycle is managed elsewhere); its default image is `bitnami/kubectl:latest` — pullable today but frozen since Bitnami retired its public catalog, so long-lived and air-gapped installs should pin their own kubectl build via `crd_hook.image`. The metrics-exporter sidecar serves Prometheus metrics for every managed cluster on port 8888 (the `metrics_endpoint` output); `service_monitor_enabled` adds a ServiceMonitor for Prometheus Operator scraping — enabling it without the Prometheus Operator CRDs fails the install.

**Naming budget** — keep the resource name at 39 characters or fewer: the modules pin the chart fullname to the resource name, and the longest generated child name adds 24 characters against the Kubernetes 63-character cap. `helm_values` merges LAST (Helm `-f` semantics) — except `fullnameOverride`, which the modules re-pin after the merge.

## Outputs and Dependencies

### What This Component Consumes

| Field | References | Purpose |
|-------|-----------|---------|
| `spec.namespace` | KubernetesNamespace (`spec.name`) | Where the operator installs |
| `spec.operator_credentials.password` | Any resource output or org secret | The operator's ClickHouse password |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `namespace` | The installation namespace | Composition, debugging |
| `release_name` | Helm release name (= metadata.name) | Helm-level operations |
| `deployment_name` | The operator Deployment (= the chart fullname, pinned to the resource name) | Monitoring, log collection |
| `credentials_secret_name` | Chart-managed Secret holding the operator's ClickHouse credentials | Auditing, credential rotation workflows |
| `metrics_endpoint` | In-cluster Prometheus metrics URL for every managed cluster; empty when metrics are disabled | Prometheus scrape configuration |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard** — one cluster-wide operator (`watch_namespaces: [".*"]`) with real credentials and modest sizing; the first thing to change is the password. Start from the **Standard** preset.

**Namespace scoped** — the tenancy posture: a single-namespace watch paired with `namespace_scoped_rbac`, one install per team. Start from the **Namespace Scoped** preset.

**Private mirror** — the air-gap posture: both images re-pointed at a private registry (including the frozen-upstream kubectl hook image) with a pull secret. Start from the **Private Mirror** preset.

## Works With

- **Kubernetes ClickHouse** — the ClickHouse clusters this operator reconciles; deploy the operator FIRST, and keep clusters inside the watched namespaces.
- **Kubernetes Namespace** — reference a managed namespace to compose governance (quotas, pod-security labels) with the installation.
- **Kubernetes Prometheus** — scrapes the fleet-wide `metrics_endpoint`; pair with `service_monitor_enabled` when the Prometheus Operator manages scrape targets.
