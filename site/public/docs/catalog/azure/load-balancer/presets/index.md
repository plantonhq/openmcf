---
title: "Presets"
description: "Ready-to-deploy configuration presets for Load Balancer"
type: "preset-list"
componentSlug: "load-balancer"
componentTitle: "Load Balancer"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-public"
    rank: "01"
    title: "Public Load Balancer"
    excerpt: "This preset creates a public (internet-facing) Azure Load Balancer: one public frontend, a `web` backend pool, an HTTP health probe, and TCP rules for ports 80 and 443 with TCP reset enabled. This is..."
  - slug: "02-internal"
    rank: "02"
    title: "Internal Load Balancer"
    excerpt: "This preset creates an internal (private VNet) Azure Load Balancer: a zone-redundant frontend with a pinned static private address in your subnet, an `app` backend pool, a TCP health probe, and a TCP..."
  - slug: "03-outbound-and-nat"
    rank: "03"
    title: "Outbound SNAT + NAT Port Forwarding"
    excerpt: "This preset shows the full traffic story on one public load balancer: inbound load balancing on port 80, explicit outbound SNAT for the pool's egress, a single-target NAT rule forwarding port 2222 to..."
---

# Load Balancer Presets

Ready-to-deploy configuration presets for Load Balancer. Each preset is a complete manifest you can copy, customize, and deploy.
