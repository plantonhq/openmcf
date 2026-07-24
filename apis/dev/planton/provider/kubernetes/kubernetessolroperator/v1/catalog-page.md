# Kubernetes Solr Operator

Installs the Apache Solr Operator — the Apache Solr project's own
operator for running SolrCloud on Kubernetes — from the official
`solr-operator` Helm chart. The operator reconciles `SolrCloud` custom
resources (declared with KubernetesSolr) into running Solr clusters
with shard-aware rolling updates, replica-moving scale operations, and
backup repositories. This component installs and configures the
engine, including the bundled zookeeper-operator dependency that
provisions ZooKeeper ensembles for Solr clusters; the clusters
themselves are declared separately, one KubernetesSolr resource per
cluster.

## What Gets Created

- **Namespace** (optional) — created and owned when `create_namespace`
  is set
- **Four CRDs** — SolrCloud, SolrBackup, SolrPrometheusExporter
  (`solr.apache.org`) and ZookeeperCluster (`zookeeper.pravega.io`),
  applied as MODULE-OWNED resources with a keep-on-uninstall posture:
  destroying the operator never deletes the CRDs, so SolrCloud
  resources are never cascade-deleted (the chart ships no CRDs; the
  bundled subchart's own CRD switch is pinned off)
- **Helm release** (official `solr-operator` chart, pinned 0.9.1,
  named `metadata.name`) — the operator Deployment and, by default,
  the bundled zookeeper-operator

## Prerequisites

- A Kubernetes namespace that already exists, or set
  `create_namespace`
- If a zookeeper-operator ALREADY runs on the cluster: set
  `zookeeper_operator.install: false` with `use_existing: true` — its
  fixed-name cluster-scoped RBAC conflicts on a second install
- For mutual-TLS Solr clusters: an existing client-certificate Secret
  for the `mtls` block

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesSolrOperator
metadata:
  name: solr-operator
spec:
  namespace:
    value: solr-operator
  create_namespace: true
```

The install waits for the operator to become Available. From that
point, KubernetesSolr resources in any namespace reconcile into
running SolrCloud clusters — including operator-provisioned ZooKeeper
ensembles, courtesy of the bundled zookeeper-operator.

## Configuration

### ZooKeeper posture

The bundled zookeeper-operator installs by default — the path that
makes provided ensembles work out of the box. Disable it only when one
already runs in the cluster (pair with `use_existing: true`) or every
Solr cluster connects to an external ensemble.

### Version pairing

Chart and operator versions move together (chart 0.9.1 ships operator
v0.9.1), and the operator is pre-1.0: minor versions can change the
CRD API, so a chart bump means restaging the matching CRDs — the
modules apply them with server-side apply, so a restage lands as an
in-place update.

### Watch scope

Empty `watch_namespaces` watches all namespaces; a list fences the
operator to exactly those.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the operator is installed into |
| `release_name` | Helm release name (= `metadata.name`) |
| `deployment_name` | The Solr operator Deployment |

## Related Components

- [KubernetesSolr](/docs/catalog/kubernetes/kubernetessolr) — declares
  the SolrCloud clusters this operator reconciles
- [KubernetesNamespace](/docs/catalog/kubernetes/kubernetesnamespace) —
  provides the installation namespace via reference
- [KubernetesCertificate](/docs/catalog/kubernetes/kubernetescertificate)
  — issues the client certificate the `mtls` block references

## Next Steps

Declare a search cluster with KubernetesSolr — replicas, ZooKeeper
wiring, persistent storage, basic-auth security, backup repositories —
and the operator reconciles it. If clusters will enforce client
certificates (`client_auth: Need`), give the operator its own mTLS
identity here first; its probes and reconciliation calls fail without
one.
