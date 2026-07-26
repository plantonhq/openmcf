---
title: "Simple Event Stream"
description: "This preset declares the standard application event topic: 6 partitions, replication factor 3, messages retained for 7 days and then deleted. It is the shape most topics should start from — durable..."
type: "preset"
rank: "01"
presetSlug: "01-simple-event-stream"
componentSlug: "kafka-topic"
componentTitle: "Kafka Topic"
provider: "kubernetes"
icon: "package"
order: 1
---

# Simple Event Stream

This preset declares the standard application event topic: 6 partitions, replication factor 3, messages retained for 7 days and then deleted. It is the shape most topics should start from — durable enough to survive a broker loss without losing acknowledged writes, sized with room for consumer parallelism, and bounded so the log cannot grow forever.

## When to Use

- Application event streams, notifications, task queues — anything where messages are consumed and eventually expire
- The default starting point when no special retention or compaction behavior is needed

## Key Configuration Choices

- **`partitions: 6`** -- six consumers in a group can read in parallel. Plan this up front: partitions can be INCREASED later but never decreased, and increasing them remaps keys on keyed topics
- **`replicas: 3`** -- the production norm; must not exceed the cluster's broker count (the topic operator rejects a higher value at reconcile time and the resource reports NotReady)
- **`min.insync.replicas: "2"`** -- with replication factor 3, one broker can be lost without losing acknowledged writes (producers using acks=all)
- **`retention.ms: "604800000"` + `cleanup.policy: delete`** -- messages older than 7 days are deleted; config values are Kafka configuration strings, so numbers are quoted
- **Same-namespace placement** -- `namespace` is the Kafka cluster's own namespace, deliberately: the topic operator watches only there, and a topic anywhere else is accepted and then silently never reconciled

## Values to Adapt

| Value | Description | Where to Find |
|---|---|---|
| `kafka` (namespace) | The Kafka cluster's own namespace | The KubernetesKafka resource's `namespace` |
| `my-kafka` (kafkaCluster) | The Kafka cluster's name | The KubernetesKafka resource's `metadata.name` or its `cluster_name` output |
| `order-events` | The topic name — also the Kubernetes resource name | Your naming convention; set `topicName` instead when the name needs `.`, `_`, or uppercase |

## Related Presets

- **02-compacted-changelog** -- Use when only the latest value per key matters
- **03-high-throughput** -- Use for high-volume streams with short retention
