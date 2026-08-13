---
title: "Isolated Spoke"
description: "This preset creates the classic isolation attachment: the spoke is routed by an isolated route table (so it never learns other isolated spokes' routes) and propagates its prefixes only to tables..."
type: "preset"
rank: "02"
presetSlug: "02-isolated-spoke"
componentSlug: "virtual-hub-connection"
componentTitle: "Virtual Hub Connection"
provider: "azure"
icon: "package"
order: 2
---

# Isolated Spoke

This preset creates the classic isolation attachment: the spoke is routed by an isolated route table (so it never learns other isolated spokes' routes) and propagates its prefixes only to tables labeled `shared` (so only shared-services networks can reach it). Pair it with the hub's **Isolated Spokes Hub** preset.

## When to Use

- Production spokes that must reach shared services but never each other
- Multi-tenant or multi-team hubs where reachability is an explicit grant

## Key Configuration Choices

- **Both halves matter** -- the association controls what this spoke can REACH; the propagation controls who can reach IT. Isolation requires getting both right
- **Reference tables by name** -- the hub surfaces its custom tables as `status.outputs.route_table_ids.<name>`; e.g. fieldPath `status.outputs.route_table_ids.isolated`
- **Labels scale the propagation** -- propagating to the label `shared` reaches every table carrying it, so shared-services tables added later need no connection changes

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-virtual-hub-arm-id>` | ARM ID of the hub | `AzureVirtualHub` status outputs (`virtual_hub_id`), or reference it with valueFrom |
| `<your-virtual-network-arm-id>` | ARM ID of the spoke VNet | `AzureVirtualNetwork` status outputs (`virtual_network_id`), or reference it with valueFrom |
| `<your-isolated-route-table-arm-id>` | ARM ID of the hub's isolated route table | `AzureVirtualHub` status outputs (`route_table_ids.<table-name>`), or reference it with valueFrom |
