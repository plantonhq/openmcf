---
title: "Presets"
description: "Ready-to-deploy configuration presets for Private Link Service"
type: "preset-list"
componentSlug: "private-link-service"
componentTitle: "Private Link Service"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-behind-load-balancer"
    rank: "01"
    title: "Behind a Load Balancer"
    excerpt: "This preset publishes a service running behind a Standard internal load balancer -- the classic Private Link shape. Consumers in other virtual networks (or tenants) connect through private endpoints..."
  - slug: "02-fixed-destination-ip"
    rank: "02"
    title: "Fixed Destination IP"
    excerpt: "This preset publishes a single-instance service at one fixed private IP -- Private Link without a load balancer in front. Consumer traffic is NATed straight to the destination address."
---

# Private Link Service Presets

Ready-to-deploy configuration presets for Private Link Service. Each preset is a complete manifest you can copy, customize, and deploy.
