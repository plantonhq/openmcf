---
title: "Presets"
description: "Ready-to-deploy configuration presets for ExpressRoute Port"
type: "preset-list"
componentSlug: "expressroute-port"
componentTitle: "ExpressRoute Port"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard-port"
    rank: "01"
    title: "Standard Port"
    excerpt: "This preset creates the common ExpressRoute Direct shape: a 10 Gbps Dot1Q port pair with both links administratively enabled and metered billing. The per-link outputs (router, interface, patch panel,..."
  - slug: "02-macsec-port"
    rank: "02"
    title: "MACsec Port"
    excerpt: "This preset creates a MACsec-encrypted ExpressRoute Direct port: layer-2 encryption on both physical links, keyed from your own Key Vault secrets through a user-assigned managed identity. Both spec..."
  - slug: "03-shared-capacity-port"
    rank: "03"
    title: "Shared Capacity Port"
    excerpt: "This preset creates the multi-tenant shape: a 100 Gbps QinQ port whose capacity is shared across subscriptions through issued authorizations. Each named authorization generates a key (surfaced,..."
---

# ExpressRoute Port Presets

Ready-to-deploy configuration presets for ExpressRoute Port. Each preset is a complete manifest you can copy, customize, and deploy.
