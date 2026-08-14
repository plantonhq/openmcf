---
title: "Flow Log with Traffic Analytics"
description: "This preset records a virtual network's flows AND processes them into a Log Analytics workspace -- queryable flows, topology maps, and threat detections instead of raw files."
type: "preset"
rank: "02"
presetSlug: "02-traffic-analytics"
componentSlug: "network-watcher-flow-log"
componentTitle: "Network Watcher Flow Log"
provider: "azure"
icon: "package"
order: 2
---

# Flow Log with Traffic Analytics

This preset records a virtual network's flows AND processes them into a Log Analytics workspace -- queryable flows, topology maps, and threat detections instead of raw files.

## When to Use

- Networks where someone actually looks at traffic (operations, security, capacity planning)
- Turning "we have flow logs" into "we can answer questions about traffic"

## Key Configuration Choices

- **Two workspace references, one workspace** -- Traffic Analytics wants both the workspace GUID (`workspace_customer_id`) and its ARM id (`workspace_id`); both come from the same AzureLogAnalyticsWorkspace
- **`workspaceRegion` is the WORKSPACE's region**, which may differ from the flow log's
- **The 60-minute interval** suits audit and forensics; switch to 10 only when someone watches in near-real-time -- ingestion cost follows the cadence

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-virtual-network>` | The AzureVirtualNetwork whose traffic is recorded | The network component's name |
| `<your-storage-account>` | The AzureStorageAccount files land in (no hand-managed lifecycle rules) | The storage component's name |
| `<your-workspace>` | The AzureLogAnalyticsWorkspace Traffic Analytics writes into | The workspace component's name |

Workspace ingestion is the cost that scales here -- scope the target and pick the interval deliberately.
