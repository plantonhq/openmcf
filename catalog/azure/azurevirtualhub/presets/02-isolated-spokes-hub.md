# Isolated Spokes Hub

This preset creates a hub prepared for spoke isolation: two custom route tables ("isolated" and "shared-services", each labeled) that connections associate with and propagate to. The tables themselves are free and carry no routes -- the isolation topology is expressed by each connection's routing block.

## When to Use

- Production spokes that must reach shared services but never each other
- Multi-team hubs where reachability is an explicit grant, not the default

## Key Configuration Choices

- **The pattern lives in the connections**: an isolated spoke's connection associates with the `isolated` table and propagates only to `shared-services` (and/or `default`), so other isolated spokes never learn its routes
- **Labels are bulk propagation**: a connection propagating to the label `shared` reaches every table carrying it -- add tables later without touching existing connections
- **Table IDs surface by name**: reference `status.outputs.route_table_ids.isolated` from a connection's `associatedRouteTableId`

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | The hub's region | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-virtual-wan-arm-id>` | ARM ID of the WAN this hub belongs to | `AzureVirtualWan` status outputs (`virtual_wan_id`), or reference it with valueFrom |
