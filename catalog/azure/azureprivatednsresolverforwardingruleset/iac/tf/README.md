# AzurePrivateDnsResolverForwardingRuleset Terraform Module

## Overview

Creates a DNS forwarding ruleset -- the rule book that steers DNS queries for chosen domains out of Azure through a DNS Private Resolver's outbound endpoints -- and its forwarding rules as composed child resources. The ruleset takes effect in a network only once that network is linked to it (`AzurePrivateDnsResolverVirtualNetworkLink`).

## Resources Created

- `azurerm_private_dns_resolver_dns_forwarding_ruleset` -- the ruleset, bound to the resolver's outbound endpoint(s)
- `azurerm_private_dns_resolver_forwarding_rule` -- one per `spec.forwarding_rules` entry, keyed by rule name

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzurePrivateDnsResolverForwardingRulesetSpec fields; the resource group and outbound endpoint references arrive as resolved literals

## Outputs

- `dns_forwarding_ruleset_id` -- the ruleset's ARM resource ID (what virtual network links reference)
- `dns_forwarding_ruleset_name` -- the ruleset's name

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **A ruleset binds at most 2 outbound endpoints, both from the SAME resolver** -- Azure enforces it at deploy time (a service rule, deliberately not mirrored as validation because Microsoft adjusts limits over time).
- **Domain names are fully qualified WITH the trailing dot** (`corp.contoso.com.`) -- ARM stores them that way; write them that way in the spec.
- **Everything on a rule updates in place EXCEPT `domain_name`**, which replaces that rule (keyed by name, so siblings are untouched). The endpoint list and tags update the ruleset in place; name, resource group, and region replace it.
- **The port is always sent explicitly** -- 53 (the standard DNS port and ARM's default) when the spec leaves it unset -- so plans stay deterministic across engines.
- **Rules carry NO tags** -- ARM gives them a free-form `metadata` map instead, modeled per rule.
- **Rulesets and rules are free at rest** -- only endpoint hours and query volume bill. Azure caps a ruleset at 1,000 rules and a rule at 6 target servers (service quotas, documented not validated).

## Required Permissions

The deploying principal needs `Microsoft.Network/dnsForwardingRulesets/*` plus read on the resolver's outbound endpoints (Network Contributor on the resource group covers both).
