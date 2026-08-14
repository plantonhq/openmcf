---
title: "Presets"
description: "Ready-to-deploy configuration presets for VPN Gateway"
type: "preset-list"
componentSlug: "vpn-gateway"
componentTitle: "VPN Gateway"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard-hub-gateway"
    rank: "01"
    title: "Standard Hub Gateway"
    excerpt: "This preset creates the hub's branch on-ramp with sensible defaults: one scale unit (500 Mbps aggregate), Microsoft-backbone routing preference, and Azure's default BGP settings (ASN 65515). **The..."
  - slug: "02-nat-overlapping-branches"
    rank: "02"
    title: "NAT for Overlapping Branches"
    excerpt: "This preset creates a gateway prepared for the classic acquisition problem: two branches that both use `192.168.10.0/24`. Each gets a static ingress NAT rule translating it to a distinct..."
---

# VPN Gateway Presets

Ready-to-deploy configuration presets for VPN Gateway. Each preset is a complete manifest you can copy, customize, and deploy.
