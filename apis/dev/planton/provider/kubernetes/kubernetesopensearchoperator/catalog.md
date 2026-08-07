# OpenSearch Operator

Install the [OpenSearch Kubernetes Operator](https://github.com/opensearch-project/opensearch-k8s-operator) — the opensearch-project's operator for running OpenSearch (the Apache-2.0 search and analytics engine) on Kubernetes — from the official `opensearch-operator` Helm chart. The operator reconciles `OpenSearchCluster` custom resources (declared with **Kubernetes OpenSearch**) into running search clusters with managed TLS, security bootstrap, safe rolling upgrades, and OpenSearch Dashboards.

This component installs and configures the **engine**. Search clusters themselves are declared with Kubernetes OpenSearch resources — one per cluster — which this operator reconciles.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Helm release** — the operator controller-manager Deployment, its RBAC (cluster-wide by default, namespace-scoped with `use_role_bindings`), ServiceAccount, and metrics Service
- **The ten OpenSearch CRDs** — pinned by the module itself, NOT release-owned: the chart templates them as release resources with no keep-on-uninstall knob, so the module owns the CRD lifecycle to guarantee uninstalling the operator never cascade-deletes OpenSearchCluster resources and their data
- **kube-rbac-proxy sidecar** (default on) — shields the metrics endpoint behind Kubernetes RBAC. The chart's own default sidecar image was deleted upstream and can never be pulled; the module re-points it at the maintainer's quay.io repository at the chart's pinned tag, so the default posture actually installs

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** — an active connection in the Connect module with credentials for the target cluster.

### Naming Budget

Keep the resource name at **27 characters or fewer**. The chart derives long child names from its fullname (the longest adds 36 characters) and Kubernetes caps object names at 63 — the chart truncates the fullname but NOT the names built from it, so an over-long name fails at the API server mid-install.

## Deploy

### Console

Open the deployment store, find **OpenSearch Operator**, and click **Deploy**. The creation wizard walks you through namespace placement, the chart pin, watch scope and RBAC, operator runtime, manager endpoints, resources and scheduling, image sourcing, and the Helm-values escape hatch. Start from the **Standard** preset in the [Presets](#presets) tab.

### CLI

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesOpenSearchOperator
metadata:
  name: opensearch-operator
  org: acme-corp
  env: dev
spec:
  namespace:
    value: opensearch-operator
  create_namespace: true
```

```shell
planton apply -f opensearch-operator.yaml
```

### InfraChart

Compose the operator with a managed namespace:

```yaml
spec:
  namespace:
    valueFrom:
      kind: KubernetesNamespace
      name: search-platform
      fieldPath: spec.name
  create_namespace: false
```

## Key Configuration

**Chart pinning** — `chart_version` defaults to **2.8.0**, the newest STABLE chart/operator pairing. The chart's default manager image tag is its appVersion, and the served 2.8.3/2.8.4/3.0.x charts all default to a PRERELEASE operator image while bundling next-generation `opensearch.org` API-group CRDs the stable operator does not serve. A version bump is a license and API-group re-check — pin those lines only after upstream cuts a stable 3.x operator release.

**Watch scope** — by default the operator watches ALL namespaces (cluster-wide RBAC): one install serves every Kubernetes OpenSearch on the cluster. Set `watch_namespace` to fence it to one namespace; pair with `use_role_bindings` to swap ClusterRoleBindings for namespace-scoped RoleBindings — the shared-cluster posture needing no cluster-admin sign-off. The pairing is enforced: namespace-scoped RBAC cannot serve a cluster-wide operator. The fence is silent on the outside — a cluster declared beyond it is never reconciled, with no event pointing at the fence.

**CRD lifecycle is module-owned** — uninstalling the operator never cascade-deletes OpenSearchCluster resources. In Helm values, `installCRDs` is re-pinned by the module after the merge, and `nameOverride`/`fullnameOverride` are off limits (they break the exported `deployment_name` output).

**DNS domain** — `dns_base` is baked into the TLS certificates the operator generates; set it only on clusters with a non-default DNS domain, or every generated certificate's SANs mismatch the service names nodes advertise.

## Outputs and Dependencies

### What This Component Consumes

| Field | References | Purpose |
|-------|-----------|---------|
| `spec.namespace` | KubernetesNamespace (`spec.name`) | Where the operator installs |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `namespace` | The installation namespace | Composition, debugging |
| `release_name` | Helm release name (= metadata.name) | Helm-level operations |
| `deployment_name` | The controller-manager Deployment name (the chart's fixed component name within the release) | Monitoring, log collection |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard** — cluster-wide watch, pinned stable chart, chart defaults everywhere. Start from the **Standard** preset.

**Namespace-scoped** — a fenced watch with RoleBindings: the multi-tenant / no-cluster-admin posture. Start from the **Namespace Scoped** preset.

**Private mirror** — the air-gap path: manager image from your own registry with pull secrets and explicit resources. Start from the **Private Mirror** preset.

## Works With

- **Kubernetes OpenSearch** — the search clusters this operator reconciles; deploy the operator FIRST (it is the registered prerequisite).
- **Kubernetes Namespace** — reference a managed namespace to compose governance (quotas, pod-security labels) with the installation.
- **Kube Prometheus Stack** — scrape the shielded metrics endpoint with an RBAC-authorized ServiceMonitor.
