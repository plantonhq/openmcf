---
title: "Forced-Tunnel Clients"
description: "This preset routes EVERYTHING from connected clients into the hub -- `internetSecurityEnabled: true` advertises the default route (0.0.0.0/0) into the tunnel, so internet-bound traffic can be..."
type: "preset"
rank: "02"
presetSlug: "02-forced-tunnel-clients"
componentSlug: "point-to-site-vpn-gateway"
componentTitle: "Point-to-Site VPN Gateway"
provider: "azure"
icon: "package"
order: 2
---

# Forced-Tunnel Clients

This preset routes EVERYTHING from connected clients into the hub -- `internetSecurityEnabled: true` advertises the default route (0.0.0.0/0) into the tunnel, so internet-bound traffic can be inspected at a hub firewall. Two scale units for the heavier per-client load forced tunneling brings.

## When to Use

- Compliance postures requiring inspected internet egress for remote users
- Hubs that already run a firewall with routing intent (the inspection point forced tunneling assumes)

## Key Configuration Choices

- **Forced tunneling** -- `internetSecurityEnabled: true` on the pool; pair it with a hub firewall via routing intent, or clients get uninspected hub egress instead of protection
- **Two scale units** -- forced tunneling multiplies per-client bandwidth through the gateway; `scaleUnit` updates in place if sizing changes
- **A separate pool** -- distinct client space (`172.16.202.0/24`) so network policy can tell secured clients apart

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group-name>` | Resource group for the gateway | `AzureResourceGroup` status outputs (`resource_group_name`), or reference it with valueFrom |
| `<your-virtual-hub-arm-id>` | ARM ID of the virtual hub the gateway deploys into | `AzureVirtualHub` status outputs (`virtual_hub_id`), or reference it with valueFrom |
| `<your-vpn-server-configuration-arm-id>` | ARM ID of the authentication policy | `AzureVpnServerConfiguration` status outputs (`vpn_server_configuration_id`), or reference it with valueFrom |

Replace the example pool (`172.16.202.0/24`) with client space your network plan reserves.
