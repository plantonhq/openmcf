---
title: "Presets"
description: "Ready-to-deploy configuration presets for Storage Container"
type: "preset-list"
componentSlug: "storage-container"
componentTitle: "Storage Container"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-private-container"
    rank: "01"
    title: "Private Container"
    excerpt: "This preset creates a private blob container -- the default posture: every read requires authorization (Entra data-plane roles or the account's keys). The right starting point for everything that is..."
  - slug: "02-public-cdn-origin"
    rank: "02"
    title: "Public CDN-Origin Container"
    excerpt: "This preset creates a container serving objects anonymously by direct URL -- the public-website and CDN-origin pattern. Listing stays disabled: a client must know an object's URL to fetch it."
  - slug: "03-tenant-scoped-encryption"
    rank: "03"
    title: "Tenant-Scoped Encryption Container"
    excerpt: "This preset creates a private container pinned to an encryption scope: every blob written to it encrypts under the scope's key, and individual writes cannot override the scope -- per-tenant key..."
---

# Storage Container Presets

Ready-to-deploy configuration presets for Storage Container. Each preset is a complete manifest you can copy, customize, and deploy.
