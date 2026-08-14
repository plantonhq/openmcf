---
title: "Standard Hub"
description: "This preset creates the common case: a Standard-tier WAN hub with ARM's defaults stated -- ExpressRoute routing preference, branch-to-branch off, the router at its capacity floor of 2 units, and no..."
type: "preset"
rank: "01"
presetSlug: "01-standard-hub"
componentSlug: "virtual-hub"
componentTitle: "Virtual Hub"
provider: "azure"
icon: "package"
order: 1
---

# Standard Hub

This preset creates the common case: a Standard-tier WAN hub with ARM's defaults stated -- ExpressRoute routing preference, branch-to-branch off, the router at its capacity floor of 2 units, and no routing customization (any-to-any through the built-in default route table).

## When to Use

- The first hub of a Virtual WAN (one per region)
- Any topology that starts simple: attach spokes, add gateways, customize routing later

## Key Configuration Choices

- **A /23 address prefix** -- Microsoft's recommendation; the minimum is /24, and the value is fixed at creation
- **Standard tier** -- the full-mesh tier; Basic exists only for site-to-site-only Basic WANs
- **No routing customization** -- connections associate with and propagate to the default table (any-to-any); add route tables or routing intent when the topology demands it
- **Cost honesty** -- the hub bills (~$0.25/hr class) from creation and takes 15-30 minutes to provision

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | The hub's region (a WAN has at most one hub per region) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-virtual-wan-arm-id>` | ARM ID of the WAN this hub belongs to | `AzureVirtualWan` status outputs (`virtual_wan_id`), or reference it with valueFrom |
