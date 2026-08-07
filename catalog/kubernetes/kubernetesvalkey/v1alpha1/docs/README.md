# KubernetesValkey: Research and Design

## Introduction

KubernetesValkey deploys Valkey — the Linux Foundation's
Redis-compatible in-memory data store — from the OFFICIAL Valkey Helm
chart (`valkey` at https://valkey.io/valkey-helm/, pinned 0.11.0,
shipping Valkey 9.1.1). The spec types the chart's meaningful
configuration surface: two topologies (standalone and primary/replica
replication), module-rendered valkey.conf durability directives,
ACL authentication with Secret-materialized passwords, TLS from an
existing certificate Secret, and Prometheus metrics.

## From Redis to Valkey

In March 2024 Redis Ltd. moved Redis off the BSD license to a
restrictive dual license (RSALv2/SSPLv1). The community response was
Valkey: a fork of the last BSD-licensed Redis (7.2), governed by the
Linux Foundation, backed by the major cloud providers and original
Redis contributors, and permanently BSD-3-Clause. Valkey is a drop-in
replacement — the wire protocol, the command surface, and every Redis
client library work unchanged. That makes the choice simple for an
open catalog: ship Valkey, and let applications keep speaking Redis.

## Why the Official Chart

The module installs the official `valkey` chart from valkey.io — not
the Bitnami Valkey chart — for two reasons:

- **Image supply chain.** Bitnami's free public images moved to a
  latest-tag-only model with versioned images relegated to an
  unmaintained legacy registry; secure pinned images require a
  commercial subscription. The official chart consumes the
  community-maintained `docker.io/valkey/valkey` images directly.
- **Upstream alignment.** The official chart is where the Valkey
  project itself lands topology work. Its roadmap — not a third
  party's — decides when Sentinel and Cluster mode arrive, and the
  spec's honest exclusions (below) are stated in terms of that
  roadmap.

The chart pin is `0.11.0`; chart and app versions move separately and
the chart pin governs. `image` overrides registry/repository/tag for
mirrors or a specific Valkey version.

## Upstream Architecture and the Naming Contract

The Helm release and chart fullname are pinned to `metadata.name`, so
several Valkey instances coexist in one cluster and every rendered
object is predictable:

- **Standalone**: a Deployment with one pod, the write Service
  `<name>`, and optionally one PVC. No headless Service is rendered.
- **Replication**: a StatefulSet (one primary plus N replicas), the
  write Service `<name>` (targets the primary's role), the read
  Service `<name>-read` (load balances across ALL pods, replicas and
  primary), and the headless Service `<name>-headless` for direct pod
  discovery. One PVC per pod via volumeClaimTemplates.

## Topologies: Standalone vs Replication

Standalone is the default and the right shape for caches and
development: one instance, optional persistence, no replication
machinery. Replication adds read scaling and a warm copy of the
dataset — and it is important to be precise about what it does NOT
add:

- **No automated failover.** The upstream chart does not ship Sentinel
  yet. If the primary pod dies, Kubernetes restarts it and the
  replicas re-attach to the restarted primary; no replica is promoted.
  Durability through that restart comes from PERSISTENCE — the
  append-only file replayed from the primary's volume — not from
  promotion. This is why persistence is REQUIRED in replication mode:
  an ephemeral primary would replicate an empty dataset after every
  restart.
- **No Cluster mode.** Sharding across multiple primaries is likewise
  an upstream chart milestone, not modeled here.

What replication DOES give you, the spec types deliberately:

- `min_replicas_to_write` (min-replicas-to-write): the primary refuses
  writes unless at least that many replicas are connected and in sync
  — so a partitioned primary stops accepting writes the replicas would
  never see. `min_replicas_max_lag` bounds how stale a replica may be
  and still count.
- `diskless_sync`: replicas full-sync from the primary's memory
  instead of an on-disk RDB snapshot — faster when the primary's disk
  is slow, at the cost of primary memory during the transfer.
- `replication_user`: the ACL user replicas authenticate to the
  primary with (must exist in `auth.users` with `+psync +replconf
  +ping`).

## Durability Model

Valkey's durability and memory behavior live in valkey.conf, which the
chart accepts only as one raw string. The module OWNS that rendering:
the typed `config` block is turned into a deterministically ordered
config string, identical on both IaC engines, with `extra_directives`
appended verbatim for anything beyond the typed fields.

- **AOF (`append_only`)**: every write is logged and replayed on
  restart — the posture that makes a pod restart lossless when paired
  with a persistent volume. Without it, a restart recovers only the
  last RDB snapshot (or nothing, if snapshots are disabled too).
- **RDB (`rdb_save_points`)**: point-in-time snapshots as standard
  `save` directives (`"900 1"` = snapshot after 900s if at least one
  key changed). Empty rides the server's built-in schedule;
  `snapshots_disabled` renders `save ""` for pure caches and AOF-only
  postures (the two are mutually exclusive, spec-enforced).
- **Memory ceiling (`max_memory` / `max_memory_policy`)**: without
  `max_memory` the pod's memory limit is the only bound — and hitting
  it is an OOM kill, not an eviction. Always set it for caches, with
  an eviction policy (`allkeys-lru` and friends); `noeviction` (the
  server default) is right for durable stores where losing keys is
  worse than failing writes.
- **Memory headroom**: size container memory ABOVE `max_memory`.
  Fork-based persistence (RDB snapshots, AOF rewrites) momentarily
  needs copy-on-write headroom, and replication buffers live outside
  the dataset — 25–50% above `max_memory` is the working rule.

## Authentication: the ACL Model and the Secret Contract

The chart ships with auth OFF — anyone with network reach has full
read/write access. That default is never silently changed; auth is
DECLARED. The `auth.users` list creates ACL users, and the spec
enforces the chart's own warning: the list MUST include the `default`
user, because `default` is what unauthenticated clients and plain
`AUTH <password>` map to — without an explicit `default`, enabling
auth is meaningless.

Passwords never appear in rendered chart values. The module
materializes the `<name>-auth` Kubernetes Secret (one key per
username) and points the chart at it via `usersExistingSecret`. The
outputs complete the contract: `username` exports `default` and
`password_secret` exports the `{name, key}` pair applications mount or
env-inject. Per-user `permissions` strings carry the ACL rule (empty =
full access, right for the admin user; scoped rules for application
and replication users).

## TLS

`tls.enabled` serves TLS from an existing kubernetes.io/tls Secret
named by `certificate_secret` — the cert-manager seam: issue a
KubernetesCertificate and wire its `status.outputs.secret_name` here.
`require_client_certificate` turns on mutual TLS. The module does not
mint certificates itself; certificate lifecycle belongs to
cert-manager, not to a data-store module.

## Metrics

`metrics.enabled` runs the redis_exporter sidecar the chart ships
(Valkey speaks the Redis protocol, so the standard exporter works
unchanged) plus its metrics Service. `service_monitor_enabled`
additionally renders a ServiceMonitor — this requires the Prometheus
operator CRDs on the cluster, and the release FAILS to install without
them, so the flag is separate and off by default.

## Version Pins

| What | Pin | Ships |
|---|---|---|
| Chart (`valkey`, https://valkey.io/valkey-helm/) | 0.11.0 | Valkey 9.1.1 |

Chart and app versions move separately; the chart pin governs. The
Valkey image itself can be overridden per instance via `image`.

## Deliberate Exclusions

- **Sentinel HA** — not modeled because the upstream chart does not
  ship it yet; it is an upstream milestone. When it lands, failover
  becomes a typed spec concern; until then the spec refuses to promise
  a posture the chart cannot deliver.
- **Cluster mode (sharding)** — the same reasoning; also an upstream
  milestone.
- **Ingress/exposure** — composed, never embedded: the store is
  in-cluster plumbing at `kube_endpoint`, and external reachability is
  a first-class exposure kind against the exported service names. The
  write and read Services' type/annotations knobs exist for the
  managed-cloud LoadBalancer recipes only.

## Render Semantics

Both engines render the same Helm values: the module-owned config
string, the auth Secret reference, and the typed topology/service/
metrics surface. `helm_values` is the escape hatch — a YAML document
merged LAST over everything the typed fields render (Helm `-f`
semantics, identical on both engines) for the chart surface beyond the
typed fields (probes, security contexts, topology spread, network
policy, extra volumes, ...). Never for secrets.

## Production Checklist

- **Declare auth** (`auth.users` with `default` plus scoped
  application users) — auth off is a development posture only, and
  even then belongs behind a NetworkPolicy.
- **Set `max_memory`** with an explicit `max_memory_policy`; size
  container memory 25–50% above it.
- **Turn on `append_only`** with persistence wherever a restart must
  not lose data; in replication mode persistence is already required.
- **Set `min_replicas_to_write: 1+`** in replication mode so a
  partitioned primary stops accepting unreplicated writes.
- **Enable the PodDisruptionBudget** in replication mode
  (`max_unavailable: 1`) so node drains cannot take the topology down
  at once.
- **Enable metrics** and, where the Prometheus operator runs, the
  ServiceMonitor.
- **Know the failover story**: there is none — plan for
  restart-and-replay (AOF on a volume), and revisit when the upstream
  chart ships Sentinel.
