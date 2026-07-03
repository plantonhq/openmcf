---
title: "Firewall Egress"
description: "This preset implements the most common user-defined-routing pattern: every attached subnet's internet-bound traffic (`0.0.0.0/0`) is steered to a network virtual appliance -- an Azure Firewall or..."
type: "preset"
rank: "01"
presetSlug: "01-firewall-egress"
componentSlug: "route-table"
componentTitle: "Route Table"
provider: "azure"
icon: "package"
order: 1
---

# Firewall Egress

This preset implements the most common user-defined-routing pattern:
every attached subnet's internet-bound traffic (`0.0.0.0/0`) is steered
to a network virtual appliance -- an Azure Firewall or third-party
firewall -- instead of going straight out. BGP route propagation is
disabled so routes learned from on-premises can never bypass the
firewall.

Attach this table to every workload subnet that must egress through
inspection. The firewall's own subnet must NOT attach it (that would
loop its egress back to itself).

## When to Use

- Hub-and-spoke topologies where the hub firewall inspects all egress
- Compliance environments requiring egress filtering and logging
- Any subnet whose direct internet path should be closed

## Key Configuration Choices

- **`VIRTUAL_APPLIANCE` + forwarding IP** -- the firewall's private IP;
  update it in place if the appliance is redeployed
- **`bgpRoutePropagationEnabled: false`** -- the standard hardening;
  without it, an on-premises-advertised default route could bypass the
  firewall
- **Empty start is fine too** -- remove the route and attach the empty
  table first if the firewall arrives later; routes update in place

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the table in | The resource group's `status.outputs.resource_group_name` |
| `<firewall-private-ip>` | The firewall's/appliance's private IP inside the network | The firewall resource's IP configuration |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |
