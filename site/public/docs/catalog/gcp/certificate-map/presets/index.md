---
title: "Presets"
description: "Ready-to-deploy configuration presets for Certificate Map"
type: "preset-list"
componentSlug: "certificate-map"
componentTitle: "Certificate Map"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-hostname-routing"
    rank: "01"
    title: "Hostname Routing"
    excerpt: "Per-domain certificates selected by SNI at the load balancer, with a wildcard PRIMARY fallback — the multi-domain TLS edge as a routing table."
  - slug: "02-wildcard-fallback"
    rank: "02"
    title: "Wildcard Fallback"
    excerpt: "One wildcard certificate serving every subdomain through a single PRIMARY entry — the simplest certificate map that can work."
---

# Certificate Map Presets

Ready-to-deploy configuration presets for Certificate Map. Each preset is a complete manifest you can copy, customize, and deploy.
