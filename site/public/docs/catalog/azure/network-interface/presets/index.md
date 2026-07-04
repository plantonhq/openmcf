---
title: "Presets"
description: "Ready-to-deploy configuration presets for Network Interface"
type: "preset-list"
componentSlug: "network-interface"
componentTitle: "Network Interface"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard"
    rank: "01"
    title: "Standard Private NIC"
    excerpt: "This preset creates a network interface with one dynamic IPv4 configuration in a referenced `AzureSubnet` and accelerated networking (SR-IOV) enabled. It is the shape virtually every VM starts from:..."
  - slug: "02-public-facing"
    rank: "02"
    title: "Public-Facing NIC with NIC-Level NSG"
    excerpt: "This preset creates a network interface whose primary configuration is fronted by a referenced `AzurePublicIp` and filtered by a NIC-level `AzureNetworkSecurityGroup`. It is the shape for a single..."
  - slug: "03-appliance-forwarding"
    rank: "03"
    title: "Network Virtual Appliance NIC (Forwarding + Static IP)"
    excerpt: "This preset creates the inside interface of a network virtual appliance -- a firewall, router, or SD-WAN box that routes other workloads' traffic. It pins a static private address (the next-hop IP..."
---

# Network Interface Presets

Ready-to-deploy configuration presets for Network Interface. Each preset is a complete manifest you can copy, customize, and deploy.
