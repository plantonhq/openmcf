---
title: "Presets"
description: "Ready-to-deploy configuration presets for DNS Zone"
type: "preset-list"
componentSlug: "dns-zone"
componentTitle: "DNS Zone"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-public-zone"
    rank: "01"
    title: "Public Zone"
    excerpt: "The 30-second public DNS zone: an empty, internet-facing zone for a domain you own, ready for records and registrar delegation. Azure's defaults cover the SOA record; governance tags carry your..."
  - slug: "02-delegation-ready-zone"
    rank: "02"
    title: "Delegation-Ready Zone"
    excerpt: "A public zone tuned for active operations: the SOA contact points at your DNS team and the negative-caching TTL is lowered so newly created records (certificate validation records especially) become..."
---

# DNS Zone Presets

Ready-to-deploy configuration presets for DNS Zone. Each preset is a complete manifest you can copy, customize, and deploy.
