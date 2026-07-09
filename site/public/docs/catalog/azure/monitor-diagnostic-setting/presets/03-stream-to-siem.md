---
title: "Security Stream to an External SIEM"
description: "This preset streams a resource's audit logs and metrics to an Event Hub -- the standard hand-off point for external SIEMs and streaming analytics pipelines that consume Azure telemetry outside Azure."
type: "preset"
rank: "03"
presetSlug: "03-stream-to-siem"
componentSlug: "monitor-diagnostic-setting"
componentTitle: "Monitor Diagnostic Setting"
provider: "azure"
icon: "package"
order: 3
---

# Security Stream to an External SIEM

This preset streams a resource's audit logs and metrics to an Event Hub -- the standard hand-off point for external SIEMs and streaming analytics pipelines that consume Azure telemetry outside Azure.

## When to Use

- Estates whose security operations run on an external SIEM (Splunk, QRadar, Elastic, Sentinel-in-another-tenant)
- Real-time processing pipelines that must react to events faster than workspace ingestion latency

## Key Configuration Choices

- **Namespace-level authorization rule** -- the `eventhubAuthorizationRuleId` is a NAMESPACE rule (e.g. RootManageSharedAccessKey), not a hub-level one; `eventhubName` then picks the hub
- **Named hub** (`eventhubName: siem-ingest`) -- omit it to let Azure route to a default hub per category; name it when the SIEM consumes one known hub
- **Audit group + metrics** -- the typical SIEM selection; widen to `allLogs` when the SIEM ingests everything

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-app-vault` | The resource whose telemetry is streamed | Any kind's status outputs (`*_id`) |
| `<subscription-id>`, `<rg>`, `<namespace>` | The Event Hub namespace's authorization rule path | The Event Hub namespace in the portal, or `AzureEventHubNamespace` outputs |

## Related Presets

- **01-logs-to-workspace** -- Keep an in-Azure queryable copy alongside the stream
- **02-archive-to-storage** -- Cheap compliance archival
