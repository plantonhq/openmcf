# Overview

The **AzureNetworkWatcherFlowLog** component records network traffic metadata -- source, destination, port, protocol, and the allow/deny verdict -- for one target (a virtual network, subnet, or network interface) into a storage account. Optionally, Traffic Analytics pipes the records into a Log Analytics workspace for query, topology, and threat views. It is the network's flight recorder: when an audit or an incident asks "who talked to what," flow logs are the answer.

## Purpose

- **Auditability by default**: traffic records exist BEFORE the incident that needs them.
- **Scoped recording**: target the whole network, one subnet, or one interface -- the narrowest scope that answers the question.
- **From files to insight**: Traffic Analytics turns raw records into queryable flows, topology maps, and threat detections.

## Key Features

- Full azurerm v5 surface: virtual network/subnet/NIC targets, retention policy, schema version 1/2, and the complete Traffic Analytics block.
- The regional Network Watcher is referenced, never created -- unset watcher fields resolve to the singleton Azure auto-creates per region, so manifests say nothing about plumbing.
- NSG targets are rejected in seconds with the retirement explained (Azure stopped accepting new NSG flow logs in June 2025) instead of failing at deploy.
- Chart-ready: references the target, storage account, and workspace by typed outputs; publishes the flow log's ARM ID and effective watcher.

## Use Cases

- **Compliance evidence**: retained traffic records for the subnets carrying regulated workloads.
- **Incident forensics**: reconstruct who reached a compromised interface, and when.
- **Cost and capacity analysis**: Traffic Analytics surfaces top talkers and cross-region flows.

## Future Enhancements

- Typed references widen as more target kinds prove useful in charts.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
