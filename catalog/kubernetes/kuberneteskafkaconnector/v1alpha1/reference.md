# KubernetesKafkaConnector

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1alpha1`

**KubernetesKafkaConnectorSpec** declares one data pipe — a
connector instance — on the Strimzi `KafkaConnector` custom
resource. The target Connect cluster's operator-managed workers run
it: a SOURCE connector streams an external system into Kafka
(Debezium CDC from Postgres/MySQL, file tails, SaaS APIs); a SINK
connector streams Kafka topics out (object stores, search indexes,
warehouses).

PLACEMENT CONTRACT (verified against the Strimzi operator): the
KafkaConnector must live in the SAME NAMESPACE as its Connect
cluster, and it binds to the cluster through the strimzi.io/cluster
label (rendered from `connect_cluster`). A connector in another
namespace, or naming a Connect cluster that does not exist there,
is accepted by the API server and then silently never reconciled —
set `namespace` to the Connect cluster's namespace.

The connector's CLASS must be available on the Connect cluster's
workers — via the stock image, a prebuilt `image`, OCI `plugins`,
or a `build` on the KubernetesKafkaConnect resource. A class the
workers do not carry fails at reconcile with a "class not found"
condition on the resource.

## Example

```yaml
# Full-surface development manifest: explicit task ceiling, plugin
# version pin, a MirrorHeartbeatConnector config set (one of the stock image's three classes),
# desired state, auto-restart policy, and both offset ConfigMap
# declarations (list target + alter source, literal ConfigMap names).
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesKafkaConnector
metadata:
  name: file-source-hack
spec:
  namespace:
    value: kafka-hack
  connectCluster:
    value: connect-hack
  connectorClass: org.apache.kafka.connect.mirror.MirrorHeartbeatConnector
  tasksMax: 2
  version: "4.3.0"
  config:
    file: /opt/kafka/LICENSE
    topic: license-lines
    key.converter: org.apache.kafka.connect.storage.StringConverter
    value.converter: org.apache.kafka.connect.storage.StringConverter
  state: running
  autoRestart:
    enabled: true
    maxRestarts: 5
  listOffsets:
    toConfigMap:
      value: file-source-offsets-list
  alterOffsets:
    fromConfigMap:
      value: file-source-offsets-override
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespace` | `string \| valueFrom` | yes |  | KubernetesNamespace (`spec.name`) |
| `spec.connectCluster` | `string \| valueFrom` | yes |  | KubernetesKafkaConnect (`status.outputs.connect_name`) |
| `spec.connectorClass` | `string` | yes |  |  |
| `spec.tasksMax` | `int32` |  |  |  |
| `spec.version` | `string` |  |  |  |
| `spec.config` | `map<string, string>` |  |  |  |
| `spec.state` | `string` |  | `running` |  |
| `spec.autoRestart` | `KubernetesKafkaConnectorAutoRestart` |  |  |  |
| `spec.autoRestart.enabled` | `bool` |  |  |  |
| `spec.autoRestart.maxRestarts` | `int32` |  |  |  |
| `spec.listOffsets` | `KubernetesKafkaConnectorListOffsets` |  |  |  |
| `spec.listOffsets.toConfigMap` | `string \| valueFrom` | yes |  | KubernetesConfigMap (`status.outputs.configmap_name`) |
| `spec.alterOffsets` | `KubernetesKafkaConnectorAlterOffsets` |  |  |  |
| `spec.alterOffsets.fromConfigMap` | `string \| valueFrom` | yes |  | KubernetesConfigMap (`status.outputs.configmap_name`) |

## Field Details

### spec.namespace

`string | valueFrom` · required

Namespace of the KafkaConnector — MUST be the Connect cluster's
own namespace (see the placement contract above). Accepts a
literal namespace name or a reference to a KubernetesNamespace
resource.

- references: KubernetesNamespace (`spec.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesNamespace, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.connectCluster

`string | valueFrom` · required

The Connect cluster this connector runs on. Accepts a literal
name (the KubernetesKafkaConnect resource's metadata.name) or a
reference to a KubernetesKafkaConnect resource. Rendered as the
strimzi.io/cluster label.

- references: KubernetesKafkaConnect (`status.outputs.connect_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesKafkaConnect, name: <that resource's name>, fieldPath: status.outputs.connect_name}} -- a bare string does not parse

### spec.connectorClass

`string` · required

Fully-qualified connector class, e.g.
"io.debezium.connector.postgresql.PostgresConnector". The class
decides whether this is a source or sink connector. The class
must be ON the Connect cluster's workers: the stock image carries
ONLY the MirrorMaker 2 connectors (verified live against the
workers' own plugin listing — Kafka's FileStream examples are not
on the classpath); anything else arrives via the Connect
cluster's image, plugins, or build arms.

- rule: {"required":true}

### spec.tasksMax

`int32` · optional (explicit presence)

Maximum number of parallel tasks the connector may run.
Empty = the Connect default (1). Tasks spread across the Connect
cluster's workers; a connector's real parallelism is also
bounded by its own semantics (e.g. one task per table set for
many CDC connectors).

- rule: {"int32":{"gte":1}}

### spec.version

`string`

Desired connector plugin version when the workers carry several
versions of the class (e.g. "3.1.0" or a Maven-style range).
Empty = the newest available. This field replaces the retired
connector.plugin.version config entry.

### spec.config

`map<string, string>`

The connector's own configuration — the exact key set its
documentation defines (connection URLs, table filters, topics,
converters). Values are configuration strings — write numbers
and booleans as strings ("5432", "false").

SECRETS in connector config (database passwords, API keys)
should not be written as literals: Connect supports configuration
providers — reference values as
${secrets:<namespace>/<secret>:<key>} after enabling the
KubernetesSecretConfigProvider through the Connect cluster's
`config` (config.providers entries).

- rule: connector.plugin.version in config is deprecated upstream with removal announced — set the connector version on the version field instead

### spec.state

`string` · optional (explicit presence)

Desired connector state: "running" (default), "paused" (retains
tasks, stops moving data), or "stopped" (deallocates tasks).

- default: `running`
- rule: state must be running, paused, or stopped

### spec.autoRestart

`KubernetesKafkaConnectorAutoRestart`

Automatic restart of failed connectors/tasks with incremental
back-off. Strongly recommended for production pipes — transient
source-system outages otherwise leave the connector FAILED until
a human intervenes.

### spec.autoRestart.enabled

`bool`

Enable automatic restarts of failed connectors and tasks.

### spec.autoRestart.maxRestarts

`int32` · optional (explicit presence)

Give up after this many consecutive restarts. Empty = the
operator default (7 attempts, back-off capped at ~30 minutes).

- rule: {"int32":{"gte":1}}

### spec.listOffsets

`KubernetesKafkaConnectorListOffsets`

Offset inspection: write the connector's current offsets to the
named ConfigMap when the resource carries the
strimzi.io/connector-offsets: list annotation (an operational
verb — the ConfigMap target is declared here, the action is
triggered by annotation).

### spec.listOffsets.toConfigMap

`string | valueFrom` · required

ConfigMap the offsets are written TO. Accepts a literal name or
a reference to a KubernetesConfigMap resource in the same
namespace (created if absent when the list action runs).

- references: KubernetesConfigMap (`status.outputs.configmap_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesConfigMap, name: <that resource's name>, fieldPath: status.outputs.configmap_name}} -- a bare string does not parse

### spec.alterOffsets

`KubernetesKafkaConnectorAlterOffsets`

Offset override: apply offsets from the named ConfigMap when the
resource carries the strimzi.io/connector-offsets: alter
annotation while the connector is stopped — the replay /
skip-poison-record / migration-cutover mechanism.

### spec.alterOffsets.fromConfigMap

`string | valueFrom` · required

ConfigMap the new offsets are read FROM. Accepts a literal name
or a reference to a KubernetesConfigMap resource in the same
namespace.

- references: KubernetesConfigMap (`status.outputs.configmap_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesConfigMap, name: <that resource's name>, fieldPath: status.outputs.configmap_name}} -- a bare string does not parse

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesKafkaConnector, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.namespace` | `string` | Namespace the KafkaConnector resource lives in (the Connect cluster's namespace). |
| `status.outputs.connector_name` | `string` | The connector's name inside the Connect cluster (metadata.name) — what the Connect REST API and consumer-group names (connect-<name>) key off. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |
| `spec.connectCluster` | KubernetesKafkaConnect | `status.outputs.connect_name` |
| `spec.listOffsets.toConfigMap` | KubernetesConfigMap | `status.outputs.configmap_name` |
| `spec.alterOffsets.fromConfigMap` | KubernetesConfigMap | `status.outputs.configmap_name` |

## See Also

- [Overview](../README.md)
