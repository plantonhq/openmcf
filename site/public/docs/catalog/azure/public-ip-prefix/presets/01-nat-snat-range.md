---
title: "NAT Gateway SNAT Range"
description: "This preset reserves a zone-redundant /28 prefix (16 contiguous public addresses) -- Azure's default production size for NAT gateway outbound SNAT. Each address contributes 64,512 SNAT ports, so the..."
type: "preset"
rank: "01"
presetSlug: "01-nat-snat-range"
componentSlug: "public-ip-prefix"
componentTitle: "Public IP Prefix"
provider: "azure"
icon: "package"
order: 1
---

# NAT Gateway SNAT Range

This preset reserves a zone-redundant /28 prefix (16 contiguous public
addresses) -- Azure's default production size for NAT gateway outbound SNAT.
Each address contributes 64,512 SNAT ports, so the /28 scales outbound
capacity sixteenfold in one allowlistable CIDR instead of managing sixteen
separate public IPs.

Wire the prefix to a NAT gateway via `publicIpPrefixIds`; subnets attach
the gateway on their own side (`AzureSubnet.natGatewayId`).

## When to Use

- Production subnets that need scalable, predictable outbound internet
  connectivity (replacing Azure's retired implicit default outbound access)
- Workloads at risk of SNAT port exhaustion (AKS, microservices, CI agents)
- Any design where egress IPs should stay in one contiguous, allowlistable
  range

## Key Configuration Choices

- **`prefixLength: 28`** -- 16 addresses; Azure's default and the usual
  starting point. Move to /27 or /26 only when port utilization demands it
- **Zone-redundant zones** -- `["1","2","3"]` survives a single
  availability-zone failure; pin to one zone only when the gateway itself is
  zonal and you want co-location
- **`ip_prefix` output** -- share the deployed CIDR with partners and
  firewalls once; it is assigned by Azure at creation and cannot be chosen
  in advance

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the prefix in | The resource group's `status.outputs.resource_group_name` |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

## Downstream Wiring

Reference this prefix from a NAT gateway:

```yaml
spec:
  publicIpPrefixIds:
    - valueFrom:
        name: my-nat-snat-prefix
```
