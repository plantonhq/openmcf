---
title: "Full-Mesh Hub"
description: "The zero-routing-configuration starting point: default association and propagation stay enabled, so every VPC attached to this gateway can reach every other one out of the box."
type: "preset"
rank: "01"
presetSlug: "01-full-mesh-hub"
componentSlug: "transit-gateway"
componentTitle: "Transit Gateway"
provider: "aws"
icon: "package"
order: 1
---

# Full-Mesh Hub

The zero-routing-configuration starting point: default association and propagation stay enabled, so every VPC attached to this gateway can reach every other one out of the box.

## When to Use

- A handful of trusted VPCs (platform, data, tooling) that should all communicate
- Replacing a VPC peering mesh without changing reachability semantics

## What It Configures

- **Nothing beyond the hub** — the gateway's defaults (auto-association, auto-propagation, DNS support, ECMP) already produce the full mesh
- Attach VPCs with `AwsTransitGatewayVpcAttachment` resources; no route tables needed

## What to Customize

- Replace `<aws-region>` with your region
- Set `amazonSideAsn` if the default 64512 collides with an on-premises ASN you plan to connect via VPN or Direct Connect
