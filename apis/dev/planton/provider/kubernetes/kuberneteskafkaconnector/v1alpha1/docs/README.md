# KubernetesKafkaConnector: Research and Design

## Introduction

KubernetesKafkaConnector declares one connector instance — a data
pipe — on the Strimzi `kafka.strimzi.io/v1` `KafkaConnector` custom
resource. The kind is deliberately thin: one CR, no namespace
creation, no satellites. Everything heavy (workers, plugins, the
Kafka connection) lives on the KubernetesKafkaConnect cluster the
connector binds to; this kind carries the pipe's class, config,
desired state, restart policy, and offset-management declarations.

## The Placement Contract

Verified against the Strimzi operator: a KafkaConnector must live in
the SAME NAMESPACE as its Connect cluster and carry the
`strimzi.io/cluster: <connect-name>` label. The operator matches
connectors by namespace + label; a connector in another namespace, or
naming a Connect cluster that does not exist there, is accepted by
the API server and then silently never reconciled — no error, no
status, nothing. The spec surfaces this as loudly as a spec can: the
`namespace` field's documentation carries the contract, and the
module renders the binding label from `connect_cluster` (a
KubernetesKafkaConnect reference resolving to its `connect_name`
output) so nothing needs manual labeling.

The deliberate corollary: this module creates NO namespace. The
namespace belongs to the KubernetesKafkaConnect resource's lifecycle;
a connector that could create its own namespace would be a connector
that can violate the placement contract.

## Declarative Management, End to End

Connect clusters from this catalog always carry
`strimzi.io/use-connector-resources: "true"`, which makes the
operator the single writer for connector configuration on the
cluster: changes made through the Connect REST API are reverted. This
resource is therefore the complete management surface — `config` is
the connector's documented key set (values as configuration strings),
`tasks_max` the parallelism ceiling (real parallelism is also bounded
by connector semantics — many CDC connectors run one task per table
set), and `state` the desired lifecycle state:

- `running` — the default.
- `paused` — tasks stay allocated, data stops moving; cheap to
  resume.
- `stopped` — tasks are deallocated; the state offset OVERRIDES
  require.

`auto_restart` (enabled + optional `max_restarts`, operator default 7
attempts with back-off capped around 30 minutes) is the production
posture: transient source outages otherwise leave the connector
FAILED until a human intervenes.

### The version field replaces a retired config entry

When workers carry several versions of a class (possible with the OCI
plugins arm), `version` pins the desired one. The old
`connector.plugin.version` config entry is deprecated upstream with
removal announced — the spec rejects it in `config` by CEL and points
authors at the typed field.

## Offset Mechanics: Declared Targets, Annotated Verbs

Strimzi models connector-offset operations as a hybrid: the CR
DECLARES the ConfigMap targets, and the operator performs the verb
when the resource carries the `strimzi.io/connector-offsets`
annotation (the annotation value selects the verb; the operator
removes the annotation when the operation completes):

- **`list_offsets.to_config_map`** — with the annotation set to
  `list`, the operator writes the connector's current offsets into
  the named ConfigMap (created if absent). Non-disruptive; the
  inspection primitive.
- **`alter_offsets.from_config_map`** — with the annotation set to
  `alter` AND the connector `stopped`, the operator applies the
  offsets read from the named ConfigMap. This is the replay
  (rewind), skip-poison-record (advance), and migration-cutover
  mechanism.

The modules render only the declarations; the annotation is an
operational verb applied with kubectl or automation — deliberately
outside IaC, which converges state rather than firing actions. Both
ConfigMap fields are foreign keys to KubernetesConfigMap for infra
charts that want the target declared too.

## Secrets in Connector Config

Connector configs routinely need credentials (database passwords,
API keys). Writing them as `config` literals would put them in the
CR, in IaC state, and in every `kubectl get kafkaconnector -o yaml`.
Connect's configuration-provider mechanism is the answer the spec
documents: reference values as
`${secrets:<namespace>/<secret>:<key>}` and enable the
KubernetesSecretConfigProvider through the CONNECT CLUSTER's worker
`config` (`config.providers: secrets`,
`config.providers.secrets.class: io.strimzi.kafka.KubernetesSecretConfigProvider`).
The workers resolve the reference at connector start; the credential
never appears in the connector resource. The provider reads Secrets
through the Connect pod's service account — RBAC on that service
account is the access boundary.

## Design Decisions

- **Untyped CustomResource on both engines.** The Strimzi CRDs type
  `config` with `x-kubernetes-preserve-unknown-fields`, which
  generated SDKs cannot carry — the same ruling as the whole Kafka
  family: Pulumi renders an untyped `apiextensions.CustomResource`,
  Terraform renders one `kubectl_manifest`, exact twins.
- **`kubectl_manifest` needs no cluster at plan time** — an infra
  chart plans the operator, the Connect cluster and its connectors in
  one run.
- **No await machinery.** Reconciliation belongs to the cluster
  operator (which drives the connector through the Connect REST API
  internally), not to applying the resource; a "class not found"
  condition surfaces on the resource status, not as an apply
  failure.
- **The binding label rides on top of the identity labels** — without
  `strimzi.io/cluster` the operator never picks the connector up, so
  the modules render it unconditionally from `connect_cluster`.

## Deliberately Unmodeled

- **Connector-level `pause` booleans from older Strimzi lines** — the
  1.x `state` field subsumes them; one typed enum instead of two
  overlapping switches.
- **Inline secret material in config** — no typed credential fields
  on this kind, deliberately: the config-provider reference pattern
  (above) is the supported path, and typed per-connector credentials
  would duplicate what Kubernetes Secrets already do.
- **The offset annotation itself** — an operational verb, not
  declarable state; modeling it would make IaC re-fire operations on
  every apply.

## Naming Contracts

| What | Value | Notes |
|---|---|---|
| CR API | `kafka.strimzi.io/v1` | The only served API from Strimzi 1.0 onward |
| Cluster binding | `strimzi.io/cluster: <connect_cluster>` | Rendered by the modules; namespace must match the Connect cluster's |
| Connector consumer group | `connect-<name>` | Kafka-side identity sink connectors consume under |
| Offset verbs | `strimzi.io/connector-offsets: list \| alter` | Annotation applied operationally; `alter` requires `state: stopped` |

## IaC Twins

Pulumi (`module/main.go`, one untyped CustomResource) and Terraform
(`locals.tf` + one `kubectl_manifest`) render identical CR bodies —
same keys rendered and omitted, tasksMax as a number, the same
autoRestart/listOffsets/alterOffsets sub-bodies. Keep them in
lockstep.
