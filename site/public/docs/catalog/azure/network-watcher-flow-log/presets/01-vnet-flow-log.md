---
title: "Virtual Network Flow Log"
description: "This preset records every flow in one virtual network into a storage account with 30-day retention -- the audit baseline: traffic records exist before the incident that needs them."
type: "preset"
rank: "01"
presetSlug: "01-vnet-flow-log"
componentSlug: "network-watcher-flow-log"
componentTitle: "Network Watcher Flow Log"
provider: "azure"
icon: "package"
order: 1
---

# Virtual Network Flow Log

This preset records every flow in one virtual network into a storage account with 30-day retention -- the audit baseline: traffic records exist before the incident that needs them.

## When to Use

- The default recording posture for any production network
- Compliance regimes that require traffic records for regulated workloads
- The starting point for incident forensics ("who reached what, and when")

## Key Configuration Choices

- **The whole network is the target** -- comprehensive but voluminous; narrow to a subnet or interface when the question is narrower
- **`version: 2` explicitly** -- flow state and byte/packet counters; version 1 exists for legacy consumers
- **30-day retention** balances forensic reach against storage cost -- `retentionPolicy.days` is the single retention dial
- **No watcher plumbing** -- the region's auto-created Network Watcher is resolved automatically

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-virtual-network>` | The AzureVirtualNetwork whose traffic is recorded | The network component's name |
| `<your-storage-account>` | The AzureStorageAccount files land in (no hand-managed lifecycle rules) | The storage component's name |

Flow logs are near-free; storage is the cost that scales with traffic and retention.
