---
title: "Public IP From Prefix"
description: "This preset shows the two-step composition: reserve a prefix first, then allocate an individual public IP from it via `publicIpPrefixId`. The address comes from the prefix's contiguous range and..."
type: "preset"
rank: "03"
presetSlug: "03-from-prefix-public-ip"
componentSlug: "public-ip-prefix"
componentTitle: "Public IP Prefix"
provider: "azure"
icon: "package"
order: 3
---

# Public IP From Prefix

This preset shows the two-step composition: reserve a prefix first, then
allocate an individual public IP from it via `publicIpPrefixId`. The
address comes from the prefix's contiguous range and remains covered by the
same partner allowlist CIDR (`status.outputs.ip_prefix` on the prefix).

The prefix and the public IP are separate resources with separate
lifecycles -- re-pointing a load balancer frontend to a new address carved
from the same prefix does not change what partners have allowlisted.

## When to Use

- Load balancer or application gateway frontends that must share an egress
  boundary with NAT or other workloads on the same prefix
- Reserving capacity once and drawing individual addresses as consumers
  appear
- Keeping partner allowlists stable while individual public IPs are
  replaced

## Key Configuration Choices

- **Prefix first** -- deploy `AzurePublicIpPrefix` and capture
  `public_ip_prefix_id` before creating consumers
- **`publicIpPrefixId` on the public IP** -- fixed at creation; omitting it
  allocates from Microsoft's general pool instead (a different allowlist
  surface)
- **Matching zones** -- align the public IP's zones with the prefix's for
  zone-redundant designs

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group for both resources | The resource group's `status.outputs.resource_group_name` |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

## Outputs to Use

| Resource | Output | Purpose |
| --- | --- | --- |
| `AzurePublicIpPrefix` | `ip_prefix` | The CIDR partners and firewalls allowlist |
| `AzurePublicIpPrefix` | `public_ip_prefix_id` | ARM ID wired into `publicIpPrefixId` |
| `AzurePublicIp` | `public_ip_id` | ARM ID for load balancers and gateways |
| `AzurePublicIp` | `ip_address` | The individual allocated address |
