# Apache Solr Operator

Install the [Apache Solr Operator](https://solr.apache.org/operator/) — the Apache Solr project's own operator for running SolrCloud on Kubernetes — from the official `solr-operator` Helm chart ([solr.apache.org/charts](https://solr.apache.org/charts)). The operator reconciles `SolrCloud` custom resources (declared with **Kubernetes Solr**) into running Solr clusters with managed rolling updates, scaling with replica movement, and backup repositories.

This component installs and configures the **engine**. Solr clusters themselves are declared with Kubernetes Solr resources — one per cluster — which this operator reconciles.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Helm release (solr-operator)** — the operator Deployment (with leader election between replicas), its RBAC, and ServiceAccount
- **The operator CRDs** — SolrCloud, SolrBackup, SolrPrometheusExporter, and ZookeeperCluster (for the bundled dependency), applied by the module itself and keyed by CRD name — the chart does not ship them. Because the module owns the CRD lifecycle, uninstalling the operator never cascade-deletes SolrCloud resources and their data
- **The bundled zookeeper-operator** (default on) — the chart dependency that provisions PROVIDED ZooKeeper ensembles, so a Kubernetes Solr resource works out of the box. Disable it only when a zookeeper-operator already runs in the cluster (its fixed-name cluster-scoped RBAC conflicts on a second install) or when every Solr cluster connects to an EXTERNAL ensemble
- **Kubernetes Namespace** — created only when `create_namespace` is `true`; otherwise the operator installs into an existing namespace

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** — an active connection in the Connect module with credentials for the target cluster.

### ZooKeeper

SolrCloud requires ZooKeeper. With the defaults you need nothing: the bundled zookeeper-operator provisions provided ensembles. Check first whether the cluster already runs a zookeeper-operator — one per cluster, ever.

## Deploy

### Console

Open the deployment store, find **Apache Solr Operator**, and click **Deploy**. The creation wizard walks you through namespace placement, the chart pin, watch scope, the bundled zookeeper-operator, operator runtime, the mTLS client identity, resources and scheduling, image sourcing, and the Helm-values escape hatch. Start from the **Standard** preset in the [Presets](#presets) tab.

### CLI

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesSolrOperator
metadata:
  name: solr-operator
  org: acme-corp
  env: dev
spec:
  namespace:
    value: solr-operator
  create_namespace: true
```

```shell
planton apply -f solr-operator.yaml
```

The defaults are the point: chart 0.9.1, the bundled zookeeper-operator installed, and a cluster-wide watch — every Kubernetes Solr in any namespace gets reconciled, including its provided ZooKeeper ensemble.

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

**Chart pinning** — `chart_version` defaults to **0.9.1**. Chart and operator versions move together (the operator image tag is `v<chart_version>`), and the chart version has **no `v` prefix** while the operator/CRD artifacts carry one. The operator is **pre-1.0**: minor versions can change the CRD API, so a version bump means applying the matching CRDs — which the deploy does, because the CRD lifecycle is module-owned. Versions must exist as SERVED charts in the repository index.

**The bundled zookeeper-operator** — `zookeeper_operator.install` defaults true: the path that makes provided ZooKeeper ensembles work out of the box. Set `install: false` together with `use_existing: true` when the cluster already runs a zookeeper-operator (fixed-name cluster-scoped RBAC conflicts on a second install); plain `install: false` alone is the external-ensemble-only posture — a provided ensemble would then hang silently waiting on ZooKeeper.

**Watch scope** — by default the operator watches ALL namespaces: one install serves every Kubernetes Solr on the cluster. Set `watch_namespaces` to fence it to an explicit list — the multi-tenant posture where several fenced installs share a cluster. The fence is silent on the outside (a Solr cluster declared beyond it is never reconciled, with no event pointing at the fence), and the watched namespaces need only exist by the time SolrCloud resources appear in them.

**mTLS client identity** — configure `mtls` only when Kubernetes Solr resources enforce `clientAuth` on their TLS listeners; the operator then presents `mtls.client_cert_secret` (required within the block) when calling Solr. `insecure_skip_verify` defaults true because the operator dials pods by pod IP — addresses that rarely appear in certificate SANs; the CA trust still authenticates both ends. All the Secret fields hold NAMES, never certificate material.

**Operator runtime** — `replicas` (default 1) adds leader-elected warm standbys, never throughput; `leader_election_enabled` (default true) should be disabled only for a single-replica dev install; `metrics_enabled` (default true) exposes the operator's Prometheus endpoint. `resources` is empty by default — the chart sets none; the operator is lightweight.

## Outputs and Dependencies

### What This Component Consumes

| Field | References | Purpose |
|-------|-----------|---------|
| `spec.namespace` | KubernetesNamespace (`spec.name`) | Where the operator installs |
| `spec.mtls.client_cert_secret` | KubernetesSecret (`metadata.name`) | The client certificate the operator presents to Solr |
| `spec.mtls.ca_cert_secret` | KubernetesSecret (`metadata.name`) | The CA to trust when calling Solr |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `namespace` | The installation namespace | Composition, debugging |
| `release_name` | Helm release name (= metadata.name) | Helm-level operations |
| `deployment_name` | The Solr operator Deployment name | Monitoring, log collection |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard** — one cluster-wide operator plus the bundled zookeeper-operator on chart 0.9.1; the first thing to change is nothing. Start from the **Standard** preset.

**Existing zookeeper-operator** — the second-install posture: `install: false` + `use_existing: true` when another release already runs the dependency. Start from the **Existing Zookeeper Operator** preset.

**Namespace fenced** — an explicit `watch_namespaces` list, a warm-standby replica, and declared resources: the multi-tenant / governed-platform posture. Start from the **Namespace Fenced** preset.

## Works With

- **Kubernetes Solr** — the SolrCloud clusters this operator reconciles; deploy the operator FIRST (it is the registered prerequisite).
- **Kubernetes Namespace** — reference a managed namespace to compose governance (quotas, pod-security labels) with the installation.
- **Kubernetes Secret** — the mTLS client and CA certificate Secrets the operator presents when Solr clusters enforce clientAuth.
