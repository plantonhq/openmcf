# Azure Event Streaming Platform

Kafka-class event streaming with the pieces production actually requires: an auto-inflating Event Hubs namespace, separate raw and processed hubs, per-application consumer groups, a schema registry governing producer/consumer contracts, and keyless capture archiving every event to Data Lake storage — plus the two alerts (throttling, capture backlog) that tell you when ingestion or archival is quietly failing.

## Who this is for

A data team ingesting events at scale — clickstreams, telemetry, transactions — that needs replayable streams today and a lakehouse landing zone tomorrow. Assembling streaming with schema governance and a durable archive path is a week of wiring done by hand; this chart deploys the whole shape, with every seam between the pieces expressed as a reference.

## Architecture

```
 producers ──▶ [ingest hub]  ── stream-processor CG ──▶ processor ──▶ [processed hub] ── dashboard CG ──▶ consumers
 (schemas       │  8 partitions, 7d retention                           4 partitions
  from the      │
  registry)     └── capture (keyless: namespace's own identity writes)
                        │  Avro, 5-min/300MB windows
                        ▼
               ADLS Gen2 archive  ◀── Storage Blob Data Contributor grant
               (HNS storage account: the lakehouse landing zone)

 schema group (AVRO, BACKWARD)      alerts: ThrottledRequests > 0 (sev 1)
                                            CaptureBacklog > 500MB (sev 2)
```

Three design decisions worth knowing:

- **Raw and processed are separate hubs, not consumer groups.** Consumer groups are offset cursors over the *same* data; the processor re-publishes *different* (enriched) data. Downstream consumers can never accidentally read raw.
- **Capture is keyless.** The namespace carries a system-assigned identity, capture authenticates as `SYSTEM_ASSIGNED`, and a role assignment on the archive account is the entire credential story — no storage key is minted for archival, and revocation is one grant.
- **The archive account is Data Lake Gen2** (hierarchical namespace), so Spark/Synapse/Databricks mount the Avro archives as a filesystem. Capture's output *is* the lakehouse's raw zone.

## Resources

| Kind | Name | Purpose |
| --- | --- | --- |
| AzureResourceGroup | `{env}-streaming` | One container for the estate |
| AzureLogAnalyticsWorkspace | `{env}-streaming-logs` | Broker logs for post-alert forensics |
| AzureMonitorDiagnosticSetting | `{env}-streaming-ns-diag` | Namespace logs + metrics into the workspace |
| AzureEventHubNamespace | `{env}-streams` | STANDARD namespace, system identity, auto-inflate |
| AzureEventHub | `{env}-ingest-hub` | Raw ingestion, capture-enabled |
| AzureEventHub | `{env}-processed-hub` | Enriched events, post-processor |
| AzureEventHubConsumerGroup | `{env}-processor-cg` / `{env}-dashboard-cg` | Independent offset cursors per application |
| AzureEventHubSchemaGroup | `{env}-event-schemas` | Avro schema registry, BACKWARD compatibility |
| AzureStorageAccount | `{env}-event-archive` | ADLS Gen2 capture archive (capture toggle) |
| AzureStorageContainer | `{env}-event-archive-container` | Where archives land (capture toggle) |
| AzureRoleAssignment | `{env}-capture-writer` | The keyless capture grant (capture toggle) |
| AzureEventHubAuthorizationRule | `{env}-ingest-send` / `{env}-processed-listen` | Per-role SAS credentials (absent in keyless-only mode) |
| AzureMonitorActionGroup | `{env}-streaming-ops` | Alert routing |
| AzureMonitorMetricAlert | throttle + capture-backlog | The two failure modes that matter |

## Parameters

| Parameter | Description | Default | Must change |
| --- | --- | --- | --- |
| `region` | Azure region | `centralus` | |
| `eventhub_namespace_name` | Globally unique namespace name | `my-event-streams` | yes |
| `archive_storage_account_name` | Globally unique archive account name | `myeventarchive` | yes |
| `auto_inflate_enabled` | Scale TUs with demand instead of throttling | `true` | |
| `max_throughput_units` | The auto-inflate (cost) ceiling, 1-40 | `10` | |
| `ingest_partition_count` | Ingest parallelism ceiling (fixed for life) | `8` | |
| `processed_partition_count` | Processed hub partitions | `4` | |
| `retention_days` | Stream replay window (STANDARD caps at 7) | `7` | |
| `capture_enabled` | The Data Lake archive arm | `true` | |
| `keyless_only_enabled` | Entra-only data plane; SAS off namespace-wide | `false` | |
| `ops_email` | Alert recipient | `ops@example.com` | yes |
| `log_retention_days` | Workspace retention | `30` | |

## After deploying

1. **Register schemas first** — producers register Avro schemas in the `event-schemas` group and serialize by schema id; consumers resolve ids from the registry. BACKWARD compatibility means consumers upgrade before producers.
2. **Producers send to `ingest`** — external producers use the `ingest-send` credential's connection string output; Azure-hosted producers should authenticate as managed identities with the *Azure Event Hubs Data Sender* role instead.
3. **Consume through your own group** — each application reads through its own consumer group (`stream-processor`, `dashboard`), never `$Default`; add a group per new consumer.
4. **Verify capture** — within minutes of the first events, Avro files appear under `event-archive/{Namespace}/{EventHub}/{PartitionId}/...` in the archive account.

## Day 2

- **New consumer**: add an `AzureEventHubConsumerGroup` on the hub it reads — offsets are isolated by group, nothing else changes.
- **Keyless everything**: once every producer/consumer runs with a managed identity, flip `keyless_only_enabled` — SAS shuts off namespace-wide and the chart's SAS credentials disappear.
- **Beyond 40 TUs**: sustained throttling at the ceiling means PREMIUM (processing units, 90-day retention, up to 1024 partitions) or a dedicated cluster.
- **Lakehouse growth**: the archive account is the raw zone; when curated zones and per-zone encryption enter the picture, the `data-lakehouse-storage` chart is the next step.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
