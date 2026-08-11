# AzurePrivateDnsResolverVirtualNetworkLink Terraform Module

## Overview

Creates the virtual network link that makes a DNS forwarding ruleset take effect in ONE virtual network. Once linked, resources in that network resolve the ruleset's domains through the resolver's outbound endpoints. One link per ruleset-network pair; spokes join and leave the topology without touching the ruleset or each other.

## Resources Created

- `azurerm_private_dns_resolver_virtual_network_link` -- the link (a child of the forwarding ruleset)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzurePrivateDnsResolverVirtualNetworkLinkSpec fields; the ruleset and virtual network references arrive as resolved literals

## Outputs

- `virtual_network_link_id` -- the link's ARM resource ID
- `virtual_network_link_name` -- the link's name

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **The linked network must be in the ruleset's region** but does NOT need to be peered with the resolver's network (hub-and-spoke: spokes link to the hub's ruleset), and cross-subscription links are allowed.
- **Everything except `metadata` is create-only** -- changing the name, ruleset, or network replaces the link (a brief forwarding gap for the affected network, nothing else).
- **The link carries NO tags** -- ARM gives it a free-form `metadata` map instead.
- **Links are free at rest**; Azure allows up to 500 per ruleset (a service quota, deliberately not mirrored as validation).

## Required Permissions

The deploying principal needs `Microsoft.Network/dnsForwardingRulesets/virtualNetworkLinks/*` plus join permissions on the linked network (Network Contributor on both resource groups covers it; cross-subscription links need it in each).
