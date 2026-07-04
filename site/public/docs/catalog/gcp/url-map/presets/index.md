---
title: "Presets"
description: "Ready-to-deploy configuration presets for URL Map"
type: "preset-list"
componentSlug: "url-map"
componentTitle: "URL Map"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-host-path-fanout"
    rank: "01"
    title: "Host and Path Fan-Out"
    excerpt: "The classic global external Application Load Balancer routing table: one host rule sends `www.example.com` traffic into a path matcher that longest-prefix matches `/api/*` and `/static/*` to..."
  - slug: "02-weighted-canary"
    rank: "02"
    title: "Weighted Canary Split"
    excerpt: "Send a small share of production traffic to a canary backend service while the stable backend keeps the majority — the standard GCP mechanism for blue/green and canary rollouts at the URL map layer."
  - slug: "03-apex-redirect"
    rank: "03"
    title: "Apex-to-WWW HTTPS Redirect"
    excerpt: "A catch-all URL map whose only job is redirecting bare-apex (`example.com`) and unmatched traffic to `www.example.com` over HTTPS — the usual first step before attaching real backends behind a global..."
---

# URL Map Presets

Ready-to-deploy configuration presets for URL Map. Each preset is a complete manifest you can copy, customize, and deploy.
