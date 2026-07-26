---
title: "Altinity Operator"
description: "Altinity Operator deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesaltinityoperator"
---

# Kubernetes Altinity Operator

Installs the Altinity ClickHouse operator — the Apache-2.0 operator
for running ClickHouse on Kubernetes — from the official
`altinity-clickhouse-operator` Helm chart. The operator reconciles
`ClickHouseInstallation` and `ClickHouseKeeperInstallation` custom
resources (declared with KubernetesClickHouse) into running clusters
with generated server configuration, rolling restarts, and per-host
StatefulSets. This component installs and configures the engine;
ClickHouse clusters are declared separately, one KubernetesClickHouse
resource per cluster.

## What Gets Created

- **Namespace** (optional) — created and owned when `create_namespace`
  is set
- **The four ClickHouse CRDs** (`clickhouse.altinity.com` and
  `clickhouse-keeper.altinity.com` API groups) — shipped in the
  chart's `crds/` directory: Helm installs them on first install and
  NEVER deletes them on uninstall, so destroying the operator never
  cascade-deletes ClickHouse clusters or their data; the chart's
  pre-install/pre-upgrade hook job server-side-applies them on every
  install and upgrade
- **Helm release** (official `altinity-clickhouse-operator` chart,
  pinned 0.27.2, named `metadata.name`) — the operator Deployment
  with the metrics-exporter sidecar, its RBAC, the chart-managed
  credentials Secret, and the `<name>-metrics` Service

## Prerequisites

- A Kubernetes namespace that already exists, or set
  `create_namespace`
- A real `operator_credentials` password outside throwaway
  environments — the chart's fallback credentials are publicly
  documented
- For `service_monitor_enabled`: the Prometheus Operator CRDs on the
  cluster

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesAltinityOperator
metadata:
  name: clickhouse-operator
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

The install waits for the operator to become Available. From that
point, KubernetesClickHouse resources in any namespace reconcile into
running clusters — because this manifest widens the watch scope to
`[".*"]`; the chart default watches only the operator's own namespace.

## Configuration

### Chart and image pinning

`chart_version` defaults to 0.27.2 — chart versions track operator
releases one-to-one, so the chart version and the operator image tag
move in lockstep. The `image` block overrides the operator image for
air-gapped and private-mirror installs; keep the tag matched to the
chart. The CRD hook's default image is `bitnami/kubectl:latest` —
pullable today but frozen since Bitnami retired its public catalog, so
long-lived and air-gapped installs should pin their own kubectl build
via `crd_hook.image`.

### Watch scope and RBAC

Empty `watch_namespaces` watches ONLY the operator's own namespace
(the chart default). Entries are regular expressions — cover every
namespace that will hold KubernetesClickHouse resources, or use
`[".*"]` for cluster-wide; a fenced operator silently ignores clusters
elsewhere. Pairing a single-namespace watch with
`namespace_scoped_rbac` avoids cluster-wide RBAC entirely on shared
clusters.

### Operator credentials

`operator_credentials` is the login the operator itself uses on every
managed ClickHouse cluster (host management, schema propagation,
metrics scraping), provisioned as a chart-managed Secret and
auto-injected as a network-restricted user into every managed cluster.
The password accepts a literal or a reference to another resource's
output.

### Metrics

The metrics-exporter sidecar (enabled by default) serves Prometheus
metrics for every managed cluster on port 8888, exported as
`metrics_endpoint`. `service_monitor_enabled` adds a ServiceMonitor
for Prometheus Operator scraping — enabling it without the Prometheus
Operator CRDs fails the install.

### Escape hatch

`helm_values` merges additional chart values LAST (Helm `-f`
semantics) — except `fullnameOverride`, which the modules re-pin to
the resource name after the merge (the chart fullname anchors every
generated child name and exported output). Keep the resource name at
39 characters or fewer — the longest generated child name adds 24
characters against the Kubernetes 63-character cap.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the operator is installed into |
| `release_name` | Helm release name (= `metadata.name`) |
| `deployment_name` | The operator Deployment (= the chart fullname, pinned to the resource name) |
| `credentials_secret_name` | Chart-managed Secret holding the operator's ClickHouse credentials |
| `metrics_endpoint` | In-cluster Prometheus metrics URL for every managed cluster; empty when metrics are disabled |

## Related Components

- [KubernetesClickHouse](/docs/catalog/kubernetes/clickhouse)
  — declares the ClickHouseInstallation (and managed-Keeper) resources
  this operator reconciles
- [KubernetesNamespace](/docs/catalog/kubernetes/namespace) —
  provides the installation namespace via reference

## Next Steps

Declare a ClickHouse cluster with KubernetesClickHouse — shards,
replicas, storage, Keeper coordination — and the operator reconciles
it. If the operator was installed with a narrow `watch_namespaces`
list, keep clusters inside the watched namespaces; they are silently
ignored anywhere else.
