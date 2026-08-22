# Compacted Changelog Topic

This preset creates a keyed changelog (latest-value-per-key) topic: pure compaction with an eager cleaner and a hard one-week compaction deadline, timestamped by the broker for consistent ordering semantics.

## When to Use

- Entity snapshots (customers, products, account state) where consumers need the CURRENT value per key, not history
- Backing topics for state stores and table-style consumers

## Key Configuration Choices

- **`cleanupPolicy: compact`** -- Kafka keeps the newest record per key; deletes ride tombstones. Every producer record needs a key, or compaction cannot do its job.
- **`minCleanableDirtyRatio: 0.4`** -- compacts earlier than Kafka's 0.5 default, trading a little CPU for a tighter topic.
- **`maxCompactionLagMs` (7 days)** -- an upper bound on how stale an uncompacted record may get, whatever the ratio says.
- **`messageTimestampType: log_append_time`** -- broker-assigned timestamps keep ordering semantics uniform across producers.

## What You Get

A replicated, self-compacting current-state topic on your referenced Kafka cluster -- storage bounded by keyspace size instead of time.
