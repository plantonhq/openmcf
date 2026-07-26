# KubernetesQdrant: Research and Design

## Introduction

KubernetesQdrant deploys one Qdrant vector-database cluster from the
official `qdrant` Helm chart (https://qdrant.github.io/qdrant-helm,
pinned 1.18.2 — the chart version tracks the Qdrant release it ships)
as a single Helm release named after `metadata.name`. The typed spec
renders into chart values; API keys ride the chart's own
existing-Secret contract so key material never lands in rendered
values; and the deployment grain is the cluster itself — collection
topology stays where Qdrant puts it, in the data plane.

## What a Vector Database Is, and Where Qdrant Sits

A vector database stores embeddings — fixed-length float arrays
produced by an ML model from text, images, audio or any other content
— and answers nearest-neighbor queries over them: "which stored
vectors are most similar to this query vector." That single primitive
is the retrieval half of RAG (retrieve the passages whose embeddings
sit closest to the question's), the engine of semantic search
(matching by meaning instead of keywords), the recall layer of
agent-memory architectures (an agent's experiences embedded and
fetched by similarity), and the candidate generator of recommendation
systems.

Qdrant (Apache-2.0) is the catalog's vector database. The
engine-level facts that matter as deployment context:

- **HNSW indexing.** Exact nearest-neighbor search is a linear scan;
  Qdrant builds an HNSW graph (Hierarchical Navigable Small World)
  over each collection's vectors for approximate search at
  logarithmic-ish cost. The index lives in memory for hot segments —
  which is why the spec tells you to size container memory for the
  vectors you plan to hold.
- **Payload filtering.** Every point carries an optional JSON payload
  and Qdrant filters DURING the vector search, not after it —
  tenant-scoped retrieval, metadata constraints and access filters
  are engine features, not application-side post-processing.
- **Quantization.** Collections can store compressed vector
  representations (scalar, product, binary) that shrink the RAM
  footprint by an order of magnitude at a recall/latency trade-off —
  declared per collection, like everything else about a collection.

All three are collection-level or query-level concerns. They shape
sizing decisions (memory, storage) but none of them is deployment
configuration — the section on the configuration split below makes
that boundary precise.

## The Deployment Landscape on Kubernetes

The official `qdrant` Helm chart is the canonical way to run
open-source Qdrant on Kubernetes. It is maintained in the qdrant-helm
repository, versioned in lockstep with the engine (chart 1.18.2 ships
app v1.18.2), and installs plain workload resources — a StatefulSet,
Services, a ConfigMap, optionally a Secret and a ServiceMonitor. No
CRDs, no operator, no admission webhooks: the chart is inert
templating, and everything it creates is inspectable with standard
tooling.

A Qdrant operator exists, but it is part of Qdrant's managed cloud
and hybrid-cloud offering rather than a standalone open-source
distribution — it is not the self-hosting path. For running OSS
Qdrant on your own cluster, the chart is both the official and the
practical answer, and it is what this component wraps.

What the chart creates per release:

- The **StatefulSet** (`replicaCount` pods, `podManagementPolicy:
  Parallel` — the chart explicitly discourages OrderedReady in
  multi-node clusters, where it deadlocks: nodes refuse to become
  Ready until all nodes are running).
- The **main ClusterIP Service** with three ports: `http` 6333
  (REST), `grpc` 6334, `p2p` 6335 — plus a headless Service
  (`-headless`) for the StatefulSet's stable pod DNS.
- The **ConfigMap** carrying `production.yaml` — the chart renders
  `.Values.config` verbatim into it — and an initialization script
  the pods run as their container command.
- The **`<name>-apikey` Secret**, only when an API key is configured.

The module pins `fullnameOverride` to `metadata.name`, so every child
name derives deterministically from the resource name and the longest
suffix stays far from the 63-character ceiling regardless of how the
release name would otherwise compose with the chart name. This is
what makes `service_name` and both endpoints exportable as plain
strings.

## Distributed Mode and the Raft Bootstrap

The chart ships `config.cluster.enabled: true` — distributed mode is
the default posture, and this component keeps it. A single replica is
simply a one-member cluster; there is no mode switch to flip later
when the workload grows.

The bootstrap mechanics: the pods run the chart's initialization
script, which gives pod 0 the role of bootstrapping the Raft
consensus and points every later pod at it over the p2p port (6335,
via the headless Service's stable DNS). Scaling out is therefore
exactly one field — `replicas` — and the new pods join the existing
consensus on their own. The chart's consensus defaults
(`config.cluster.consensus.tick_period_ms: 100`) apply unless
overridden through `helm_values`.

Two operational consequences worth stating:

- **Quorum arithmetic applies.** Raft tolerates the loss of a
  minority: 3 nodes survive 1 loss, 5 survive 2. A 2-node cluster is
  worse than 1 for availability (either loss breaks quorum). The
  production posture is 3.
- **Joining is not replicating.** A new pod joins consensus with no
  data. Collections replicate per their own `replication_factor` —
  raising `replicas` from 1 to 3 makes room for replicas; it does not
  create them. That boundary is the next section.

## Collection-Level vs Deployment-Level Configuration

This is the central design boundary of the component, and it explains
what is NOT in the spec.

In Qdrant, a collection is created through the API (REST or gRPC),
and the collection carries its own topology: `shard_number` (how many
shards the collection splits into), `replication_factor` (how many
copies of each shard), `write_consistency_factor`, quantization
config, HNSW parameters, on-disk options. These are DATA operations —
they belong to whoever owns the data model, they differ per
collection on the same cluster, and they are changed at runtime
through the same API (including resharding and replica moves).

The deployment, by contrast, owns the substrate: how many nodes
exist, how much memory and disk each has, whether the listeners
require keys, whether they speak TLS, where the pods schedule. That
is the spec.

Putting `shard_number` or `replication_factor` in the deployment spec
would be a category error on three axes:

1. **Grain** — they are per collection, and one cluster hosts many
   collections with different settings.
2. **Lifecycle** — they change through the Qdrant API at runtime; a
   manifest copy would drift from reality immediately.
3. **Authority** — the engine is the source of truth for collection
   topology; a deployment manifest that claimed it would be lying the
   moment an operator ran a reshard.

What the spec DOES own is cluster-wide collection DEFAULTS, and only
through the escape hatch: Qdrant's engine config accepts defaults
applied when a collection is created without explicit values, and the
chart renders anything under `config:` verbatim into
`production.yaml` — so those defaults ride `helm_values`, keeping the
typed spec honest about what is deployment configuration and what is
merely deployment-delivered engine tuning.

## The API-Key Auth Model and the Chart-Owned Secret Contract

Upstream Qdrant ships with authentication OFF — the listeners accept
unauthenticated requests. The component keeps that default (a private
dev cluster inside a namespace boundary is a legitimate posture) and
makes turning auth on a first-class, secret-safe declaration.

Qdrant's auth model is two static keys, not users and roles:

- The **read-write key** (`api_key`) gates everything.
- The **read-only key** (`read_only_api_key`) permits queries and
  reads only — the key to hand to query-side consumers so a leaked
  retrieval credential cannot mutate or delete collections.

The spec enforces the one combination that makes no sense: a
read-only key WITHOUT a read-write key is rejected by CEL validation,
because with no read-write key configured the cluster is
unauthenticated and the read-only key protects nothing — writes would
stay open.

Each key takes one of two arms, both riding the chart's own
mechanics:

- **`generate: true`** renders the chart value `apiKey: true` (or
  `readOnlyApiKey: true`). The chart generates a random key ONCE at
  first install and keeps it stable across upgrades via its lookup —
  a redeploy does not rotate credentials out from under every
  connected application.
- **`existing_secret {name, key}`** renders the chart's
  `valueFrom.secretKeyRef` shape. The chart reads the referenced
  Secret AT TEMPLATE TIME — the Secret must exist BEFORE the install,
  and a missing Secret fails the install rather than deploying an
  unauthenticated cluster that was declared authenticated.

Either way, the chart owns key materialization: the key material
lands in the chart-owned **`<name>-apikey`** Secret — key `api-key`
for the read-write key, `read-only-api-key` for the read-only key —
and NEVER in rendered Helm values. For the existing-secret arm the
chart copies the referenced key into its own Secret, which is why the
stack outputs export `<name>-apikey` for both arms alike. The chart
also supports a plaintext-string arm (`apiKey: "<literal>"`); the
spec deliberately does not expose it — a literal key in a manifest in
git is the exact failure mode the two typed arms exist to prevent.

The stack outputs (`api_key_secret_name`,
`read_only_api_key_secret_name`) are the composition handles:
applications mount or env-inject the Secret and the credential never
transits a manifest.

## TLS: One Typed Seam, One Escape-Hatch Surface

Qdrant has two TLS surfaces, and the component types exactly one of
them:

- **Client-listener TLS** (the typed `tls` block): the modules mount
  an existing TLS Secret at `/qdrant/tls` and render
  `config.service.enable_tls: true` with `config.tls.cert`/`key`
  pointed at the mounted `tls.crt`/`tls.key`. That enables HTTPS on
  6333 and gRPC-TLS on 6334, and the chart's probes switch scheme off
  the same flag — no probe overrides needed. The Secret shape is the
  standard Kubernetes TLS Secret (`tls.crt` + `tls.key`), which is
  exactly what cert-manager produces — the field accepts a
  KubernetesCertificate reference and consumes its issued Secret
  directly, with no key-name bridging.
- **Inter-node p2p TLS** (`config.cluster.p2p.enable_tls` plus its
  cert config): a separate upstream surface with its own certificate
  distribution concerns, left on `helm_values`. In-cluster p2p
  between pods of one StatefulSet inside one namespace is the 10% of
  the 90/10 split.

Empty `tls` means plaintext in-cluster; the exported `http_endpoint`
scheme follows the block (`http` vs `https`). TLS for external
clients composes at the exposure layer instead.

## Storage and the Snapshots Volume

Each pod gets its own PersistentVolumeClaim — vectors, payloads, and
the write-ahead log live on the `storage` volume (spec default 10Gi;
the chart's persistence block always renders a PVC). `storage_class`
is set only when declared; empty means the cluster's default class.

The optional `snapshots` volume exists because of a specific failure
mode upstream documents in the chart values themselves: snapshots
(backups, and the snapshot-based shard transfer mechanism used when
shards move between nodes) are written to disk, and a snapshot of a
large collection written onto the data volume can fill it and crash
the node. Upstream's recommendation, carried into the spec: when
using snapshots, declare the separate volume and size it like the
main volume. The chart also notes the snapshots StorageClass is a
good place for cold storage — the spec exposes it. When `snapshots`
is not declared, the chart puts snapshots on the storage volume,
which is fine until the first big one.

Sizing memory belongs here too: Qdrant keeps hot segments and HNSW
indexes in RAM, so the `resources` block — not the storage size — is
what bounds how many vectors serve at full speed. Quantization and
on-disk options (collection-level) stretch that budget; the
deployment's job is to declare it honestly.

## Scaling, Anti-Affinity and Scheduling

`replicas` is the only scaling control, by design (see the Raft
section). The `scheduling` block carries the placement surface:

- **`node_selector`** and **`tolerations`** — the standard pinning
  primitives, passed through to the chart's values.
- **`pod_anti_affinity: true`** renders the chart's own documented
  anti-affinity recipe (the commented example in the chart values):
  REQUIRED anti-affinity on `kubernetes.io/hostname` matching the
  release's own pods — every member on a different node, so one node
  loss takes one member, not the quorum. Required (not preferred)
  because a production cluster that silently co-schedules two of
  three members has an availability property it does not know it
  lost. Meaningful from 2 replicas up; the chart default is none.
- **`priority_class_name`** — eviction ordering for clusters that run
  priority tiers.

Beyond the typed block, `topologySpreadConstraints` (zone-level
spreading) and the chart's PodDisruptionBudget (`podDisruptionBudget`
— disabled by default upstream) ride `helm_values`.

## Monitoring

Qdrant exposes Prometheus metrics on `/metrics` over the REST
listener. `service_monitor_enabled: true` renders the chart's
`metrics.serviceMonitor.enabled` — a ServiceMonitor scraping the http
port at the chart's defaults (30s interval). The flag requires the
Prometheus Operator CRDs to exist in the cluster; without them the
install fails on the unknown kind, which is why the flag is opt-in
and defaults to the chart's own false. ServiceMonitor tuning
(relabelings, scrape authorization when API keys guard the metrics
endpoint) stays on `helm_values`.

## What the Typed Spec Covers vs What Rides helm_values

The typed spec is the chart's meaningful configuration surface —
the decisions every real deployment makes:

| Concern | Typed field |
|---|---|
| Placement | `namespace`, `create_namespace` |
| Version | `chart_version` |
| Cluster size | `replicas` |
| Sizing | `resources` |
| Data volume | `storage` (size, class) |
| Snapshots volume | `snapshots` (size, class) |
| Auth | `api_key`, `read_only_api_key` |
| Client TLS | `tls` |
| Placement detail | `scheduling` |
| Metrics | `service_monitor_enabled` |
| Air-gap / hardened image | `image` (repository, tag, `use_unprivileged`) |

Everything else rides `helm_values`, merged LAST with Helm `-f`
semantics (maps deep-merge with the later document winning, lists
replace — identical on both engines):

- The **`config:` engine document** — collection defaults, optimizer
  and WAL tuning, consensus tick period, p2p TLS. The chart renders
  it verbatim into `production.yaml`.
- **Probes** (the chart ships readiness on, liveness/startup off),
  **PodDisruptionBudget**, **topology spread constraints**,
  **lifecycle** (the chart's default preStop sleep),
  **sidecarContainers**, **additionalVolumes/-Mounts**, **env**,
  service annotations and type changes, image pull secrets.
- **snapshotRestoration** — the chart's restore-from-snapshot path,
  an operational maneuver rather than steady-state configuration.

Never secrets: key material belongs in the `api_key` /
`read_only_api_key` arms, where the Secret contract keeps it out of
rendered values and state.

## The 90/10 Rationale

The split follows one test: does the setting decide what gets
deployed (typed), or does it tune how the deployed thing behaves
(escape hatch)?

Ninety percent of real deployments are fully described by the typed
fields — a namespace, a size, volumes, keys, TLS, placement, metrics.
Those fields are validated (quantity formats, the replicas floor, the
read-only-requires-read-write rule), defaulted (chart version,
volume sizes), and cross-referenced (namespace, StorageClass,
certificate references) — the platform can reason about them.

The remaining ten percent is real but long-tail: a team that tunes
`default_segment_number`, runs a backup sidecar, or enables p2p TLS
knows exactly which chart value it wants, and `helm_values` hands it
the whole chart surface without the spec growing a field per knob.
Merged last, it can even override a typed rendering in an emergency —
a safety valve, never the primary interface. The one thing it must
never carry is a secret, and the docs say so at every mention.

## Version Pins and Naming Contracts

| What | Value | Notes |
|---|---|---|
| Chart | `qdrant` at https://qdrant.github.io/qdrant-helm | Pinned 1.18.2; chart version = app version |
| Release / Service name | `metadata.name` | `fullnameOverride` pinned; headless twin `-headless` |
| Ports | REST 6333, gRPC 6334, p2p 6335 | SDKs default to gRPC; p2p is cluster-internal |
| API-key Secret | `<name>-apikey`, keys `api-key` / `read-only-api-key` | Chart-owned; populated for both arms; exported in outputs |
| TLS mount | `/qdrant/tls` (`tls.crt`/`tls.key`) | cert-manager Secret shape, consumed directly |
| Storage | 10Gi spec default, PVC per pod | `snapshots` volume optional; size like storage |
| Image | `docker.io/qdrant/qdrant` at chart appVersion | `use_unprivileged` appends `-unprivileged` to the TAG (e.g. `v1.18.2-unprivileged`) and skips the root-owned ownership init container |
| Pod management | `Parallel` (chart default) | OrderedReady deadlocks multi-node bootstrap |

## IaC Twins

Pulumi (`module/values.go` + `module/locals.go`) and Terraform
(`locals.tf` + `helm_release`) render identical chart values — the
same fullnameOverride, the same api-key arm shapes, the same TLS
mount constants, the same anti-affinity recipe — and export the same
outputs. The Terraform module reaches the helm_values merge natively
(values = [yamlencode(typed), helm_values], provider-merged in that
order); the Pulumi module merges in Go with the same semantics. Keep
the typed-value rendering in lockstep.
