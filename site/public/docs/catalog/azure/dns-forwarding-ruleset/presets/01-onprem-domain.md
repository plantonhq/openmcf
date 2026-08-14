---
title: "On-Premises Domain"
description: "This preset deploys the everyday shape: one ruleset forwarding the corporate Active Directory namespace to two datacenter DNS servers, bound to the hub resolver's outbound endpoint by reference."
type: "preset"
rank: "01"
presetSlug: "01-onprem-domain"
componentSlug: "dns-forwarding-ruleset"
componentTitle: "DNS Forwarding Ruleset"
provider: "azure"
icon: "package"
order: 1
---

# On-Premises Domain

This preset deploys the everyday shape: one ruleset forwarding the corporate Active Directory namespace to two datacenter DNS servers, bound to the hub resolver's outbound endpoint by reference.

## When to Use

- Azure workloads need to resolve on-premises names (domain controllers, internal services) over VPN or ExpressRoute
- The standard hub-and-spoke DNS pattern: rules written once on the hub's ruleset, consumed by every linked network

## Key Configuration Choices

- **Domains carry the trailing dot** (`corp.contoso.com.`) -- ARM stores them fully qualified; the rule captures the domain and everything under it
- **Targets are tried in order** -- put the primary datacenter's server first; both must be reachable FROM the outbound endpoint's subnet over your tunnel
- **The ruleset steers nothing until networks are linked** -- deploy AzurePrivateDnsResolverVirtualNetworkLink per network afterwards, including the hub's own

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The AzureResourceGroup the ruleset is created in | The resource-group component's name |
| `<your-dns-resolver>` | The AzurePrivateDnsResolver whose outbound endpoint the rules steer | The resolver component's name |

Replace the target `ipAddress` values with your actual on-premises DNS servers. Rulesets and rules are free at rest.
