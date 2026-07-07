---
title: "Presets"
description: "Ready-to-deploy configuration presets for Front Door Route"
type: "preset-list"
componentSlug: "front-door-route"
componentTitle: "Front Door Route"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-catch-all-https"
    rank: "01"
    title: "Catch-All HTTPS Route"
    excerpt: "This preset creates the endpoint's default route: every path, both protocols, HTTP redirected to HTTPS at the edge, no caching. The standard production entry rule for a dynamic application."
  - slug: "02-static-assets-cached"
    rank: "02"
    title: "Cached Static Assets Route"
    excerpt: "This preset creates a route dedicated to static content: matched paths are cached at every Front Door edge location and text-based types are compressed on the way out. The origin sees each asset..."
  - slug: "03-api-https-only"
    rank: "03"
    title: "HTTPS-Only API Route"
    excerpt: "This preset creates a strict API route: HTTPS is the only accepted protocol (plain HTTP fails instead of redirecting), the origin leg is pinned to HTTPS, and the public path prefix is rewritten on..."
---

# Front Door Route Presets

Ready-to-deploy configuration presets for Front Door Route. Each preset is a complete manifest you can copy, customize, and deploy.
