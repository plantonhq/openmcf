---
title: "NAT Gateway with a Prefix SNAT Range"
description: "This preset SNATs through a referenced `AzurePublicIpPrefix` instead of individual addresses: one contiguous, pre-communicated CIDR that partners and firewalls allowlist once. A /28 prefix (16..."
type: "preset"
rank: "02"
presetSlug: "02-prefix-snat-range"
componentSlug: "nat-gateway"
componentTitle: "NAT Gateway"
provider: "azure"
icon: "package"
order: 2
---

# NAT Gateway with a Prefix SNAT Range

This preset SNATs through a referenced `AzurePublicIpPrefix` instead of individual addresses: one contiguous, pre-communicated CIDR that partners and firewalls allowlist once. A /28 prefix (16 addresses) provides 16 × 64,512 SNAT ports -- the scalable shape for high-connection-volume egress.

## When to Use

- High-throughput workloads that would exhaust a single address's 64,512 SNAT ports
- Egress whose source range partners must allowlist (one CIDR instead of N addresses)
- Architectures that reserve egress ranges up front and grow into them

## Key Configuration Choices

- **`publicIpPrefixIds` reference** -- resolves to the prefix's `public_ip_prefix_id` output; the prefix must be in the same region and zone as the gateway
- **Prefix length trades cost for ports** -- /31 (2 addresses) is the smallest reservation; /28 (16) suits heavy egress. Every reserved address bills whether or not traffic flows
- **Prefixes and individual addresses can mix** -- a gateway may carry both `publicIpIds` and `publicIpPrefixIds` simultaneously

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (must match the subnets it will serve) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-public-ip-prefix-resource-name>` | Planton metadata name of the `AzurePublicIpPrefix` (zone-matched) | Your prefix resource |
