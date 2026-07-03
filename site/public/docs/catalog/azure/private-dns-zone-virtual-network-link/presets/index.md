---
title: "Presets"
description: "Ready-to-deploy configuration presets for Private DNS Zone Virtual Network Link"
type: "preset-list"
componentSlug: "private-dns-zone-virtual-network-link"
componentTitle: "Private DNS Zone Virtual Network Link"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-privatelink-zone-link"
    rank: "01"
    title: "Private Link Zone Attachment"
    excerpt: "This preset links a Private Link zone (e.g. `privatelink.postgres.database.azure.com`) to a virtual network so workloads inside it resolve the PaaS service's FQDN to its private endpoint IP instead..."
  - slug: "02-internal-zone-autoregistration"
    rank: "02"
    title: "Internal Zone with VM Auto-Registration"
    excerpt: "This preset links a custom internal zone (e.g. `corp.internal`) to a virtual network with VM auto-registration on: every virtual machine in the network gets an A record at boot and loses it at..."
  - slug: "03-public-fallback"
    rank: "03"
    title: "Shared Zone with Public DNS Fallback"
    excerpt: "This preset links a zone with `NX_DOMAIN_REDIRECT` resolution: names the private zone cannot answer are retried against public DNS instead of failing. The fallback pattern for privatelink zones..."
---

# Private DNS Zone Virtual Network Link Presets

Ready-to-deploy configuration presets for Private DNS Zone Virtual Network Link. Each preset is a complete manifest you can copy, customize, and deploy.
