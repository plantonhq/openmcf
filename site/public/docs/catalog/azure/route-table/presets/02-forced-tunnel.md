---
title: "Forced Tunneling On-Premises"
description: "This preset sends all internet-bound traffic from attached subnets back on-premises through the virtual network gateway (VPN or ExpressRoute), where corporate egress controls inspect it -- the..."
type: "preset"
rank: "02"
presetSlug: "02-forced-tunnel"
componentSlug: "route-table"
componentTitle: "Route Table"
provider: "azure"
icon: "package"
order: 2
---

# Forced Tunneling On-Premises

This preset sends all internet-bound traffic from attached subnets back
on-premises through the virtual network gateway (VPN or ExpressRoute),
where corporate egress controls inspect it -- the classic
"forced tunneling" compliance pattern.

Note the trade-off: subnets attached to this table lose direct Azure
internet egress entirely. Azure service traffic that must stay direct
(and not hairpin on-premises) can be carved out with additional
service-tag routes (`addressPrefix: "AzureBackup"` →
`nextHopType: INTERNET`).

## When to Use

- Regulated environments where ALL egress must pass corporate inspection
- Networks connected on-premises via ExpressRoute or site-to-site VPN
  with a default-route advertisement policy

## Key Configuration Choices

- **`VIRTUAL_NETWORK_GATEWAY` next hop** -- no forwarding IP; the
  gateway is resolved from the network's gateway subnet
- **`bgpRoutePropagationEnabled: false`** -- keeps the user-defined
  default authoritative even when on-premises advertises its own
- **Service-tag carve-outs** -- add `INTERNET` routes per Azure service
  tag for traffic classes allowed to stay direct

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the table in | The resource group's `status.outputs.resource_group_name` |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |
