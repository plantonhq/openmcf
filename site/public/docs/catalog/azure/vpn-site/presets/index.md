---
title: "Presets"
description: "Ready-to-deploy configuration presets for VPN Site"
type: "preset-list"
componentSlug: "vpn-site"
componentTitle: "VPN Site"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-single-link-branch"
    rank: "01"
    title: "Single-Link Branch"
    excerpt: "This preset describes the classic branch: one ISP with a static public IP, and a static list of the prefixes reachable behind the branch. The site is free and deploys nothing at the branch -- it is..."
  - slug: "02-dual-link-bgp-branch"
    rank: "02"
    title: "Dual-Link BGP Branch"
    excerpt: "This preset describes a resilient branch: two ISPs, each a separate link with its own BGP speaker, and no static prefix list -- the branch advertises its routes. A connection then builds one tunnel..."
---

# VPN Site Presets

Ready-to-deploy configuration presets for VPN Site. Each preset is a complete manifest you can copy, customize, and deploy.
