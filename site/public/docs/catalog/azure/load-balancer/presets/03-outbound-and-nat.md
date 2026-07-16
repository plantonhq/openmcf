---
title: "Outbound SNAT + NAT Port Forwarding"
description: "This preset shows the full traffic story on one public load balancer: inbound load balancing on port 80, explicit outbound SNAT for the pool's egress, a single-target NAT rule forwarding port 2222 to..."
type: "preset"
rank: "03"
presetSlug: "03-outbound-and-nat"
componentSlug: "load-balancer"
componentTitle: "Load Balancer"
provider: "azure"
icon: "package"
order: 3
---

# Outbound SNAT + NAT Port Forwarding

This preset shows the full traffic story on one public load balancer: inbound load balancing on port 80, explicit outbound SNAT for the pool's egress, a single-target NAT rule forwarding port 2222 to one instance's SSH, and a pool-style NAT rule giving every pool member its own frontend SSH port.

## When to Use

- Production pools that make outbound connections: implicit SNAT has a small, exhaustion-prone port budget; explicit outbound rules size it deliberately
- Admin access to specific instances behind the LB without public IPs on the instances (single-target NAT + a NIC-side association)
- Per-instance access across a scale set (pool-style NAT: each member gets a dedicated frontend port from the range)

## Key Configuration Choices

- **`disableOutboundSnat: true` on the rule** -- required to combine a load-balancing rule and an outbound rule on the same pool; the outbound rule becomes the only SNAT path
- **`allocatedOutboundPorts: 2048`** -- sized explicitly (64,000 ports per frontend IP divided by maximum expected instances, in multiples of 8); `0` lets Azure divide the budget but reallocation churns connections on scale events
- **Single-target NAT (`ssh-admin`)** -- the LB declares the port forward; a NIC ip_configuration references `status.outputs.natRuleIds.ssh-admin` to pick the receiving instance
- **Pool-style NAT (`per-instance-ssh`)** -- the frontend range 50000-50099 maps one port to each pool member's port 22

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (must match backend resources) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-lb-name>` | Name for the load balancer (unique within resource group) | Your naming convention |
| `<public-ip-resource-id>` | Full ARM resource ID of a Standard SKU public IP | `AzurePublicIp` status outputs |

## Related Presets

- **01-public** -- inbound-only public load balancing
- **02-internal** -- private VNet load balancing
