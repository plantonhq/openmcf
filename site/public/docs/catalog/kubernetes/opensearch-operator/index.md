---
title: "OpenSearch Operator"
description: "OpenSearch Operator deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesopensearchoperator"
---

# Kubernetes OpenSearch Operator

Installs the OpenSearch Kubernetes Operator — the opensearch-project's
operator for running OpenSearch on Kubernetes — from the official
`opensearch-operator` Helm chart. The operator reconciles
`OpenSearchCluster` custom resources (declared with
KubernetesOpenSearch) into running search clusters with managed TLS,
security bootstrap, safe rolling upgrades, and OpenSearch Dashboards.
This component installs and configures the engine; search clusters are
declared separately, one KubernetesOpenSearch resource per cluster.

## What Gets Created

- **Namespace** (optional) — created and owned when `create_namespace`
  is set
- **The ten OpenSearch CRDs** (`opensearch.opster.io` API group) —
  applied as MODULE-OWNED resources with a keep-on-uninstall posture:
  destroying the operator never deletes the CRDs, so
  OpenSearchCluster resources and their data are never
  cascade-deleted (`installCRDs` is pinned false; the chart never
  touches them)
- **Helm release** (official `opensearch-operator` chart, pinned
  2.8.0, named `metadata.name`) — the operator controller-manager
  Deployment, its RBAC, and the kube-rbac-proxy metrics shield

## Prerequisites

- A Kubernetes namespace that already exists, or set
  `create_namespace`
- For `use_role_bindings`: a `watch_namespace` — namespace-scoped RBAC
  cannot serve a cluster-wide operator

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesOpenSearchOperator
metadata:
  name: opensearch-operator
spec:
  namespace:
    value: opensearch-operator
  create_namespace: true
```

The install waits for the operator to become Available. From that
point, KubernetesOpenSearch resources in any namespace reconcile into
running clusters (the default watch scope is cluster-wide; set
`watch_namespace` to fence it).

## Configuration

### Chart and image pinning

`chart_version` defaults to 2.8.0 — the newest served chart whose
default manager image is a stable operator release. The newer served
charts (2.8.3+, 3.0.x) default to a prerelease manager image, and the
3.x line migrates the CRDs from `opensearch.opster.io` to
`opensearch.org`; pin those lines only after upstream cuts a stable
3.x operator release. The `image` block overrides the manager image
for air-gapped and private-mirror installs.

### Watch scope and RBAC

Empty `watch_namespace` watches all namespaces on cluster-wide RBAC.
Setting it fences the operator to one namespace; pairing it with
`use_role_bindings` avoids cluster-wide RBAC entirely on shared
clusters.

### Escape hatch

`helm_values` merges additional chart values LAST (Helm `-f`
semantics) — except `installCRDs`, which the modules re-pin false
after the merge (the module owns the CRD lifecycle), and
`nameOverride`/`fullnameOverride`, which would break the exported
`deployment_name`.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the operator is installed into |
| `release_name` | Helm release name (= `metadata.name`) |
| `deployment_name` | The operator controller-manager Deployment |

## Related Components

- [KubernetesOpenSearch](/docs/catalog/kubernetes/opensearch)
  — declares the OpenSearchCluster resources this operator reconciles
- [KubernetesNamespace](/docs/catalog/kubernetes/namespace) —
  provides the installation namespace via reference

## Next Steps

Declare a search cluster with KubernetesOpenSearch — node pools,
security posture, Dashboards, snapshot repositories — and the operator
reconciles it. If the operator was installed with `watch_namespace`,
keep clusters inside that namespace; they are silently ignored
anywhere else.
