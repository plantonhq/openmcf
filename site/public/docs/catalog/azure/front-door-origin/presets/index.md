---
title: "Presets"
description: "Ready-to-deploy configuration presets for Front Door Origin"
type: "preset-list"
componentSlug: "front-door-origin"
componentTitle: "Front Door Origin"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-app-service-origin"
    rank: "01"
    title: "App Service Origin"
    excerpt: "This preset creates an origin pointing at an Azure App Service -- the most common Front Door backend. The defaults are exactly right for it: the Host header falls back to the origin hostname (which..."
  - slug: "02-weighted-canary"
    rank: "02"
    title: "Weighted Canary Origin"
    excerpt: "This preset adds a low-weight origin beside the main backend: same priority, ~5% of traffic. Ramping the canary is a weight change -- an in-place update -- and rolling back is deleting one resource."
  - slug: "03-private-link-origin"
    rank: "03"
    title: "Private Link Origin"
    excerpt: "This preset creates an origin that Front Door reaches over Azure Private Link -- the backend never sees the public internet and can disable public network access entirely."
---

# Front Door Origin Presets

Ready-to-deploy configuration presets for Front Door Origin. Each preset is a complete manifest you can copy, customize, and deploy.
