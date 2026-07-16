---
title: "Presets"
description: "Ready-to-deploy configuration presets for LB Listener Rule"
type: "preset-list"
componentSlug: "lb-listener-rule"
componentTitle: "LB Listener Rule"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-path-based-routing"
    rank: "01"
    title: "Path-Based Routing"
    excerpt: "This preset routes one URL prefix to one service: requests whose path matches `/api/*` forward to the service's target group, while everything else falls through to lower-priority rules and the..."
  - slug: "02-host-based-routing"
    rank: "02"
    title: "Host-Based Routing"
    excerpt: "This preset routes one hostname to one service: requests whose Host header matches (say `api.example.com`) forward to the service's target group. Each service gets its own subdomain while sharing the..."
  - slug: "03-canary-weighted"
    rank: "03"
    title: "Canary Weighted Forward"
    excerpt: "This preset splits one route's traffic between two target groups -- 95% to the stable version, 5% to the canary -- with group stickiness so a client stays on whichever version served its first..."
---

# LB Listener Rule Presets

Ready-to-deploy configuration presets for LB Listener Rule. Each preset is a complete manifest you can copy, customize, and deploy.
