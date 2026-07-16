---
title: "Presets"
description: "Ready-to-deploy configuration presets for Front Door Profile"
type: "preset-list"
componentSlug: "front-door-profile"
componentTitle: "Front Door Profile"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard-delivery"
    rank: "01"
    title: "Standard Delivery"
    excerpt: "This preset creates a Standard-tier Front Door profile -- the container for a global CDN deployment. Endpoints, origin groups, origins, and routes compose against it as their own resources, so the..."
  - slug: "02-premium-private-origins"
    rank: "02"
    title: "Premium with Private Origins"
    excerpt: "This preset creates a Premium-tier profile with a system-assigned managed identity -- the shape for locked-down architectures where backends disable public access entirely and Front Door reaches them..."
  - slug: "03-compliance-log-scrubbing"
    rank: "03"
    title: "Compliance Log Scrubbing"
    excerpt: "This preset creates a Standard profile with access-log scrubbing turned all the way up -- client IP addresses, request URIs, and query-string arguments are masked before Front Door writes its logs."
---

# Front Door Profile Presets

Ready-to-deploy configuration presets for Front Door Profile. Each preset is a complete manifest you can copy, customize, and deploy.
