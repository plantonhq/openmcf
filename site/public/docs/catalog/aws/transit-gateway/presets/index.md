---
title: "Presets"
description: "Ready-to-deploy configuration presets for Transit Gateway"
type: "preset-list"
componentSlug: "transit-gateway"
componentTitle: "Transit Gateway"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-full-mesh-hub"
    rank: "01"
    title: "Full-Mesh Hub"
    excerpt: "The zero-routing-configuration starting point: default association and propagation stay enabled, so every VPC attached to this gateway can reach every other one out of the box."
  - slug: "02-segmented-hub"
    rank: "02"
    title: "Segmented Hub"
    excerpt: "The isolation-first posture: both default-table dials off, so nothing routes until an `AwsTransitGatewayRouteTable` explicitly associates and propagates attachments. Production cannot see..."
  - slug: "03-hybrid-connectivity-hub"
    rank: "03"
    title: "Hybrid Connectivity Hub"
    excerpt: "A hub prepared for on-premises connectivity: a deliberately chosen Amazon-side ASN that will not collide with your data center's BGP, ECMP enabled so parallel VPN tunnels aggregate bandwidth, and..."
---

# Transit Gateway Presets

Ready-to-deploy configuration presets for Transit Gateway. Each preset is a complete manifest you can copy, customize, and deploy.
