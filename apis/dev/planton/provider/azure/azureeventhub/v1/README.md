# AzureEventHub

An event hub inside an Azure Event Hubs namespace: one partitioned,
replayable event stream. Producers append events to partitions;
consumers read through consumer groups, each keeping its own offset --
so the same stream feeds real-time processing, batch analytics, and
archival independently. Consumer groups, hub-scoped SAS rules, and
hub-level data-plane RBAC compose onto it as first-class kinds.

## When to Use

Use AzureEventHub when you need:

- **A replayable event stream** -- telemetry, logs, clickstreams,
  change-data capture; consumers that fall behind or need to reprocess
  rewind within the retention window
- **A Kafka topic on Azure** -- the hub's name IS the topic name on the
  namespace's Kafka endpoint (STANDARD and above)
- **Streaming-to-batch archival** -- the folded capture block lands
  every event in Blob Storage as Avro on a size-or-interval cadence,
  with no consumer application to run
- **Changelog/table semantics** -- Kafka-style log compaction keeps the
  latest event per partition key forever (COMPACT cleanup policy)

## Key Configuration

- `event_hub_name` -- unique within the namespace; the Kafka topic name
  (ForceNew; renaming replaces the hub and its retained events)
- `partition_count` -- the unit of parallelism and ordering (1-32 on
  shared namespaces, up to 1024 on PREMIUM/dedicated); never
  decreasable, increasable only on PREMIUM/dedicated -- size for peak
  up front on shared namespaces
- `message_retention` XOR `retention_description` -- simple days (tier
  caps 1/7/90) or the richer model: hour-granular DELETE windows, or
  COMPACT with a tombstone window (cleanup policy fixed at creation)
- `status` -- ACTIVE (default), DISABLED, or SEND_DISABLED
  (receive-only drain mode)
- `capture_description` -- AVRO/AVRO_DEFLATE encoding, 60-900s window
  interval, 10-500 MB window size, and a Blob Storage destination with
  service-managed SAS or managed-identity authentication

## Composition

```yaml
namespaceId:
  valueFrom:
    kind: AzureEventHubNamespace
    name: telemetry-hubs
    fieldPath: status.outputs.namespace_id
```

Consumer groups and hub-scoped authorization rules reference
`status.outputs.event_hub_id`; it is also the scope for hub-level
data-plane role grants (Azure Event Hubs Data Receiver/Sender on exactly
this hub).

## Documentation

- [Design research](docs/README.md) -- field mapping, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)
