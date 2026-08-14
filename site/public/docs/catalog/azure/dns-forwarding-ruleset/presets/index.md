---
title: "Presets"
description: "Ready-to-deploy configuration presets for DNS Forwarding Ruleset"
type: "preset-list"
componentSlug: "dns-forwarding-ruleset"
componentTitle: "DNS Forwarding Ruleset"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-onprem-domain"
    rank: "01"
    title: "On-Premises Domain"
    excerpt: "This preset deploys the everyday shape: one ruleset forwarding the corporate Active Directory namespace to two datacenter DNS servers, bound to the hub resolver's outbound endpoint by reference."
  - slug: "02-staged-multi-domain"
    rank: "02"
    title: "Staged Multi-Domain"
    excerpt: "This preset carries two namespaces on one rule book: the live corporate domain, plus an acquired company's domain PARKED with `enabled: false` until its connectivity is ready -- flipping one field..."
---

# DNS Forwarding Ruleset Presets

Ready-to-deploy configuration presets for DNS Forwarding Ruleset. Each preset is a complete manifest you can copy, customize, and deploy.
