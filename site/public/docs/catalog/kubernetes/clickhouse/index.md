---
title: "ClickHouse"
description: "ClickHouse deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesclickhouse"
---

# Kubernetes ClickHouse

Declares a ClickHouse cluster — the columnar OLAP database built for
analytical queries over billions of rows — as a
`ClickHouseInstallation` reconciled by the Altinity ClickHouse
operator. One resource carries the shard×replica topology, per-host
persistent storage, users with Secret-delivered passwords, settings
profiles and quotas, and the operational verbs (stop, disruption
budget, replica anti-affinity). Coordination takes care of itself: a
managed ClickHouse Keeper is deployed automatically whenever the
topology needs one. Workloads connect at the exported native (9000)
and HTTP (8123) endpoints; exposure composes from ingress and gateway
kinds over the exported service handle.

> **Replication is enforced honestly**: `replicas > 1` without a
> coordination service is rejected at validation — ReplicatedMergeTree
> cannot sync without ClickHouse Keeper or ZooKeeper, so the spec
> refuses to let that cluster exist half-broken.

## What Gets Created

- **Namespace** (optional) — created and owned when `create_namespace`
  is set
- **ClickHouseInstallation** (`clickhouse.altinity.com/v1`, named
  `metadata.name`) — topology, storage, users, settings, placement
- **ClickHouseKeeperInstallation** (`<name>-keeper`) — the managed
  coordination ensemble, created exactly when the topology needs one
  (or as configured)
- **Auth Secret** (`<name>-clickhouse-auth`) — one key per declared
  user; passwords never appear in the custom resource

The operator reconciles these into one single-pod StatefulSet per
host (`chi-<name>-<cluster>-<shard>-<replica>`), the cluster-wide
client Service (`clickhouse-<name>`, ClusterIP), generated
configuration ConfigMaps, and a PodDisruptionBudget.

## Prerequisites

- The Altinity ClickHouse operator on the cluster
  (KubernetesAltinityOperator) — its chart default watches ONLY its
  own namespace, so set `watch_namespaces` to cover this cluster's
  namespace
- A StorageClass for the data volumes (most managed clusters provide
  a default; or reference a KubernetesStorageClass)
- Keep `metadata.name` within 48 characters (generated child names
  embed it against the Kubernetes 63-character Service cap)

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesClickHouse
metadata:
  name: analytics
spec:
  namespace:
    value: data
  create_namespace: true
  version: "25.3"
  replicas: 3
  disk_size: 100Gi
  users:
    - name: app
      password:
        value: change-me
      grants:
        - GRANT SELECT, INSERT ON analytics.*
```

Because `replicas` is 3 and `coordination` is unset, a three-node
managed Keeper deploys alongside the cluster automatically. Workloads
connect at `clickhouse-analytics.data.svc.cluster.local:9000` (native
protocol) or `:8123` (HTTP) with the password from the
`analytics-clickhouse-auth` Secret.

## Data Safety: Read This Before Production

Two switches deserve conscious ownership. **`retain_volumes_on_delete`**
defaults off — the operator's own reclaim policy DELETES every data
volume with the resource; turn it on in production so a re-created
cluster with the same name re-attaches its data (retained volumes are
never garbage-collected). And **password rotation is two steps**:
secret-sourced passwords reach ClickHouse through pod environment
variables, so after rotating the Secret, any spec change triggers the
re-reconcile that rolls the rotation out.

## Configuration

### Topology

`shards` × `replicas` is the cluster's shape: replicas are full
copies kept in sync through ReplicatedMergeTree (durability — turn on
`spread_replicas_across_nodes` in production), shards are disjoint
slices Distributed tables query in parallel (capacity). 1×1 is a dev
cluster; 1×3 is the production durability posture; scale shards only
when one shard cannot carry the dataset or the write rate.

### Coordination

Leave `coordination` unset and the managed Keeper appears exactly
when the topology needs it. Size it explicitly with `managed_keeper`
(quorum of 1, 3, or 5), point at existing ensembles with
`external_keeper` / `external_zookeeper`, or opt out with `none`
(single-replica only).

### Users and Access

Named users carry Secret-delivered passwords, profile and quota
references, network allowlists, and declarative SQL grants. The
built-in `default` user stays operator-managed: passwordless but
restricted to the cluster's own pods — every real client gets a named
user.

### Server Settings

ClickHouse's full configuration vocabulary passes through the CHI's
own path-keyed maps: `settings` (server config, e.g.
`max_concurrent_queries`), `profiles` and `quotas` (named bundles
users reference), and `files` (raw config-file drop-ins) — keys are
`/`-separated XML paths exactly as the upstream defines them.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the cluster runs in |
| `chi_name` | ClickHouseInstallation resource name (= `metadata.name`) |
| `cluster_name` | Logical cluster name — the `ON CLUSTER` target |
| `service_name` | Cluster-wide client Service (`clickhouse-<name>`) |
| `tcp_endpoint` | In-cluster native-protocol endpoint (port 9000) |
| `http_endpoint` | In-cluster HTTP endpoint (port 8123) |
| `auth_secret_name` | `<name>-clickhouse-auth` (one key per user); empty when no users |
| `keeper_name` | Managed Keeper resource (`<name>-keeper`); empty when external or none |
| `keeper_service_name` | Managed Keeper client Service (`keeper-<name>-keeper`); empty when external or none |
| `port_forward_command` | Workstation access when no exposure is composed |

## Related Components

- [KubernetesAltinityOperator](/docs/catalog/kubernetes/altinity-operator)
  — the engine; must be installed and watching this namespace
- [KubernetesNamespace](/docs/catalog/kubernetes/namespace) —
  provides the target namespace via reference
- [KubernetesStorageClass](/docs/catalog/kubernetes/storage-class)
  — referenced by the data and Keeper volume storage classes
- [KubernetesIngress](/docs/catalog/kubernetes/ingress) —
  composes external exposure over the exported service handle

## Next Steps

Turn on `retain_volumes_on_delete` and `spread_replicas_across_nodes`
before real data arrives. Replace placeholder passwords — user
passwords accept references to other resources' outputs, so generated
credentials flow in without being written down. Compose exposure from
KubernetesIngress or Gateway API kinds over `service_name` — this
component never embeds it.
