---
title: "Presets"
description: "Ready-to-deploy configuration presets for DNS Private Resolver"
type: "preset-list"
componentSlug: "dns-private-resolver"
componentTitle: "DNS Private Resolver"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-hybrid-resolver"
    rank: "01"
    title: "Hybrid Resolver"
    excerpt: "This preset deploys the full hybrid shape: one inbound endpoint (on-premises resolves Azure names through it) and one outbound endpoint (Azure resolves on-premises names through it, steered by a..."
  - slug: "02-inbound-only"
    rank: "02"
    title: "Inbound Only (Pinned IP)"
    excerpt: "This preset deploys the one-way shape: on-premises (or another cloud) resolves Azure's private names by forwarding to the inbound endpoint, whose address is STATICALLY pinned so forwarder..."
---

# DNS Private Resolver Presets

Ready-to-deploy configuration presets for DNS Private Resolver. Each preset is a complete manifest you can copy, customize, and deploy.
