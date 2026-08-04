# KubernetesKafkaTopic

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1`

**KubernetesKafkaTopicSpec** declares a Kafka topic on the Strimzi
`KafkaTopic` custom resource. The target cluster's TOPIC OPERATOR
(enabled by default on KubernetesKafka) reconciles it into a real
topic — creating it, growing partitions, and applying configuration
changes declaratively.

PLACEMENT CONTRACT (verified against the Strimzi operator): the
KafkaTopic must live in the SAME NAMESPACE as its Kafka cluster,
and it binds to the cluster through the strimzi.io/cluster label
(rendered from `kafka_cluster`). A topic in another namespace, or
naming a cluster that does not exist there, is accepted by the API
server and then silently never reconciled — set `namespace` to the
Kafka cluster's namespace.

Deleting the resource deletes the TOPIC AND ITS DATA (the topic
operator propagates deletion to Kafka).

## Example

```yaml
# Full-surface development manifest: topic-name override (characters a
# Kubernetes resource name cannot carry), explicit partitions/replicas,
# and a config set spanning retention, compaction and size bounds.
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesKafkaTopic
metadata:
  name: orders-v1
spec:
  namespace:
    value: kafka-hack
  kafkaCluster:
    value: kafka-hack
  topicName: orders.v1_events
  partitions: 12
  replicas: 3
  config:
    retention.ms: "1209600000"
    cleanup.policy: delete
    max.message.bytes: "2097152"
    min.insync.replicas: "2"
    segment.bytes: "1073741824"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespace` | `string \| valueFrom` | yes |  | KubernetesNamespace (`spec.name`) |
| `spec.kafkaCluster` | `string \| valueFrom` | yes |  | KubernetesKafka (`status.outputs.cluster_name`) |
| `spec.topicName` | `string` |  |  |  |
| `spec.partitions` | `int32` |  |  |  |
| `spec.replicas` | `int32` |  |  |  |
| `spec.config` | `map<string, string>` |  |  |  |

## Field Details

### spec.namespace

`string | valueFrom` · required

Namespace of the KafkaTopic — MUST be the Kafka cluster's own
namespace (the topic operator watches only there; see the
placement contract above). Accepts a literal namespace name or a
reference to a KubernetesNamespace resource.

- references: KubernetesNamespace (`spec.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesNamespace, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.kafkaCluster

`string | valueFrom` · required

The Kafka cluster this topic belongs to. Accepts a literal
cluster name (the KubernetesKafka resource's metadata.name) or a
reference to a KubernetesKafka resource. Rendered as the
strimzi.io/cluster label.

- references: KubernetesKafka (`status.outputs.cluster_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesKafka, name: <that resource's name>, fieldPath: status.outputs.cluster_name}} -- a bare string does not parse

### spec.topicName

`string`

Kafka topic name. Empty = the resource's metadata.name. Set it
when the topic name cannot be a Kubernetes resource name — Kafka
allows '.', '_' and uppercase (e.g. "orders.v1_DLQ"), Kubernetes
names do not. Topic names are limited to 249 characters of
alphanumerics, '.', '_' and '-'; names differing only by '.' vs
'_' collide in Kafka's internal metrics — avoid coexisting pairs.

- rule: topic name may use alphanumerics, '.', '_' and '-', up to 249 characters

### spec.partitions

`int32` · optional (explicit presence)

Number of partitions. Empty = the cluster's num.partitions
default. Partitions can be INCREASED later but never decreased
(Kafka has no partition shrink); increasing partitions changes
key-to-partition mapping for keyed topics — plan partition counts
up front for topics with semantic partitioning.

- rule: {"int32":{"gte":1}}

### spec.replicas

`int32` · optional (explicit presence)

Replication factor. Empty = the cluster's
default.replication.factor. Cannot exceed the cluster's broker
count — the topic operator rejects it at reconcile time (the
resource reports NotReady, nothing is created). 3 is the
production norm (pairs with min.insync.replicas=2 to survive one
broker loss without losing acknowledged writes).

- rule: {"int32":{"lte":32767,"gte":1}}

### spec.config

`map<string, string>`

Topic configuration (Kafka topic-level entries), e.g.
"retention.ms", "cleanup.policy", "max.message.bytes",
"min.insync.replicas". Values are Kafka configuration strings —
write numbers and booleans as strings ("604800000", "false").

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesKafkaTopic, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.namespace` | `string` | Namespace the KafkaTopic resource lives in (the Kafka cluster's namespace). |
| `status.outputs.topic_name` | `string` | The actual Kafka topic name (spec.topic_name when set, otherwise metadata.name) — what producers and consumers subscribe to. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |
| `spec.kafkaCluster` | KubernetesKafka | `status.outputs.cluster_name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
