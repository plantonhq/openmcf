---
title: "High Throughput"
description: "This preset declares a wide, short-lived firehose topic: 24 partitions for consumer parallelism, 1 GiB segments so high-volume partitions roll less often, and a 2-day retention window because..."
type: "preset"
rank: "03"
presetSlug: "03-high-throughput"
componentSlug: "kafka-topic"
componentTitle: "Kafka Topic"
provider: "kubernetes"
icon: "package"
order: 3
---

# High Throughput

This preset declares a wide, short-lived firehose topic: 24 partitions for consumer parallelism, 1 GiB segments so high-volume partitions roll less often, and a 2-day retention window because telemetry-class data has a short shelf life. The durability posture (replication factor 3, `min.insync.replicas` 2) stays the same as the other presets — throughput tuning changes partitioning and log mechanics, not replication.

## When to Use

- Telemetry, metrics, clickstream, log ingestion — high message volume, consumed near real time
- Streams where a large consumer group needs wide parallel read fan-out
- Data that is aggregated or shipped elsewhere quickly, so long retention only burns disk

## Key Configuration Choices

- **`partitions: 24`** -- up to 24 consumers in one group read in parallel. Plan the count up front: partitions can be INCREASED later but never decreased, and increasing them remaps keys on keyed topics
- **`retention.ms: "172800000"`** -- 2 days; whichever retention threshold Kafka reaches first triggers cleanup, and short retention keeps disk usage proportional to recent volume
- **`segment.bytes: "1073741824"`** -- 1 GiB segments (Kafka's broker-level default made explicit): larger active segments contain more messages and are rolled less often
- **`replicas: 3` + `min.insync.replicas: "2"`** -- unchanged from the default posture; must not exceed the cluster's broker count

## Values to Adapt

| Value | Description | Where to Find |
|---|---|---|
| `kafka` (namespace) | The Kafka cluster's own namespace | The KubernetesKafka resource's `namespace` |
| `my-kafka` (kafkaCluster) | The Kafka cluster's name | The KubernetesKafka resource's `metadata.name` or its `cluster_name` output |
| `telemetry-ingest` | The topic name — also the Kubernetes resource name | Your naming convention; set `topicName` instead when the name needs `.`, `_`, or uppercase |

## Related Presets

- **01-simple-event-stream** -- The balanced default for ordinary event streams
- **02-compacted-changelog** -- Use when only the latest value per key matters
