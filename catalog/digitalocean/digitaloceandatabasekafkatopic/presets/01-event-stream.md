# Event Stream Topic

This preset creates a durable event-stream topic: seven-day time-based retention, three replicas with `min_insync_replicas: 2` (pair with producers sending `acks=all` for real durability), and six partitions of headroom. Compression is left to producers.

## When to Use

- Order/event streams consumed by one or more services within a bounded replay window
- Any delete-policy topic where losing acknowledged writes is not acceptable

## Key Configuration Choices

- **`retentionMs: 604800000` (7 days)** -- bounds disk growth on the cluster; consumers get a week of replay.
- **`replicationFactor: 3` + `minInsyncReplicas: 2`** -- survives a broker loss without losing acknowledged writes (needs a 3-node cluster; the replication ceiling is the cluster's node count).
- **`partitionCount: 6`** -- parallelism headroom; remember partitions can only ever be ADDED.
- **`cleanupPolicy: delete` stated explicitly** -- the provider would otherwise seed `compact_delete` when any config block is present.

## What You Get

A replicated, time-bounded event topic on your referenced Kafka cluster, addressed by its `topic_name` output.
