---
title: "Forced-Tunneling Firewall (Private Data Path)"
description: "This preset creates a firewall whose DATA path carries no public IP at all -- outbound traffic is forced on-premises (via ExpressRoute/VPN and a 0.0.0.0/0 route on the AzureFirewallSubnet) instead of..."
type: "preset"
rank: "02"
presetSlug: "02-forced-tunneling"
componentSlug: "firewall"
componentTitle: "Firewall"
provider: "azure"
icon: "package"
order: 2
---

# Forced-Tunneling Firewall (Private Data Path)

This preset creates a firewall whose DATA path carries no public IP at
all -- outbound traffic is forced on-premises (via ExpressRoute/VPN and a
0.0.0.0/0 route on the AzureFirewallSubnet) instead of egressing to the
internet from Azure. The separate management configuration gives Azure
the control-plane path it requires.

## When to Use

- Regulated environments where ALL egress must transit on-premises
  inspection before the internet
- Hybrid networks whose security boundary is the on-prem perimeter

## Key Configuration Choices

- **No `publicIpAddressId` on the data path** -- legal exactly because
  the management configuration exists (validated at authoring time)
- **`managementIpConfiguration`** -- needs its own subnet named exactly
  `AzureFirewallManagementSubnet` (/26+) and its own Standard static
  public IP; the block is FIXED AT CREATION (adding it later replaces
  the firewall), so decide on forced tunneling up front
- **BASIC tier note** -- Basic firewalls require this block regardless of
  tunneling

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group | Its `status.outputs.resource_group_name` |
| `<firewall-subnet-name>` | The AzureSubnet named exactly `AzureFirewallSubnet` | That subnet's Planton resource name |
| `<firewall-mgmt-subnet-name>` | The AzureSubnet named exactly `AzureFirewallManagementSubnet` | That subnet's Planton resource name |
| `<firewall-mgmt-pip-name>` | The management path's Standard static AzurePublicIp | That address's Planton resource name |
| `<firewall-policy-name>` | The AzureFirewallPolicy to enforce | That policy's Planton resource name |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |
