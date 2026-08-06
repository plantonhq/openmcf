# KubernetesRabbitMq: Research and Design

## Introduction

KubernetesRabbitMq declares one RabbitMQ cluster as a
`rabbitmq.com/v1beta1` `RabbitmqCluster` custom resource reconciled
by the RabbitMQ Cluster Operator (KubernetesRabbitMqOperator, the
registry prerequisite — pinned at operator v2.22.3). One resource
carries the whole cluster story: node count, storage, TLS, RabbitMQ's
own configuration vocabulary, placement, and the operator's tuning
knobs — and the operator generates the administrator credentials
itself, so they never pass through the spec.

## The Deployment Landscape

RabbitMQ is the queue-messaging broker: task queues, work
distribution, RPC, per-message acknowledgement and routing over AMQP
0-9-1/1.0, with MQTT and STOMP available as plugins. It is
deliberately NOT the event-streaming log — for append-only streams
and replayable history at scale, the catalog's answer is
KubernetesKafka. The boundary is the consumption model: a queue
delivers each message to a consumer and forgets it on acknowledgement;
a log retains everything for re-reading.

Running RabbitMQ on Kubernetes without an operator means hand-rolling
peer discovery, the Erlang cookie shared across nodes, rolling
restarts that do not lose messages, and the generated configuration
every node needs. The RabbitMQ Cluster Operator — upstream's own —
encodes all of that, so this kind is deliberately thin: it renders
ONE RabbitmqCluster (plus the optional namespace) and exports the
operator's deterministic names.

## Upstream Architecture

The operator reconciles the declared resource into:

- **The StatefulSet** (`<name>-server`) — `replicas` broker pods,
  each with its own data volume claim
  (`persistence-<name>-server-<i>`).
- **Two Services** — the client Service `<name>` (ClusterIP by
  default; the spec's `service` block switches type and carries
  annotations) and the headless inter-node Service `<name>-nodes`
  (peer discovery and stable pod DNS).
- **The generated credentials Secret** (`<name>-default-user`, keys:
  username, password, host, port, connection_string, ...) and the
  erlang-cookie Secret. The operator generates the admin credentials;
  the spec never carries them.

That naming contract — verified in the operator source at the pinned
release — is what makes the stack outputs plain strings: the modules
compute `service_name`, `headless_service_name`, both endpoints, and
`default_user_secret_name` without reading anything back from the
cluster.

### Replicas and quorum

Production clusters run an ODD replica count (3, 5, 7) so quorum
queues and the Raft-based metadata store survive node loss; a 2-node
cluster loses availability when either node fails. The operator does
NOT support scaling down — removed brokers strand their queue
replicas — so sizing down is a migration to a new cluster. The spec
documents both facts where the decision is made (`replicas`).

### The image contract

The server image must be a `-management` variant (default
`rabbitmq:4.2.6-management` at the pinned operator release, or the
operator kind's fleet-wide `default_rabbitmq_image`): the operator's
generated configuration expects the management plugin, and the plain
`rabbitmq:<version>` tags do not carry it. Three plugins are always
on — rabbitmq_management, rabbitmq_prometheus,
rabbitmq_peer_discovery_k8s — and
`configuration.additional_plugins` extends the list (changing it
rolls the cluster).

### Memory and the high watermark

RabbitMQ derives its memory high watermark from the container memory
LIMIT — so requests and limits must carry the SAME memory value, or
the broker's flow control triggers at the wrong threshold. The
operator's own defaults follow the rule (requests 1 CPU / 2Gi, limits
2 CPU / 2Gi), and the spec repeats it at the `resources` field.

### Storage

Each node gets a data PVC (`disk_size`, default 10Gi — the operator
default) carrying queues, quorum-queue Raft logs and the metadata
store, on an optional `storage_class` reference. `ephemeral: true` is
the CR's own storage-0 + emptyDir posture: every queue, message and
user vanishes with each pod restart — throwaway dev only. A CEL rule
makes the postures mutually exclusive: `ephemeral` with `disk_size`
or `storage_class` is rejected at validation.

### TLS

`tls.secret_name` references a kubernetes.io/tls Secret in the
cluster's namespace — the cert-manager seam; a KubernetesCertificate
resource's secret output plugs in directly. `ca_secret_name` (key
ca.crt) adds MUTUAL TLS. `disable_non_tls_listeners: true` closes
every plain port — AMQP moves to 5671, management to 15671, and the
plain ports of any enabled plugin (MQTT, STOMP, their WebSocket
forms) close too; the exported endpoints follow the effective ports.
The certificate must cover the client Service DNS names
(`<name>.<namespace>.svc` and its cluster-domain form).

### Placement

`node_selector` renders as REQUIRED node affinity with one In-match
per label — the RabbitmqCluster CR has no nodeSelector field, and the
affinity form is behaviorally identical for exact matches; both
engines render the same shape. `spread_across_nodes: true` renders
REQUIRED pod anti-affinity on the operator's own
`app.kubernetes.io/name: <cluster>` pod label over the
`kubernetes.io/hostname` topology — off by default so single-node dev
clusters schedule, on in production (a cluster with more replicas
than schedulable nodes sits Pending). `tolerations` are typed.

### Tuning knobs

- **`termination_grace_period_seconds`** — operator default 604800
  (7 DAYS, deliberately generous so draining nodes never lose
  messages); dev clusters lower it for fast teardown.
- **`delay_start_seconds`** — startup readiness delay (operator
  default 30, the DNS-propagation guard).
- **`skip_post_deploy_steps`** — skips post-upgrade queue
  rebalancing; only for externally managed rebalancing.
- **`auto_enable_all_feature_flags`** — enables all stable feature
  flags after upgrades, keeping the cluster upgrade-ready.

### Secret backends

The default-user credentials can bypass the operator-generated
Kubernetes Secret entirely (a oneof): **`vault`** (Kubernetes-auth
role, `default_user_path`, extra Vault agent annotations, optional
`pki_issuer_path` for Vault-issued TLS; the credential-updater
sidecar rotates), or **`external_secret_name`** (a pre-existing
Secret — e.g. one a KubernetesExternalSecret materialized from a
cloud secret manager). With the Vault arm the
`default_user_secret_name` output is empty — the credentials live at
the Vault path instead.

### Configuration layers

Typed fields cover topology, storage, TLS and placement. RabbitMQ's
own configuration vocabulary flows through the CRD's own model:
`additional_config` (rabbitmq.conf lines APPENDED to the
operator-generated file — memory watermarks, default vhost, consumer
timeouts), `advanced_config` (Erlang terms), `env_config`
(rabbitmq-env.conf; the CRD itself rejects shell command
substitution, mirrored in the spec's CEL), and `erlang_inet_config`.
Changing config or plugins triggers a rolling restart. One caveat is
documented at the field: `default_user` / `default_pass` lines in
`additional_config` override the operator-GENERATED credentials and
put the password in plaintext on the CR — migration-only.

## Design Decisions

- **One CR, everything else operator-created.** The StatefulSet,
  Services, credentials and erlang-cookie Secrets, and PVCs are all
  the operator's; the modules render the RabbitmqCluster (plus the
  optional namespace) and export the operator's deterministic names.
- **Credentials are generated, never declared.** The spec has no
  username/password fields; consumers read the `<name>-default-user`
  Secret or the declared secret backend. The one override path
  (additional_config) is documented as discouraged rather than
  hidden.
- **No ingress resources, by design.** The client Service's type and
  annotations are the cloud-exposure surface; everything else
  composes from first-class kinds over the exported handles.
- **`node_selector` becomes required node affinity.** The upstream
  CRD has no nodeSelector field; the affinity rendering is
  behaviorally identical for exact matches and kept byte-identical
  across both engines.
- **Value-based rendering for the operator-default knobs.**
  Terraform's tfvars contract flattens proto presence, so the two
  knobs whose spec defaults equal the operator defaults (termination
  grace 604800, delay start 30) render only when they DIFFER from the
  operator defaults — the common contract that keeps the CR bodies
  byte-identical across engines.
- **BACKGROUND deletion propagation, explicitly.** The operator's own
  deletion finalizer is the cascade; foreground propagation deadlocks
  against operators that keep reconciling children during deletion
  (verified live on sibling operator-owned CRs). Pulumi carries the
  `pulumi.com/deletionPropagationPolicy` annotation, Terraform sets
  `delete_cascade = "Background"`.
- **No await machinery, deliberately.** Cluster readiness depends on
  the operator (image pulls, peer discovery, feature-flag sync) that
  is not part of applying the resource — the
  never-block-on-a-controller posture of every operator-CR kind in
  the catalog.

## Version Pins and Naming Contracts

| What | Value | Notes |
|---|---|---|
| CR API | `rabbitmq.com/v1beta1` `RabbitmqCluster` | Reconciled by the RabbitMQ Cluster Operator |
| Operator | v2.22.3 (via KubernetesRabbitMqOperator) | Watches ALL namespaces by default |
| Server image | `rabbitmq:4.2.6-management` (operator default) | Must stay a `-management` variant |
| Client Service | `<name>` (ClusterIP default) | Exported as `service_name` |
| Headless Service | `<name>-nodes` | Peer discovery; exported as `headless_service_name` |
| Credentials Secret | `<name>-default-user` (keys: username, password, host, port, connection_string, ...) | Operator-generated; empty output under the Vault backend |
| StatefulSet / PVCs | `<name>-server` / `persistence-<name>-server-<i>` | Per-node data volumes |
| Ports | AMQP 5672 (TLS 5671), management 15672 (TLS 15671) | Endpoints follow the effective ports |
| Termination grace | 604800s (7 days, operator default) | Rendered only when it differs |
| Delay start | 30s (operator default) | Rendered only when it differs |

## IaC Twins

The Pulumi module renders the RabbitmqCluster with the typed
crd2pulumi SDK (field/structure drift against the pinned CRD fails at
compile time); the Terraform module applies through `kubectl_manifest`
(alekc/kubectl), which needs no cluster connection at plan time — an
infra chart can plan the operator and its clusters in one run, before
the CRD exists. Unset optionals are omitted entirely so the operator
applies its own defaults; presence discipline and the rendered CR
body are kept byte-identical across engines. Every resolution in
`module/locals.go` has an exact twin in `locals.tf` — keep them in
lockstep.

## Validation Status

The component is offline-validated: spec-level tests exercise the
validation rules (the ephemeral/storage CEL exclusion, the env_config
shell-substitution guard, format patterns), and both engines carry
plan/preview proofs across full and minimal shapes against the CRD at
the pinned operator release. Three live E2E lanes are defined — the
minimal message round-trip, the tuned full surface, and the
quorum-queue durability proof (a marker served DURING a broker pod
outage) — but have not run yet.
