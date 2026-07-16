---
title: "Presets"
description: "Ready-to-deploy configuration presets for Router NAT"
type: "preset-list"
componentSlug: "router-nat"
componentTitle: "Router NAT"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-all-subnets-auto"
    rank: "01"
    title: "All-Subnets Auto-Allocated NAT"
    excerpt: "This preset creates a Cloud Router with a NAT gateway that covers all subnets in the region using automatically allocated external IPs. This is the simplest and most common Cloud NAT configuration,..."
  - slug: "02-static-ip-allowlisting"
    rank: "02"
    title: "Static-IP Allowlisted Egress"
    excerpt: "This preset creates a NAT gateway whose egress IPs are referenced `GcpAddress` reservations — stable addresses a partner, payment processor, or compliance regime can allowlist — scoped to the..."
  - slug: "03-private-nat"
    rank: "03"
    title: "Private NAT for Network Connectivity Center"
    excerpt: "This preset creates a PRIVATE-type NAT gateway that translates traffic between VPC networks attached as Network Connectivity Center spokes — the mechanism that lets two spokes with overlapping subnet..."
---

# Router NAT Presets

Ready-to-deploy configuration presets for Router NAT. Each preset is a complete manifest you can copy, customize, and deploy.
