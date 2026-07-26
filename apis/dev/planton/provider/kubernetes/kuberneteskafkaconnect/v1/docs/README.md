# KubernetesKafkaConnect: Research and Design

## Introduction

KubernetesKafkaConnect declares one Kafka Connect cluster on the
Strimzi `kafka.strimzi.io/v1` `KafkaConnect` custom resource,
reconciled by the Strimzi cluster operator
(KubernetesStrimziKafkaOperator must be on the cluster and watching
the namespace; the operator line is pinned at 1.1.0). Connect is the
pluggable integration engine of the Kafka ecosystem — Debezium CDC
from databases, sinks into object stores and search indexes — and the
CDC/integration arm of this catalog's Kafka family: this kind is the
worker fleet, KubernetesKafkaConnector is each declared pipe.

## Declarative Connector Management

The module stamps `strimzi.io/use-connector-resources: "true"` on the
CR unconditionally. With that annotation the operator reconciles
`KafkaConnector` resources against this cluster and REVERTS any
change made through the Connect REST API — the REST surface (exported
as `rest_api_endpoint`, port 8083 on `<name>-connect-api`) becomes
read-only inspection. This is a deliberate one-way door: mixing
declarative and REST-driven management is how Connect clusters drift,
and KubernetesKafkaConnector resources only work against clusters
carrying the annotation.

## The Four Plugin-Delivery Arms

A connector's class must be on the workers' plugin path before a
KafkaConnector can reference it. The spec models all four upstream
delivery mechanisms:

1. **The stock image** (no fields set) — ONLY the MirrorMaker 2
   connector classes (MirrorSource/MirrorCheckpoint/MirrorHeartbeat;
   verified against the workers' own plugin listing — Kafka's
   FileStream examples are NOT on the distribution's classpath).
   Zero machinery; what dev lanes and smoke tests run on, via a
   MirrorSource self-mirror.
2. **`image`** — a prebuilt Connect image already carrying plugins:
   a vendor-published image or the output of a previous `build`. The
   fastest production path when the artifact set is someone else's
   problem.
3. **`plugins`** — OCI-artifact plugins mounted as Kubernetes image
   volumes. Plugins ship as container images and mount directly into
   the workers' plugin path: no image build, no registry push, per-
   plugin versioning. The catch is CLUSTER-LEVEL: the ImageVolume
   feature must be enabled AND the container runtime must support it
   — verified live: workers fail to schedule with an image-volume
   admission error on clusters without it. The artifact `type` is
   always the literal `image`, so the module owns it rather than
   asking the author to repeat it.
4. **`build`** — the operator builds a custom image from declared
   artifacts and pushes it to `output.image`. On Kubernetes the build
   runs Kaniko (Buildah on OpenShift, where `output.type:
   imagestream` also becomes meaningful); `push_secret` names a
   docker-registry Secret in the Connect namespace with push
   credentials — the SECRET'S NAME, not a credential value. Artifacts
   are `jar`/`tgz`/`zip`/`other` URLs (with strongly-recommended
   `sha512sum` so a tampered download fails the build instead of
   running in the workers) or `maven` coordinates resolved at build
   time (Maven Central by default).

**`image` XOR `build` is CEL-enforced at the spec.** Verified in the
operator source: when `build` is configured the operator deploys the
image IT builds and a declared `image` is silently overridden —
validation turns that silent override into an authoring error.

## The Group Identity Contract

Connect workers coordinate through a group.id and three internal
Kafka topics (configs, status, offsets). The spec defaults all four
from `metadata.name` (`<name>`, `<name>-connect-configs`,
`<name>-connect-status`, `<name>-connect-offsets`) and the module
renders them EXPLICITLY into the CR — leaving them to the CRD's
defaults would derive the same values, but rendering them keeps the
uniqueness contract visible in the applied object. Two Connect
clusters (MirrorMaker 2 instances included — same protocol) sharing a
group.id or a storage topic corrupt each other's state; distinct
metadata.names give distinct identities for free.

### Storage replication factors on small clusters

Connect's internal-topic replication factor defaults to 3, which a
single-broker dev cluster cannot satisfy — the workers wedge creating
their topics. The e2e lanes set the three
`*.storage.replication.factor` entries to `"-1"` (broker default);
the presets carry the same teaching: `"-1"` for dev, `"3"` for
production.

## The Kafka Connection: Shared Client Blocks

`tls` and `authentication` use the shared Strimzi client messages
(`strimzi_kafka_client.proto`) — the same shapes MirrorMaker 2's
target and sources read, so drift between the family's connection
blocks is structurally impossible. The composition defaults encode
the common wiring: trusted-certificate Secret names default to a
KubernetesKafka's `cluster_ca_cert_secret_name` output;
authentication Secret names default to a KubernetesKafkaUser's
credential Secret (`user.crt`/`user.key` for mutual TLS, key
`password` for SCRAM). The client side supports the full SASL family
(scram-sha-512/-256, plain, custom) — wider than what listeners
accept, because external clusters (Confluent, MSK) offer mechanisms
Strimzi listeners do not.

## Worker Configuration: Strings, and Who Owns What

`config` entries are Connect worker configuration
(connect-distributed.properties) written as strings; the operator
serializes them into Java properties. The operator OWNS connection,
identity and listener configuration — entries prefixed
`bootstrap.servers`, `group.id`, the three `*.storage.topic` keys,
`ssl.`, `sasl.`, `security.`, `listeners`, `plugin.path`, `rest.`,
the interceptor keys and `prometheus.metrics.reporter.` are IGNORED
with an operator log warning (the `ssl.endpoint.identification.algorithm`
/ `ssl.cipher.suites` / `ssl.protocol` / `ssl.enabled.protocols`
client-tuning keys are the exception). `connector.plugin.version` is
additionally rejected at the spec by CEL — this Strimzi line does not
accept it in worker config; plugin versions are declared on each
KubernetesKafkaConnector's `version` field.

The secrets story for CONNECTOR config rides worker config too:
enabling `config.providers` entries (the
KubernetesSecretConfigProvider) here is what lets connectors
reference `${secrets:<namespace>/<secret>:<key>}` instead of literal
passwords — see the KubernetesKafkaConnector research doc.

## Design Decisions

- **Untyped CustomResources on both engines.** The Strimzi CRDs type
  `spec.config` (and the custom authentication's config block) with
  `x-kubernetes-preserve-unknown-fields`, which generated SDKs
  flatten into shapes that cannot hold the free-typed bodies — the
  Pulumi module renders untyped `apiextensions.CustomResource` bodies
  and the Terraform module renders `kubectl_manifest` resources, kept
  as exact twins (same keys rendered and omitted, numbers as numbers,
  booleans as booleans; the same ruling as the whole Kafka family).
- **`kubectl_manifest` over `kubernetes_manifest`.** The alekc/kubectl
  provider needs no cluster connection at plan time, so an infra
  chart can plan the operator and its Connect clusters in one run,
  before the CRDs exist.
- **No await machinery, deliberately.** Worker readiness depends on
  the operator (image pulls or an operator-driven image BUILD, group
  formation) that is not part of applying the resources — the
  never-block-on-a-controller posture of every operator-CR kind.
- **The module owns the metrics ConfigMap.** `metrics.enabled`
  renders the canonical Strimzi connect JMX Prometheus rules as
  `<name>-connect-metrics` (key `metrics-config.yml`) and wires it as
  `metricsConfig`; the CR only points at it.
- **`node_selector` translates to node affinity.** The Strimzi pod
  template carries affinity and tolerations but no nodeSelector — the
  modules render a required node affinity with one matchExpressions
  entry per label, sorted for determinism (identical on both
  engines).
- **Per-arm rendering.** Each authentication type renders only its
  own credential fields; build artifacts render exactly the declared
  type's keys (a maven artifact never carries a url and vice versa —
  CEL enforces the same partition the operator's per-type sub-schema
  expects).

## Deliberately Unmodeled

Kept off the typed spec on purpose; every item remains reachable by
declaring the raw `KafkaConnect` CR through KubernetesManifest:

- **Tracing (OpenTelemetry)** — upstream's `tracing` block (only
  OpenTelemetry remains in 1.x; Jaeger's type is gone) enables the
  agent but is useless without endpoint/service-name environment
  variables plumbed through the container template — a two-surface
  configuration (tracing block + template env) with no typed template
  surface to hang it on. All-or-nothing through KubernetesManifest.
- **`jmxOptions`** — opens an authenticated remote JMX port on the
  workers; the modeled JMX Prometheus metrics path covers
  observability without exposing a management protocol.
- **The `logging` block** — per-category log4j tuning is an
  operational debugging surface, not a deployment posture; a typed
  map would freeze log4j category names into the API.
- **`externalConfiguration`** — REMOVED upstream in 1.x (the field is
  gone from the 1.1.0 API; only a stale property-order entry
  remains). Its use cases moved to config providers (worker `config`)
  and template volumes — nothing to model.
- **Custom pod templates beyond node_selector/tolerations** — the
  upstream `template` trees are an unbounded surface; the spec models
  the two scheduling knobs that cover the real placement need.
- **`externalAccess` / listener exposure** — Connect has no client
  listener surface; its REST API is in-cluster plumbing, exported as
  an endpoint and composed behind exposure kinds if ever needed.

## Version Pins and Naming Contracts

| What | Value | Notes |
|---|---|---|
| CR API | `kafka.strimzi.io/v1` | The only served API from Strimzi 1.0 onward |
| Connector-management annotation | `strimzi.io/use-connector-resources: "true"` | Module-owned, unconditional |
| Group identity | `<name>`, `<name>-connect-{configs,status,offsets}` | metadata.name-derived defaults; unique per Kafka cluster |
| REST API Service | `<name>-connect-api` (port 8083) | Exported as `rest_api_service_name` / `rest_api_endpoint` |
| Metrics ConfigMap | `<name>-connect-metrics` (key `metrics-config.yml`) | Module-owned, rendered when `metrics.enabled` |
| `version` | empty = operator default | Strimzi 1.1 supports Kafka 4.3.0 and 4.2.1 |

## IaC Twins

Pulumi (untyped CustomResource, `module/connect.go`) and Terraform
(`kubectl_manifest` + null-prune locals, `locals.tf`) render
identical CR bodies, the same module-owned metrics ConfigMap, and the
same metadata.name-derived group identity. Keep
`locals.go`/`connect.go` and `locals.tf` in lockstep.
