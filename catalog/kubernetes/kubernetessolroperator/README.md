# Kubernetes Solr Operator

## When NOT to Use This

**This component installs the ENGINE, not a search cluster.** The
Apache Solr Operator reconciles `SolrCloud` custom resources into
running Solr clusters; those clusters are declared with KubernetesSolr
— one resource per cluster. Install the operator once per Kubernetes
cluster, then declare SolrCloud clusters against it.

Also not the right component when:

- **You want a Solr cluster** — that is KubernetesSolr; this component
  is the controller that reconciles it.
- **You want a managed cloud search service** — use the host cloud
  provider's managed search kinds; this component is for running
  SolrCloud ON the Kubernetes cluster itself.
- **A zookeeper-operator already runs on the cluster** — do not
  install a second one: its fixed-name cluster-scoped RBAC conflicts
  on a second install. Set `zookeeper_operator.install: false` with
  `use_existing: true` instead.

## Overview

**KubernetesSolrOperator** installs the Apache Solr Operator — the
Apache Solr project's own operator for running SolrCloud on Kubernetes
— from the official `solr-operator` Helm chart
(https://solr.apache.org/charts). The operator reconciles `SolrCloud`
custom resources into running Solr clusters with managed rolling
updates, scaling with replica movement, and backup repositories — plus
`SolrBackup` and `SolrPrometheusExporter` resources.

**Key design points:**

- **ZooKeeper comes bundled.** SolrCloud requires ZooKeeper, and the
  chart bundles the zookeeper-operator as a dependency (installed by
  default) — so a KubernetesSolr resource can simply declare a
  provided ensemble and the operator provisions it. Disable the
  bundled install only when a zookeeper-operator already runs in the
  cluster (pair `install: false` with `use_existing: true`) or every
  Solr cluster will connect to an EXTERNAL ensemble.
- **The module owns the CRDs.** Unlike most operator charts, the
  solr-operator chart ships NO CRDs — they are separate release
  artifacts. The modules apply all four themselves (SolrCloud,
  SolrBackup, SolrPrometheusExporter, and the ZookeeperCluster CRD of
  the bundled dependency), keyed by CRD name, with a keep-on-uninstall
  posture: Terraform's `kubectl_manifest` `apply_only = true` (the
  provider's Delete is a no-op) and Pulumi's `retainOnDelete`
  transformation. Destroying the operator never cascade-deletes
  SolrCloud resources. The bundled subchart's own CRD switch is pinned
  off (`zookeeper-operator.crd.create: false`) so the ZookeeperCluster
  CRD never falls under Helm's delete-on-uninstall lifecycle.
- **Chart and CRD versions pair.** Chart and operator versions move
  together (chart 0.9.1 ships operator v0.9.1 — the chart version has
  no `v` prefix, the operator/CRD artifacts carry one), and the
  operator is pre-1.0: minor versions can change the CRD API. Restage
  the matching CRD files when upgrading — server-side apply lets a
  restaged CRD apply as an in-place update.
- **mTLS client identity is modeled.** When KubernetesSolr clusters
  enforce client certificates on their TLS listeners
  (`client_auth: Need`), the operator needs its own certificate to
  reach them — `mtls.client_cert_secret` (required whenever the block
  is declared), an optional CA secret, hostname-verification and
  rotation-watch knobs.
- **Watch scope defaults to cluster-wide.** `watch_namespaces` fences
  the operator to an explicit set; the watched namespaces need only
  exist by the time SolrCloud resources appear in them.
- **`helm_values` is the escape hatch** — additional chart values
  merged LAST over everything the typed fields render (Helm `-f`
  semantics, identical on both engines); a safety valve, never the
  primary interface.

## Essential Configuration Fields

### Required

- **`spec.namespace`**: namespace to install the operator into —
  literal or a KubernetesNamespace reference (`create_namespace` to
  own it)

### Common

- **`spec.chart_version`**: chart pin (default `0.9.1`); a bump means
  restaging the matching CRDs
- **`spec.zookeeper_operator`**: the bundled dependency — `install`
  (default true) and `use_existing` (tell the Solr operator a
  zookeeper-operator is present even though this release does not
  install one)
- **`spec.watch_namespaces`**: watch exactly these namespaces; empty =
  all (the chart default)
- **`spec.replicas`**: operator replicas (default 1; extras are
  leader-elected warm standbys)
- **`spec.mtls`**: the operator's client identity for mutual-TLS Solr
  clusters — `client_cert_secret` (tls.crt/tls.key), `ca_cert_secret`
  + `ca_cert_secret_key` (default `ca-cert.pem`),
  `insecure_skip_verify` (chart default true — pod-IP calls rarely
  match certificate SANs), `watch_for_updates` (restart on rotation,
  default true)
- **`spec.leader_election_enabled` / `spec.metrics_enabled`**: chart
  defaults true; disable leader election only for a single-replica dev
  install
- **`spec.resources`**: empty = the chart defaults (none set — the
  operator is lightweight)
- **`spec.node_selector` / `spec.tolerations`**: operator pod
  scheduling
- **`spec.image` / `spec.image_pull_secret`**: air-gap path (empty =
  `apache/solr-operator` at the chart's matching tag; the chart
  accepts exactly ONE pull secret)
- **`spec.helm_values`**: the escape hatch

## Stack Outputs

| Output | Purpose |
|---|---|
| `namespace` | Namespace the operator is installed into |
| `release_name` | Helm release name of the operator install (= `metadata.name`) |
| `deployment_name` | The Solr operator Deployment (the chart's fullname helper) |

## Composing in Infra Charts

- **`spec.namespace`** is a foreign key (default kind
  KubernetesNamespace, field path `spec.name`);
  **`mtls.client_cert_secret`** / **`ca_cert_secret`** reference
  existing Secrets (KubernetesSecret by default) — issue the client
  certificate with a KubernetesCertificate and reference its Secret.
- **KubernetesSolr resources depend on this component**: the operator
  (and, for provided ensembles, the bundled zookeeper-operator) must
  be running before SolrCloud resources reconcile. With
  `watch_namespaces` set, clusters outside the list are silently
  ignored.
- **The install is deliberately blocking**: the Helm release waits for
  the operator to become Available (atomic, 600s timeout), so an
  unpullable image fails THIS apply instead of surfacing later as
  SolrCloud resources that never reconcile.

## Examples

### Standard install (bundled zookeeper-operator)

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesSolrOperator
metadata:
  name: solr-operator
spec:
  namespace:
    value: solr-operator
  create_namespace: true
```

### Cluster already runs a zookeeper-operator

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesSolrOperator
metadata:
  name: solr-operator
spec:
  namespace:
    value: solr-operator
  create_namespace: true
  zookeeper_operator:
    install: false
    use_existing: true
```

### Fenced to team namespaces, with an mTLS client identity

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesSolrOperator
metadata:
  name: solr-operator
spec:
  namespace:
    value: solr-operator
  create_namespace: true
  watch_namespaces:
    - search
    - analytics
  mtls:
    client_cert_secret:
      value: solr-operator-client-tls
    ca_cert_secret:
      value: solr-ca
```

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
