---
title: "Presets"
description: "Ready-to-deploy configuration presets for Private DNS Zone"
type: "preset-list"
componentSlug: "private-dns-zone"
componentTitle: "Private DNS Zone"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "02-internal-zone"
    rank: "02"
    title: "Custom Internal DNS Zone"
    excerpt: "This preset creates a custom internal zone (e.g. `corp.internal`) for private name resolution, with the two SOA fields worth pinning: a real DNS-admin contact and a 30-second negative-caching TTL so..."
---

# Private DNS Zone Presets

Ready-to-deploy configuration presets for Private DNS Zone. Each preset is a complete manifest you can copy, customize, and deploy.
