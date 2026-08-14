---
title: "Spoke Link"
description: "This preset attaches one spoke network to the hub's forwarding ruleset -- the moment it deploys, resources in the spoke resolve the ruleset's domains through the hub resolver's outbound endpoint. No..."
type: "preset"
rank: "01"
presetSlug: "01-spoke-link"
componentSlug: "dns-resolver-virtual-network-link"
componentTitle: "DNS Resolver Virtual Network Link"
provider: "azure"
icon: "package"
order: 1
---

# Spoke Link

This preset attaches one spoke network to the hub's forwarding ruleset -- the moment it deploys, resources in the spoke resolve the ruleset's domains through the hub resolver's outbound endpoint. No resolver, endpoints, or peering needed in the spoke.

## When to Use

- Every spoke network in a hub-and-spoke topology that should resolve on-premises names
- The hub's OWN network too -- it is never linked implicitly, and forgetting it is the classic "works from spokes, fails from the hub" mystery

## Key Configuration Choices

- **Name the link after the network it attaches** -- the ruleset's link list becomes its audience roster
- **`metadata` records the network's owner** -- when the spoke is decommissioned, the link to clean up is self-identifying
- **The network must be in the ruleset's region** (cross-subscription is fine, cross-region is not)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-forwarding-ruleset>` | The AzurePrivateDnsResolverForwardingRuleset to turn on | The ruleset component's name |
| `<your-spoke-vnet>` | The AzureVirtualNetwork that starts forwarding | The network component's name |

Links are free at rest. One link per ruleset-network pair -- deploy one of these per network, owned alongside that network.
