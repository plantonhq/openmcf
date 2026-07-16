---
title: "Presets"
description: "Ready-to-deploy configuration presets for Firewall Policy Rule Collection Group"
type: "preset-list"
componentSlug: "firewall-policy-rule-collection-group"
componentTitle: "Firewall Policy Rule Collection Group"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-egress-baseline-rules"
    rank: "01"
    title: "Egress Baseline Rules"
    excerpt: "This preset creates the platform team's baseline rules group: DNS and NTP as network rules, and the package registries every build needs as application rules. It is the first group most policies..."
  - slug: "02-dnat-publish-service"
    rank: "02"
    title: "DNAT: Publish an Internal Service"
    excerpt: "This preset creates a DNAT group that publishes an internal service through the firewall's public IP -- inbound HTTPS arriving at the firewall is translated to an internal frontend (typically an..."
---

# Firewall Policy Rule Collection Group Presets

Ready-to-deploy configuration presets for Firewall Policy Rule Collection Group. Each preset is a complete manifest you can copy, customize, and deploy.
