---
title: "Presets"
description: "Ready-to-deploy configuration presets for Front Door Endpoint"
type: "preset-list"
componentSlug: "front-door-endpoint"
componentTitle: "Front Door Endpoint"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-production-endpoint"
    rank: "01"
    title: "Production Endpoint"
    excerpt: "This preset creates an enabled endpoint on an existing Front Door profile -- the public entry hostname client traffic arrives at. Routes attach to it; custom-domain DNS records CNAME onto its..."
  - slug: "02-maintenance-disabled"
    rank: "02"
    title: "Pre-Provisioned (Disabled) Endpoint"
    excerpt: "This preset creates the endpoint dark: fully provisioned, hostname generated, but not accepting traffic. Flipping `enabled` is a fast in-place update -- the launch switch."
---

# Front Door Endpoint Presets

Ready-to-deploy configuration presets for Front Door Endpoint. Each preset is a complete manifest you can copy, customize, and deploy.
