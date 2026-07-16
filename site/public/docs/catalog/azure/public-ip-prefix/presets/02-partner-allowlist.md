---
title: "Partner Firewall Allowlist"
description: "This preset reserves a smaller /30 prefix (4 contiguous public addresses) optimized for partner and third-party firewall allowlisting. Partners pin one CIDR (`status.outputs.ip_prefix`) instead of..."
type: "preset"
rank: "02"
presetSlug: "02-partner-allowlist"
componentSlug: "public-ip-prefix"
componentTitle: "Public IP Prefix"
provider: "azure"
icon: "package"
order: 2
---

# Partner Firewall Allowlist

This preset reserves a smaller /30 prefix (4 contiguous public addresses)
optimized for partner and third-party firewall allowlisting. Partners pin
one CIDR (`status.outputs.ip_prefix`) instead of chasing individual
addresses as workloads scale, while the reservation cost stays minimal.

Use this when egress volume is modest but external systems require a stable,
documented source range -- SaaS integrations, payment gateways, legacy
on-premises firewalls.

## When to Use

- Third-party APIs or partners that require IP allowlisting
- Compliance environments where egress sources must be documented as a
  single CIDR
- Low-to-moderate outbound volume where a /28 would over-reserve (and
  over-bill) unused addresses

## Key Configuration Choices

- **`prefixLength: 30`** -- 4 addresses; enough headroom to allocate a few
  public IPs from the range while keeping the allowlist surface small
- **Zone-redundant zones** -- production resilience even for a small range
- **`ip_prefix` is the deliverable** -- hand partners the output CIDR after
  deployment, not the prefix name or ARM ID

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the prefix in | The resource group's `status.outputs.resource_group_name` |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |
