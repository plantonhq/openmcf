# KubernetesKafkaMirrorMaker2: Research and Design

## Introduction

KubernetesKafkaMirrorMaker2 declares a MirrorMaker 2 replication
engine on the Strimzi `kafka.strimzi.io/v1` `KafkaMirrorMaker2`
custom resource (operator line pinned at 1.1.0). MirrorMaker 2 is
Kafka's own cross-cluster replication engine, built on Connect: each
mirror runs a MirrorSourceConnector (records + topic configuration)
and a MirrorCheckpointConnector (consumer-group offset translation).
In this catalog it is the migrate-from-Confluent/MSK story — the
engine that moves a running estate into a Strimzi-managed
KubernetesKafka cluster without stopping it.

## The 1.x CRD Shape: One Target, Per-Mirror Sources

The Strimzi 1.x API restructured KafkaMirrorMaker2 (verified in the
1.1.0 API source): the spec carries ONE `target` cluster — a typed
`KafkaMirrorMaker2TargetClusterSpec`, where the Connect-style engine
runs and keeps its state — and a `mirrors` list whose entries each
carry their OWN `source` connection. The old shape (a flat `clusters`
list with a `connectCluster` selector and per-mirror
source/target alias pairs) is gone. The spec models the new shape
directly: `target` with connection + group identity, `mirrors[]` each
with `source`, scope patterns, and two connector-tuning blocks.

Two consequences worth stating:

- **One resource mirrors INTO one cluster.** Fan-out to several
  targets is several resources — a cleaner failure and identity
  boundary than the old shared-cluster-list model.
- **The heartbeat connector is gone from the CRD on this line.** The
  1.x MirrorSpec carries only `sourceConnector` and
  `checkpointConnector` (verified in the API source) — upstream
  dropped the MirrorHeartbeatConnector block, so this spec models
  exactly two connector-tuning surfaces, not three.

## Alias Semantics and Topic Naming

Every cluster in the flow carries an alias. The DEFAULT replication
policy (`DefaultReplicationPolicy`) prefixes mirrored topics with the
source alias: `orders` from source `prod-msk` lands as
`prod-msk.orders`. That is the right default for active/active and
fan-in topologies — provenance is visible and cycles are detectable —
but wrong for migrations, where consumers expect original names.

**The IdentityReplicationPolicy recipe**: set
`replication.policy.class:
org.apache.kafka.connect.mirror.IdentityReplicationPolicy` in the
mirror's `source_connector.config` AND `checkpoint_connector.config`.
Both, always — the checkpoint connector maps offsets onto TARGET
topic names, and a policy mismatch between the two connectors
translates offsets onto topic names that do not exist. The spec
comments carry this contract; the presets demonstrate it.

Alias validation is CEL-enforced at the spec: source aliases must be
unique (they identify clusters in the flow and prefix topic names),
and the target alias must differ from every source alias — remember
the target alias DEFAULTS to `target`, so a source literally named
"target" is rejected even when the target block never sets one.

## Checkpointing: Why Cutover Is Seamless

Offsets are cluster-local: position 41523 on the source is
meaningless on the target, whose log for the same topic has its own
offsets. The MirrorCheckpointConnector continuously translates each
mirrored consumer group's committed source offsets into target-space
checkpoints; with `sync.group.offsets.enabled: "true"` it writes them
straight into the target's consumer-group state (for groups idle on
the target). At cutover, consumers repoint their bootstrap at the
target and resume from translated positions — no reprocessing, no
loss. `groups_pattern` scopes which groups checkpoint; the e2e
migration lane proves the produce-on-source → consume-on-target path
live.

## The Engine Is Connect: Group Identity and Storage Topics

MirrorMaker 2 workers form a Connect group against the TARGET
cluster. The spec's group identity fields (`target.group_id` and the
three storage topics) default from metadata.name
(`<name>`, `<name>-mirrormaker2-{configs,status,offsets}`) and share
the Connect-family uniqueness contract: unique among ALL
Connect-protocol workloads (KubernetesKafkaConnect clusters included)
sharing the target cluster.

### Storage replication factors on small clusters

The engine's internal-topic replication factor defaults to 3; on a
single-broker dev target the workers wedge creating their topics. The
e2e lanes set the three `*.storage.replication.factor` entries to
`"-1"` (broker default) through `target.config`; production targets
carry `"3"`. The presets teach the same line.

## Connection Blocks: the Shared Client Shapes

Target and every source use the shared Strimzi client messages
(`strimzi_kafka_client.proto`) — identical to KubernetesKafkaConnect's
connection block by construction. Migration postures map cleanly:

- **MSK (SCRAM)**: `scram-sha-512` + the AWS-provided CA bundle in a
  Secret (TLS trust via `pattern: "*.crt"`).
- **Confluent Cloud**: `plain` with the API key as username and the
  API secret in a Secret — SASL PLAIN inside the TLS session is
  Confluent's client contract.
- **Strimzi-to-Strimzi**: KubernetesKafka references for bootstrap
  and CA, a KubernetesKafkaUser Secret for credentials.

## Design Decisions

- **Untyped CustomResource on both engines** — the CRDs type the
  cluster and connector `config` blocks with
  `x-kubernetes-preserve-unknown-fields`; same ruling as the whole
  Kafka family (Pulumi: untyped CustomResource; Terraform:
  `kubectl_manifest`; exact twins).
- **The group identity always renders** — resolved in locals from
  spec overrides with metadata.name-derived fallbacks, so the
  uniqueness contract is visible in the applied object.
- **Replication policy stays in connector config** — a typed
  `identity_replication` boolean was considered and rejected: the
  policy class is one of several coupled connector-config knobs
  (refresh intervals, sync toggles, emit intervals), and typing one
  while leaving its siblings untyped invites half-configured
  mirrors. The presets carry the recipe instead.
- **No await machinery** — engine readiness depends on the operator
  (image pulls, worker group formation, connector startup); the
  never-block-on-a-controller posture of every operator-CR kind.
- **`node_selector` translates to node affinity** — the Strimzi pod
  template carries no nodeSelector; one matchExpressions entry per
  label, sorted, identical on both engines.
- **The module owns the metrics ConfigMap** — `<name>-mm2-metrics`,
  the canonical Strimzi connect rule set; mirroring lag rides the
  standard Connect/MirrorMaker metric families.

## Deliberately Unmodeled

Reachable by declaring the raw CR through KubernetesManifest:

- **Tracing (OpenTelemetry)** — same two-surface problem as Connect:
  the tracing block needs template env plumbing the spec does not
  carry.
- **`jmxOptions`** — remote JMX; the JMX Prometheus metrics path
  covers observability without a management protocol.
- **The `logging` block** — per-category log4j tuning; an
  operational debugging surface, not a deployment posture.
- **`clientRackInitImage` / per-connector `pause`** — narrow
  operational knobs; `rack.topology_key` and connector `state`
  semantics on the Connect side cover the modeled need.
- **Custom pod templates beyond node_selector/tolerations** — the
  unbounded upstream `template` trees.

## Version Pins and Naming Contracts

| What | Value | Notes |
|---|---|---|
| CR API | `kafka.strimzi.io/v1` | The only served API from Strimzi 1.0 onward |
| Target alias default | `target` | Must differ from every source alias (CEL-enforced) |
| Group identity | `<name>`, `<name>-mirrormaker2-{configs,status,offsets}` | Unique per Connect-protocol workload on the target cluster |
| Identity policy | `org.apache.kafka.connect.mirror.IdentityReplicationPolicy` | Set on BOTH connectors of a mirror |
| REST endpoint | `http://<name>-mirrormaker2-api.<ns>.svc.cluster.local:8083` | Exported as `rest_api_endpoint`; read-only inspection |
| Metrics ConfigMap | `<name>-mm2-metrics` (key `metrics-config.yml`) | Module-owned, rendered when `metrics.enabled` |
| `version` | empty = operator default | Strimzi 1.1 supports Kafka 4.3.0 and 4.2.1 |

## IaC Twins

Pulumi (`module/mirrormaker2.go`, untyped CustomResource) and
Terraform (`locals.tf` + `kubectl_manifest`) render identical CR
bodies — target body with resolved group identity, per-mirror
source/connector bodies, the shared client TLS/authentication
shapes, the affinity translation — plus the same module-owned metrics
ConfigMap. Keep them in lockstep.
