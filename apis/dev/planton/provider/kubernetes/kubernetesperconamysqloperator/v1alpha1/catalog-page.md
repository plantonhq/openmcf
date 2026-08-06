# Kubernetes Percona MySQL Operator

Installs the Percona Operator for MySQL — based on Percona XtraDB
Cluster — from the official `pxc-operator` Helm chart, with a typed
spec over the chart's meaningful configuration surface. The operator
reconciles `PerconaXtraDBCluster` custom resources into highly
available MySQL clusters: Galera synchronous multi-primary replication
with automated failover, HAProxy or ProxySQL query routing, scheduled
XtraBackup backups with point-in-time recovery, and TLS. This component
installs the ENGINE; the databases themselves are declared with
KubernetesMysql resources — one per MySQL cluster.

## What Gets Created

- **Namespace** (optional) — the installation namespace, created and
  owned when `create_namespace` is set
- **Helm Release** (named `metadata.name`) — the operator Deployment,
  its RBAC, and the PerconaXtraDBCluster CRDs
  (PerconaXtraDBCluster / PerconaXtraDBClusterBackup /
  PerconaXtraDBClusterRestore). The CRDs ship in the chart's
  Helm-native `crds/` directory: installed on first install, never
  upgraded or deleted by Helm — uninstalling the release NEVER
  cascade-deletes the database clusters. The same posture means a
  `chart_version` upgrade runs new operator code against the EXISTING
  CRDs — apply the new release's CRDs yourself when an upgrade's
  release notes call for it.
- **Validation webhook** (widened-watch arms only) — the operator's
  cluster-scoped `percona-xtradbcluster-webhook`, rendered by the
  module rather than left to the operator's own startup registration.
  Module ownership is what makes uninstall clean: the upstream-created
  object outlives its operator (Kubernetes never garbage-collects a
  cluster-scoped dependent of a namespaced owner) and would brick every
  future PerconaXtraDBCluster admission once its Service is gone. The
  operator still manages the certificate: it refreshes the CA bundle
  into the module-declared object at startup. Because the webhook's
  name is fixed upstream, one cluster carries at most ONE widened-watch
  operator.

## Watch Scope

By default the operator watches ITS OWN namespace only (the upstream
default — databases live beside their operator). Two knobs widen it,
mutually exclusive:

- **`watch.cluster_wide: true`** — one operator manages databases in
  every namespace (cluster-wide RBAC)
- **`watch.namespaces: [...]`** — an explicitly fenced set of
  namespaces

A KubernetesMysql database declared in a namespace the operator does
not watch is silently never reconciled — install the operator in the
databases' namespace, or widen the watch.

## Prerequisites

- The cluster must not already run the Percona MySQL operator with
  overlapping watch scope — the CRDs are cluster-scoped singletons
- Nothing else: databases, backup stores, and TLS issuers are concerns
  of the KubernetesMysql resources composed on top

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesPerconaMysqlOperator
metadata:
  name: pxc-operator
spec:
  namespace:
    value: percona-mysql
  create_namespace: true
  disable_telemetry: true
```

The operator becomes Available and starts reconciling; KubernetesMysql
resources declared in the `percona-mysql` namespace turn into running
Galera clusters.

## Configuration Surface

- **`chart_version`** — the chart pin (default `1.20.0`; chart and
  operator versions move together). Pin deliberately; upgrades re-run
  the release with the new chart
- **`replicas`** — operator pods (default 1); extras are
  leader-elected warm standbys for the OPERATOR itself, not
  reconciliation throughput
- **`watch`** — the scope block above
- **`max_concurrent_reconciles`** — PerconaXtraDBCluster resources
  reconciled concurrently (default 1); raise on control planes
  managing many databases
- **`s3_workers_limit`** — concurrent S3 workers for backup
  uploads/deletes across the managed databases (default 10)
- **`log`** — `structured` (JSON) and `level` (DEBUG / INFO / ERROR)
- **`disable_telemetry`** — stop the anonymous version/feature pings
  to check.percona.com
- **`leader_election`** — election tuning between operator replicas
  (meaningful only with `replicas` above 1)
- **`xtrabackup_sidecar`** — run XtraBackup as a sidecar of the
  database pods instead of separate backup jobs (chart feature gate)
- **`resources`, `node_selector`, `tolerations`,
  `image_pull_secrets`, `image`** — operator pod placement, sizing,
  and image override (registry mirrors, air-gapped clusters)
- **`helm_values`** — escape hatch: extra chart values merged LAST
  (Helm `-f` semantics) for the surface beyond the typed fields —
  never the primary interface, never for secrets

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the operator runs in |
| `release_name` | Helm release name of the operator (equals `metadata.name`) |

## Next Steps

Declare databases as KubernetesMysql resources — one per MySQL
cluster, in a namespace the operator watches; the operator reconciles
them into Galera nodes, proxy Services, and credential Secrets. Pin
`chart_version` deliberately and upgrade the operator on the
platform's schedule — removing this component never deletes the CRDs
or the databases behind them.
