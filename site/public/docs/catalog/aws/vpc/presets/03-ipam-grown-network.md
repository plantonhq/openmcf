---
title: "IPAM-Governed, Growing Network"
description: "This preset creates a VPC whose address space is governed by AWS IPAM rather than hand-picked CIDRs: the primary /16 is allocated from an IPAM pool, additional ranges (an IPAM-sized /20, a pinned..."
type: "preset"
rank: "03"
presetSlug: "03-ipam-grown-network"
componentSlug: "vpc"
componentTitle: "VPC"
provider: "aws"
icon: "package"
order: 3
---

# IPAM-Governed, Growing Network

This preset creates a VPC whose address space is governed by AWS IPAM rather than hand-picked CIDRs: the primary /16 is allocated from an IPAM pool, additional ranges (an IPAM-sized /20, a pinned block from the 100.64.0.0/10 shared space, and an IPAM-allocated IPv6 /56) attach as in-place associations, and VPC Encryption Control enforces encryption in transit with exclusions for the gateway paths that cannot encrypt.

## When to Use

- Organizations running AWS IPAM for non-overlapping, auditable address allocation across many VPCs and accounts
- Networks expected to outgrow their first range -- secondary CIDRs add address space without recreating the VPC or anything inside it
- Compliance postures that require encryption in transit inside the VPC (with a deliberate, reviewed exclusion list)

## Key Configuration Choices

- **IPAM-allocated primary** (`ipv4IpamPoolId` + `ipv4NetmaskLength`) — IPAM picks the block; changing the pool replaces the VPC, so govern from day one
- **Mixed secondary sources** (`secondaryIpv4Cidrs`) — each entry is its own association: IPAM-sized entries grow governed space; the pinned `100.64.0.0/16` carves the carrier-grade NAT shared range for pods or appliances that must not consume routable space
- **IPv6 by IPAM** (`secondaryIpv6Cidrs`) — a /56 sized from your IPv6 pool; entries also accept `assignGenerated: true` (Amazon-provided) or a BYOIP `ipv6Pool`
- **Encryption enforced, exclusions explicit** (`encryptionControl`) — `enforce` blocks unencrypted traffic; internet and NAT gateway paths are excluded here because their traffic terminates outside the VPC. Run `mode: monitor` first and review findings before enforcing
- **Pool ids are literal today** — no IPAM pool catalog kind exists yet; supply pool ids from your IPAM console

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aws-region>` | AWS region for the VPC | Your deployment region |
| `<ipam-pool-id>` | IPv4 IPAM pool (ipam-pool-...) | VPC console → IPAM → Pools |
| `<ipv6-ipam-pool-id>` | IPv6 IPAM pool | VPC console → IPAM → Pools |

## Related Presets

- `01-production-dual-stack` — hand-picked CIDR with Amazon-provided IPv6
- `02-development` — the minimal single-range network
