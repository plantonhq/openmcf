---
title: "Compacted Changelog Stream"
description: "This preset creates a log-compacted hub: the latest event per partition key is kept forever -- Kafka-style compacted-topic semantics for entity changelogs, materialized views, and cache warming."
type: "preset"
rank: "03"
presetSlug: "03-compacted-changelog"
componentSlug: "event-hub"
componentTitle: "Event Hub"
provider: "azure"
icon: "package"
order: 3
---

# Compacted Changelog Stream

This preset creates a log-compacted hub: the latest event per partition
key is kept forever -- Kafka-style compacted-topic semantics for entity
changelogs, materialized views, and cache warming.

## When to Use

- Change-data-capture streams where consumers need the CURRENT state of
  every key, not the full history
- Rebuilding caches/read models by replaying the compacted log from the
  beginning

## Key Configuration Choices

- **`cleanupPolicy: COMPACT` is ForceNew** -- a hub cannot switch
  between delete and compact retention after creation
- **`tombstoneRetentionTimeInHours: 24`** -- consumers must replay
  within 24 hours to observe deletions (a tombstone is a null-value
  event marking a key deleted); size it to your slowest consumer's
  catch-up window
- **STANDARD tier or higher** -- Azure rejects compaction on BASIC
  namespaces at apply time

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `entity-changelog` | The hub name | Your stream taxonomy |
| `24` | The tombstone window in hours | Your consumer lag budget |
