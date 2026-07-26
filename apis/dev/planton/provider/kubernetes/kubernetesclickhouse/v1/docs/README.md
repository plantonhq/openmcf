# KubernetesClickHouse: Research and Design

## Introduction

KubernetesClickHouse declares one ClickHouse cluster as a
`clickhouse.altinity.com/v1` `ClickHouseInstallation` (CHI) custom
resource reconciled by the Altinity ClickHouse operator
(KubernetesAltinityOperator, the registry prerequisite — pinned at
the chart 0.27.2 / operator 0.27.2 pairing). One resource carries the
whole cluster story: shard×replica topology, coordination, storage,
users, settings layers, and the operational verbs — plus, when the
topology needs coordination, one companion
`ClickHouseKeeperInstallation` the module manages alongside it.

## The Deployment Landscape

ClickHouse is the columnar OLAP database: data stored by column,
vectorized execution, compression ratios and scan speeds that make
analytical queries over billions of rows interactive. It is the
storage engine behind observability platforms, event analytics, and
data-warehouse workloads — and it is deliberately NOT an OLTP
database.

Running it on Kubernetes without an operator means hand-rolling a
multi-dimensional topology: N shards × M replicas where every host
needs its own stable identity, its own generated XML configuration,
membership in `remote_servers`, and a coordination service wired in
for replication and distributed DDL. The Altinity operator encodes
all of that — it has been reconciling production ClickHouse for
years, is Apache-2.0 end to end, and ships a second CRD that manages
ClickHouse Keeper the same way. So this kind is deliberately thin: it
renders ONE ClickHouseInstallation (plus one managed
ClickHouseKeeperInstallation when the topology calls for it) and
exports the operator's deterministic names.

## Upstream Architecture

The operator reconciles the declared resource into:

- **One single-pod StatefulSet per host**
  (`chi-<name>-<cluster>-<shard>-<replica>`) — the operator's own
  layout decision: every shard×replica host is individually
  addressable, individually restartable, and carries its own PVCs.
  Each host also gets a headless Service of the same name.
- **The cluster-wide client Service** (`clickhouse-<name>`) covering
  every ready host, and a per-cluster Service
  (`cluster-<name>-<cluster_name>`). All of them ClusterIP — the
  operator's own default, verified in the operator source; there is
  no LoadBalancer surprise.
- **Generated ConfigMaps** carrying the rendered ClickHouse XML —
  common config, users config, and per-host macros (shard/replica
  identity for ReplicatedMergeTree paths).
- **A PodDisruptionBudget per cluster** (`pdb_max_unavailable`).

The `cluster_name` is a segment of every generated child name, which
is why the upstream CRD caps it at 15 characters; with the default
("main"), keeping `metadata.name` within 48 characters keeps every
generated Service name under the Kubernetes 63-character cap.

### The coordination design

Replication is ReplicatedMergeTree, and ReplicatedMergeTree cannot
sync without a coordination service — so the spec enforces it: a CEL
rule rejects `replicas > 1` with coordination type `none` at
validation, never at runtime.

The design puts the common case at zero configuration:

- **Unset (the recommended default)**: the module deploys a managed
  ClickHouse Keeper — a `ClickHouseKeeperInstallation` named
  `<name>-keeper`, reconciled by the same operator — automatically
  whenever the topology needs coordination (`replicas > 1` or
  `shards > 1`), and nothing otherwise. Keeper is ClickHouse's own
  Raft-based ZooKeeper replacement: C++ (no JVM), protocol-compatible,
  the upstream-recommended default.
- **`managed_keeper`** sizes it explicitly: quorum of 1 (dev), 3
  (production — survives one loss), or 5, plus resources and a
  metadata-only volume (default 10Gi).
- **`external_keeper` / `external_zookeeper`** point at existing
  ensembles (`nodes`, optional `root` znode path for shared
  ensembles, optional digest-auth `identity`).
- **`none`** opts out — valid only single-replica; a multi-shard
  cluster without coordination loses `ON CLUSTER` DDL (run DDL per
  shard) but Distributed queries keep working.

The CHI references the managed Keeper through the CRD's NATIVE keeper
reference rather than a hardcoded host list — the operator resolves
the endpoints itself. At the pinned release the operator resolves
that reference to the Keeper's ready-only CLIENT service, which is
split from the Raft peer service (the peer service publishes
not-ready addresses so quorum can form; clients never see a host that
cannot serve).

### The users truth

Passwords never land in the CHI. The module writes every declared
user's password into one Kubernetes Secret
(`<name>-clickhouse-auth`, key = user name) and the CHI's users
section carries only `valueFrom.secretKeyRef` references the operator
injects. Two consequences, both upstream-documented:

1. **Secret-sourced passwords reach ClickHouse through pod
   environment variables** — so rotating the Secret alone does not
   re-render the configuration. A spec change triggers the
   re-reconcile that rolls the rotation out.
2. **The built-in `default` user stays operator-managed**:
   passwordless but network-restricted to the cluster's own pods (the
   operator generates the IP allowlist from the pod addresses). It
   exists for the operator's and the cluster's internal traffic —
   every real client gets a named user.

Declared users carry the full upstream vocabulary: `profile` and
`quota` references, `networks` allowlists, SQL `grants` executed at
config render, `access_management` for administrative users, and
per-user path-keyed `settings`.

Two grant behaviors to know (verified against a live cluster): a user
declared **without** `grants` receives ClickHouse's unrestricted
default access — declare grants to constrain every real user. And a
grants-constrained user running distributed DDL (`... ON CLUSTER`,
how schema reaches replicated and sharded topologies) additionally
needs `GRANT CLUSTER ON *.*` — database-scoped `CREATE` alone is
rejected with `ACCESS_DENIED`.

### Server configuration layers

Typed fields cover what every deployment decides: topology, storage,
users, placement, and the operational verbs. Everything else —
hundreds of server settings, settings profiles, quota intervals, raw
config-file drop-ins — rides the CHI's own path-keyed maps
(`settings`, `profiles`, `quotas`, `files`, per-user `settings`)
where keys are `/`-separated XML paths exactly as the upstream CRD
defines them. This is not an escape hatch bolted onto a typed API:
the upstream IS maps at this layer, and the component passes them
through unaltered. `files` names may carry the upstream's placement
prefixes (`{common}`, `{users}`, `{hosts}`) to choose between
config.d, users.d and conf.d.

### Storage and the reclaim policy

Every host gets a data PVC (`disk_size`, optional `storage_class`
reference) mounted at /var/lib/clickhouse, and optionally a separate
log volume (`log_disk_size`) at /var/log/clickhouse-server.
`retain_volumes_on_delete` maps to the operator's PVC reclaim policy:
the upstream default is Delete — PVCs vanish with their StatefulSet —
and Retain keeps them, so a re-created cluster with the same name
re-attaches its data. Retained PVCs are never garbage-collected;
deleting them becomes a deliberate manual act.

### Operational verbs

- **`stopped`** maps to the CHI `stop` verb: every host StatefulSet
  scales to zero (pods and Services go away, every PVC stays), and
  flipping it back brings the same data up — the declarative pause
  switch for expensive dev/staging clusters.
- **`spread_replicas_across_nodes`** maps to the operator's
  ShardAntiAffinity pod distribution: replicas of the same shard
  never co-locate on a Kubernetes node. Off by default so single-node
  dev clusters schedule; on in production, where co-located replicas
  make replication pointless against node loss.
- **`auto_inter_node_secret`** maps to the CHI `secret.auto`
  mechanism: the operator generates a shared secret that
  authenticates distributed queries between the cluster's own hosts.
  Default true, rendered only for multi-host topologies; only
  ClickHouse versions below 20.10 predate it.

### Distributed DDL after deploys and topology changes

The in-server cluster definition (`remote_servers`) reaches each host
through mounted configuration and can lag the installation reaching
its `Completed` state. `ON CLUSTER` DDL initiated in that window
silently executes on only the hosts the initiating server can see —
and still returns success (verified against a live cluster: the
distributed-DDL queue task listed one of two replicas). Before running
migrations or any distributed DDL immediately after a deploy or a
shard/replica change, confirm the initiator's view has converged:

```sql
SELECT count() FROM system.clusters WHERE cluster = '<cluster_name>'
```

must report every declared host (`shards × replicas`). Note that a
grants-constrained user needs `GRANT SELECT ON system.clusters` to
run this check — declared grants fence the system tables too.

## Design Decisions

- **One CHI, plus at most one managed CHK.** Host StatefulSets,
  Services, ConfigMaps and the PDB are all operator-created; the
  modules render the custom resources, the auth Secret, and export
  the operator's deterministic names.
- **Coordination unset means "do the right thing."** The managed
  Keeper appears exactly when the topology needs it and never
  otherwise — a dev cluster stays one pod, and turning `replicas` up
  brings coordination with it instead of a runtime replication
  failure. The CEL rule closes the remaining gap: replication without
  coordination is unrepresentable.
- **Passwords ride a module-owned Secret, never the CHI.** The CHI is
  inspectable without leaking credentials, and the rotation caveat
  (env-var delivery) is documented rather than hidden.
- **No ingress resources.** All generated Services are ClusterIP;
  exposure composes from first-class kinds over the exported
  `service_name`, and `service_annotations` is the cloud-controller
  injection surface.
- **`version` is required.** The operator's own fallback is the
  `latest` tag, which turns pod restarts into implicit upgrades — the
  spec refuses to inherit that default. `image` overrides the
  repository for mirrors; the tag still defaults to `version`.
- **Typed CR rendering on both engines.** The Pulumi module renders
  the CHI and CHK with the typed crd2pulumi SDK (field/structure
  drift against the pinned CRDs fails at compile time); the Terraform
  module applies through `kubectl_manifest` (alekc/kubectl), which
  needs no cluster connection at plan time — an infra chart can plan
  the operator and its clusters in one run, before the CRDs exist.
  Unset optionals are omitted entirely so the apiserver applies the
  CRD's own defaults; presence discipline is kept identical across
  engines.
- **No await machinery, deliberately.** Cluster readiness depends on
  the operator (image pulls, Keeper quorum, config rollout) that is
  not part of applying the resource — the never-block-on-a-controller
  posture of every operator-CR kind in the catalog.

## Version Pins and Naming Contracts

| What | Value | Notes |
|---|---|---|
| CR APIs | `clickhouse.altinity.com/v1`, `clickhouse-keeper.altinity.com/v1` | CHI and CHK |
| Operator | 0.27.2 (via KubernetesAltinityOperator) | Chart default watches ONLY its own namespace — set `watch_namespaces` |
| Server image | `clickhouse/clickhouse-server` at `spec.version` | Pin required; overridable for air-gap |
| Client Service | `clickhouse-<name>` (ClusterIP, ports 8123/9000) | Exported as `service_name` |
| Cluster Service | `cluster-<name>-<cluster_name>` | Per logical cluster |
| Host StatefulSets | `chi-<name>-<cluster_name>-<shard>-<replica>` | Single-pod each; matching headless Services |
| Managed Keeper | CHK `<name>-keeper`, client Service `keeper-<name>-keeper` | Exported as `keeper_name` / `keeper_service_name` |
| Auth Secret | `<name>-clickhouse-auth` (key = user name) | Module-owned; empty output when no users |
| `cluster_name` | Default "main", max 15 chars (upstream CRD cap) | Keep `metadata.name` ≤ 48 chars with the default |
| Ports | 9000 native, 8123 HTTP | Exported as `tcp_endpoint` / `http_endpoint` |

## IaC Twins

Pulumi (typed crd2pulumi SDK) and Terraform (`kubectl_manifest` +
null-prune locals) render identical CHI/CHK bodies — the same keys
rendered and omitted, the same coordination auto-selection, the same
Secret and output names. Keep the two modules' locals in lockstep.

## Validation Status

The component is offline-validated: spec-level tests exercise the
validation rules (topology bounds, the coordination CEL rules, format
patterns), and both engines carry offline rendering proofs that pin
the CHI/CHK bodies and output names byte-for-byte. Live end-to-end
validation against a running cluster is pending.
