---
title: "Presets"
description: "Ready-to-deploy configuration presets for Public IP"
type: "preset-list"
componentSlug: "public-ip"
componentTitle: "Public IP"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard-static"
    rank: "01"
    title: "Standard Static Public IP"
    excerpt: "This preset creates a zone-redundant Azure Public IP with the Standard SKU and static allocation. Standard with static allocation is the only supported configuration (Azure retired the Basic SKU in..."
  - slug: "02-dns-labeled-endpoint"
    rank: "02"
    title: "DNS-Labeled Endpoint"
    excerpt: "This preset creates a zone-redundant Standard public IP with an Azure-managed DNS name: `{label}.{region}.cloudapp.azure.com`. The scope-hashed label (`domainNameLabelScope`) lets the same label..."
  - slug: "03-allowlisted-from-prefix"
    rank: "03"
    title: "Allowlisted Address from a Prefix"
    excerpt: "This preset allocates the public IP from a reserved `AzurePublicIpPrefix` instead of Microsoft's general pool. Every address drawn from the prefix falls inside one contiguous, pre-communicated CIDR..."
---

# Public IP Presets

Ready-to-deploy configuration presets for Public IP. Each preset is a complete manifest you can copy, customize, and deploy.
