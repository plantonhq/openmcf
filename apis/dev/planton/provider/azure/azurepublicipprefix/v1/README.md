# AzurePublicIpPrefix

## Overview

`AzurePublicIpPrefix` provisions an Azure Public IP Prefix: a reserved,
contiguous range of public IP addresses. Individual public IPs allocated
from the prefix are guaranteed to come from the same known range -- which
lets partners and firewalls allowlist a single CIDR instead of chasing
individual addresses, and gives a NAT gateway a predictable, scalable SNAT
pool.

## Why a First-Class Resource?

A prefix is real infrastructure with its own lifecycle:

- **One range, many consumers** -- the same /28 serves a NAT gateway's
  outbound SNAT and any public IPs carved from it; the range is reserved
  once and referenced by ARM ID
- **Independent lifecycle** -- scale egress or add frontends by allocating
  from the prefix without replacing the gateway or re-allowlisting partners
- **Known allowlist surface** -- `ip_prefix` (e.g. `20.42.0.16/28`) is the
  CIDR partners and firewalls pin; it is assigned by Azure at creation and
  exported as an output

The prefix is referenced, never created inline, by `AzurePublicIp`
(`public_ip_prefix_id`) and `AzureNatGateway` (`public_ip_prefix_ids`).

## Key Features

- **Contiguous range sizing** -- `prefix_length` from /21 (2,048 addresses)
  down to /31 (2) for IPv4; Azure defaults to /28 (16 addresses)
- **Standard and StandardV2 SKUs** -- StandardV2 is required for StandardV2
  NAT gateways; GLOBAL tier requires Standard SKU
- **Zone-redundant anchoring** -- multiple zones (`"1"`, `"2"`, `"3"`) make
  the range resilient across availability-zone failures
- **BYOIP carve-out** -- optional `custom_ip_prefix_id` to allocate from a
  bring-your-own IP range already onboarded to Azure
- **Composable** -- the resource group is referenced by name, defaulting to
  an `AzureResourceGroup`'s output in composed environments

## Spec Fields

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `region` | string | Yes | -- | Azure region (must match consumers in the same region) |
| `resource_group` | StringValueOrRef | Yes | -- | Resource group name (defaults to an AzureResourceGroup reference) |
| `name` | string | Yes | -- | Prefix name, unique within the resource group (1-80 chars) |
| `prefix_length` | int32 | No | 28 | CIDR length of the range to reserve; fixed at creation |
| `ip_version` | enum | No | IPv4 | `IPV4` or `IPV6`; fixed at creation |
| `sku` | enum | No | STANDARD | `STANDARD` or `STANDARD_V2`; fixed at creation |
| `sku_tier` | enum | No | REGIONAL | `REGIONAL` or `GLOBAL` (cross-region LB only); fixed at creation |
| `zones` | list(string) | No | -- | Availability zones (`"1"`, `"2"`, `"3"`); fixed at creation |
| `custom_ip_prefix_id` | string | No | -- | ARM ID of a Custom IP Prefix for BYOIP; fixed at creation |
| `tags` | map | No | -- | User tags, merged over Planton-derived tags (user wins) |

## Outputs

| Output | Description |
|--------|-------------|
| `public_ip_prefix_id` | Full ARM ID of the prefix -- the join key `AzurePublicIp` and `AzureNatGateway` reference |
| `ip_prefix` | The actual reserved CIDR (e.g. `20.42.0.16/28`) -- the value partners allowlist |
| `public_ip_prefix_name` | The prefix's name as deployed |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzurePublicIpPrefix
metadata:
  name: prod-egress-prefix
  org: mycompany
  env: production
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: network-rg
  name: prod-egress-prefix
  prefixLength: 28
  zones:
    - "1"
    - "2"
    - "3"
```

Wire it to a NAT gateway:

```yaml
spec:
  publicIpPrefixIds:
    - valueFrom:
        name: prod-egress-prefix
```

## Lifecycle Notes

- Every field except `tags` is **fixed at creation**; changing name,
  region, length, SKU, tier, or zones **replaces the prefix** and assigns
  a new IP range -- treat that as a coordinated migration, not a casual edit
- A prefix **cannot be deleted** while any of its addresses are in use by
  public IPs or NAT gateway associations
- Bill for every address in the reserved range whether used or not; smaller
  `prefix_length` values reserve bigger (more expensive) ranges
