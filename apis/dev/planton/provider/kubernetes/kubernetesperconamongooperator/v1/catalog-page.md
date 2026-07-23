# Kubernetes Percona Mongo Operator

Installs the Percona Operator for MongoDB from the official
`psmdb-operator` Helm chart, with a typed spec over the chart's
meaningful configuration surface. The operator reconciles
`PerconaServerMongoDB` custom resources into highly available MongoDB
deployments: replica sets with automated failover, sharded clusters
(mongos + config servers), scheduled logical/physical/incremental
backups with point-in-time recovery via Percona Backup for MongoDB,
TLS, and user management. This component installs the ENGINE; the
databases themselves are declared with KubernetesMongodb resources —
one per MongoDB cluster.

## What Gets Created

- **Namespace** (optional) — the installation namespace, created and
  owned when `create_namespace` is set
- **Helm Release** (named `metadata.name`) — the operator Deployment,
  its RBAC, and the PerconaServerMongoDB CRDs
  (PerconaServerMongoDB / PerconaServerMongoDBBackup /
  PerconaServerMongoDBRestore). The CRDs ship in the chart's
  Helm-native `crds/` directory: installed on first install, never
  upgraded or deleted by Helm — uninstalling the release NEVER
  cascade-deletes the database clusters. The same posture means a
  `chart_version` upgrade runs new operator code against the EXISTING
  CRDs — apply the new release's CRDs yourself when an upgrade's
  release notes call for it.

## Watch Scope

By default the operator watches ITS OWN namespace only (the upstream
default — databases live beside their operator). Two knobs widen it,
mutually exclusive:

- **`watch.cluster_wide: true`** — one operator manages databases in
  every namespace (cluster-wide RBAC)
- **`watch.namespaces: [...]`** — an explicitly fenced set of
  namespaces

A KubernetesMongodb database declared in a namespace the operator does
not watch is silently never reconciled — install the operator in the
databases' namespace, or widen the watch.

## Prerequisites

- The cluster must not already run the Percona MongoDB operator with
  overlapping watch scope — the CRDs are cluster-scoped singletons
- Nothing else: databases, backup stores, and TLS issuers are concerns
  of the KubernetesMongodb resources composed on top

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesPerconaMongoOperator
metadata:
  name: psmdb-operator
spec:
  namespace:
    value: percona-mongo
  create_namespace: true
  disable_telemetry: true
```

The operator becomes Available and starts reconciling;
KubernetesMongodb resources declared in the `percona-mongo` namespace
turn into running replica sets.

## Configuration Surface

- **`chart_version`** — the chart pin (default `1.22.0`; chart and
  operator versions move together). Pin deliberately; upgrades re-run
  the release with the new chart
- **`replicas`** — operator pods (default 1); extras are
  leader-elected warm standbys for the OPERATOR itself, not
  reconciliation throughput
- **`watch`** — the scope block above
- **`max_concurrent_reconciles`** — PerconaServerMongoDB resources
  reconciled concurrently (default 1); raise on control planes
  managing many databases
- **`log`** — `structured` (JSON) and `level` (DEBUG / INFO / ERROR)
- **`disable_telemetry`** — stop the anonymous version/feature pings
  to check.percona.com
- **`resources`, `node_selector`, `tolerations`,
  `image_pull_secrets`, `image`** — operator pod placement, sizing
  (the chart ships no requests/limits by default), and image override
  (registry mirrors, air-gapped clusters)
- **`helm_values`** — escape hatch: extra chart values merged LAST
  (Helm `-f` semantics) for the surface beyond the typed fields —
  never the primary interface, never for secrets

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the operator runs in |
| `release_name` | Helm release name of the operator (equals `metadata.name`) |

## Next Steps

Declare databases as KubernetesMongodb resources — one per MongoDB
cluster, in a namespace the operator watches; the operator reconciles
them into replica-set members, Services, and credential Secrets. Pin
`chart_version` deliberately and upgrade the operator on the
platform's schedule — removing this component never deletes the CRDs
or the databases behind them.
