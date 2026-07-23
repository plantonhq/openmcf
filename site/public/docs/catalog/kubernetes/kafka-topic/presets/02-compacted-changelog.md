---
title: "Compacted Changelog"
description: "This preset declares a compacted topic: instead of deleting messages by age, Kafka's log cleaner retains the LATEST value for each message key and discards older versions. The topic behaves as a..."
type: "preset"
rank: "02"
presetSlug: "02-compacted-changelog"
componentSlug: "kafka-topic"
componentTitle: "Kafka Topic"
provider: "kubernetes"
icon: "package"
order: 2
---

# Compacted Changelog

This preset declares a compacted topic: instead of deleting messages by age, Kafka's log cleaner retains the LATEST value for each message key and discards older versions. The topic behaves as a durable, replayable table of current state — the shape for entity snapshots, connector offsets, and any keyed stream where only the most recent update matters.

Records must be keyed: Kafka guarantees the latest message per key is retained, and a null-valued message (a tombstone) marks a key for deletion.

## When to Use

- Changelog / snapshot topics where consumers rebuild current state by replaying from the beginning
- Keyed streams where old versions of a value are noise, not history
- NOT for event streams where every occurrence matters — use the delete policy (01-simple-event-stream) there

## Key Configuration Choices

- **`cleanup.policy: compact`** -- the topic-level switch to compaction; the cluster-wide default stays `delete`
- **`min.cleanable.dirty.ratio: "0.4"`** -- compaction kicks in once 40% of the log is uncompacted, keeping the latest-value guarantee fresher than the laxer default
- **`segment.bytes: "268435456"`** -- 256 MiB segments, deliberately smaller than Kafka's 1 GiB default: the ACTIVE segment is never compacted, so rolling segments sooner makes records eligible for compaction sooner
- **`partitions: 6` / `replicas: 3` / `min.insync.replicas: "2"`** -- the same durability posture as the event-stream preset; compaction changes what is retained, not how it is replicated

## Values to Adapt

| Value | Description | Where to Find |
|---|---|---|
| `kafka` (namespace) | The Kafka cluster's own namespace | The KubernetesKafka resource's `namespace` |
| `my-kafka` (kafkaCluster) | The Kafka cluster's name | The KubernetesKafka resource's `metadata.name` or its `cluster_name` output |
| `customer-profiles` | The topic name — also the Kubernetes resource name | Your naming convention; set `topicName` instead when the name needs `.`, `_`, or uppercase |

## Related Presets

- **01-simple-event-stream** -- The delete-policy default for event streams
- **03-high-throughput** -- Use for high-volume streams with short retention
