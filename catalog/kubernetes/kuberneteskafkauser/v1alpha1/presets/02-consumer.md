# Consumer (SCRAM)

This preset declares a consumer identity: SCRAM-SHA-512 authentication with an operator-generated password, and the two-part ACL grant every consumer needs — Read/Describe on the topic it consumes AND Read on the consumer group it coordinates under. Forgetting the group half is the classic consumer-ACL mistake: the client authenticates, subscribes, and then fails on group coordination.

## When to Use

- A service that consumes one topic under one consumer group and writes nothing
- SCRAM listeners — username/password is the simplest client wiring, the common choice for in-cluster application traffic

## Key Configuration Choices

- **`type: scram-sha-512`** -- must match a scram-sha-512 listener on the target cluster; the generated Secret carries `password` plus a ready `sasl.jaas.config`
- **Topic `Read` + `Describe`** -- the canonical consumer grant on the topic side (matches the upstream consumer example)
- **Group `Read`** -- consumer-group coordination is its own resource type; the group id here must match the client's `group.id`
- **Literal matching** (the default) -- one topic, one group; switch to `patternType: prefix` for topic or group families
- **ACLs enforced only with cluster-side `simple` authorization** -- declared against a cluster without it, the rules are rejected at reconcile and the resource reports NotReady

## Values to Adapt

| Value | Description | Where to Find |
|---|---|---|
| `kafka` (namespace) | The Kafka cluster's own namespace | The KubernetesKafka resource's `namespace` |
| `my-kafka` (kafkaCluster) | The Kafka cluster's name | The KubernetesKafka resource's `metadata.name` or its `cluster_name` output |
| `order-events` | The topic this consumer reads | The KubernetesKafkaTopic resource's `topic_name` output |
| `order-processing` | The consumer group id | The client application's `group.id` configuration |

## Related Presets

- **01-producer** -- The write-side counterpart, with a prefix ACL
- **03-mtls-service** -- Certificate-based identity with quotas, for tls-auth listeners
