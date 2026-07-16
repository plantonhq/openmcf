---
title: "Presets"
description: "Ready-to-deploy configuration presets for Virtual Network"
type: "preset-list"
componentSlug: "virtual-network"
componentTitle: "Virtual Network"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard"
    rank: "01"
    title: "Standard Virtual Network"
    excerpt: "This preset creates a general-purpose virtual network with a single /16 address space -- room for dozens of /24 subnets -- and Azure's default DNS resolver, which is what private DNS zone resolution..."
  - slug: "02-hub-custom-dns"
    rank: "02"
    title: "Hub Network with Custom DNS"
    excerpt: "This preset shapes the network as a hub for hybrid or hub-and-spoke topologies: two address-space blocks (shared services grow fast in hubs), custom DNS servers for on-premises integration, and a..."
  - slug: "03-ddos-protected"
    rank: "03"
    title: "DDoS-Protected Edge Network"
    excerpt: "This preset hardens a network that fronts public IPs: an attached DDoS Protection Plan (always-on traffic monitoring with adaptive mitigation for every public IP in the network) and virtual network..."
---

# Virtual Network Presets

Ready-to-deploy configuration presets for Virtual Network. Each preset is a complete manifest you can copy, customize, and deploy.
