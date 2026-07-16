---
title: "Presets"
description: "Ready-to-deploy configuration presets for NAT Gateway"
type: "preset-list"
componentSlug: "nat-gateway"
componentTitle: "NAT Gateway"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard"
    rank: "01"
    title: "Standard NAT Gateway"
    excerpt: "This preset creates a zonal Standard NAT Gateway that SNATs through one referenced `AzurePublicIp`, giving every subnet that attaches it (via the subnet's `natGatewayId`) stable outbound connectivity..."
  - slug: "02-prefix-snat-range"
    rank: "02"
    title: "NAT Gateway with a Prefix SNAT Range"
    excerpt: "This preset SNATs through a referenced `AzurePublicIpPrefix` instead of individual addresses: one contiguous, pre-communicated CIDR that partners and firewalls allowlist once. A /28 prefix (16..."
  - slug: "03-zone-redundant-v2"
    rank: "03"
    title: "Zone-Redundant StandardV2 NAT Gateway"
    excerpt: "This preset creates Azure's next-generation StandardV2 NAT Gateway: zone-redundant automatically, with no zone pinning (`zones` must be left empty -- the spec enforces it). Where a Standard gateway..."
---

# NAT Gateway Presets

Ready-to-deploy configuration presets for NAT Gateway. Each preset is a complete manifest you can copy, customize, and deploy.
