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

- **`VIRTUAL_APPLIANCE` + forwarding address by reference** -- the route
  references the AzureFirewall's `private_ip_address` output, so the
  table always follows the firewall (redeploy the firewall and the route
  updates with it). A literal IP (`value: "10.0.1.4"`) works too, for a
  third-party appliance or an address managed outside Planton
- **`bgpRoutePropagationEnabled: false`** -- the standard hardening;
  without it, an on-premises-advertised default route could bypass the
  firewall
- **Empty start is fine too** -- remove the route and attach the empty
  table first if the firewall arrives later; routes update in place

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the table in | The resource group's `status.outputs.resource_group_name` |
| `<firewall-resource-name>` | The AzureFirewall resource the route forwards to (its `private_ip_address` output is resolved automatically) | The firewall resource's name; for a non-firewall appliance, use a literal `value:` with its IP instead |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |
