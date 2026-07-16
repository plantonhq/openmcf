---
title: "Presets"
description: "Ready-to-deploy configuration presets for Public IP Prefix"
type: "preset-list"
componentSlug: "public-ip-prefix"
componentTitle: "Public IP Prefix"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-nat-snat-range"
    rank: "01"
    title: "NAT Gateway SNAT Range"
    excerpt: "This preset reserves a zone-redundant /28 prefix (16 contiguous public addresses) -- Azure's default production size for NAT gateway outbound SNAT. Each address contributes 64,512 SNAT ports, so the..."
  - slug: "02-partner-allowlist"
    rank: "02"
    title: "Partner Firewall Allowlist"
    excerpt: "This preset reserves a smaller /30 prefix (4 contiguous public addresses) optimized for partner and third-party firewall allowlisting. Partners pin one CIDR (`status.outputs.ip_prefix`) instead of..."
  - slug: "03-from-prefix-public-ip"
    rank: "03"
    title: "Public IP From Prefix"
    excerpt: "This preset shows the two-step composition: reserve a prefix first, then allocate an individual public IP from it via `publicIpPrefixId`. The address comes from the prefix's contiguous range and..."
---

# Public IP Prefix Presets

Ready-to-deploy configuration presets for Public IP Prefix. Each preset is a complete manifest you can copy, customize, and deploy.
