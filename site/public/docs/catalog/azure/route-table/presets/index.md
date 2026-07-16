---
title: "Presets"
description: "Ready-to-deploy configuration presets for Route Table"
type: "preset-list"
componentSlug: "route-table"
componentTitle: "Route Table"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-firewall-egress"
    rank: "01"
    title: "Firewall Egress"
    excerpt: "This preset implements the most common user-defined-routing pattern: every attached subnet's internet-bound traffic (`0.0.0.0/0`) is steered to a network virtual appliance -- an Azure Firewall or..."
  - slug: "02-forced-tunnel"
    rank: "02"
    title: "Forced Tunneling On-Premises"
    excerpt: "This preset sends all internet-bound traffic from attached subnets back on-premises through the virtual network gateway (VPN or ExpressRoute), where corporate egress controls inspect it -- the..."
  - slug: "03-blackhole-guardrails"
    rank: "03"
    title: "Black-Hole Guardrails"
    excerpt: "This preset drops traffic to a protected prefix at the routing layer (`nextHopType: NONE`) -- a coarse, cheap guardrail that stops lateral movement toward a sensitive tier (a data subnet, a..."
---

# Route Table Presets

Ready-to-deploy configuration presets for Route Table. Each preset is a complete manifest you can copy, customize, and deploy.
