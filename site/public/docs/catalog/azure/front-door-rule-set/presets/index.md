---
title: "Presets"
description: "Ready-to-deploy configuration presets for Front Door Rule Set"
type: "preset-list"
componentSlug: "front-door-rule-set"
componentTitle: "Front Door Rule Set"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-security-headers"
    rank: "01"
    title: "Security Headers Policy"
    excerpt: "This preset creates a rule set with one condition-less rule that stamps the standard security headers on every response and strips the backend's technology fingerprint -- the baseline every..."
  - slug: "02-https-and-caching"
    rank: "02"
    title: "HTTPS Upgrade + Tiered Caching"
    excerpt: "This preset creates a three-rule delivery policy: permanent HTTPS upgrade at the edge, aggressive caching for immutable static assets (ignoring tracking parameters), and caching turned off for API..."
  - slug: "03-path-rewrite-and-canary"
    rank: "03"
    title: "Path Rewrite + Cookie Canary"
    excerpt: "This preset creates a delivery policy that decouples the public URL surface from the backend's real path layout (a rewrite the client never sees) and routes cookie-flagged requests to a canary origin..."
---

# Front Door Rule Set Presets

Ready-to-deploy configuration presets for Front Door Rule Set. Each preset is a complete manifest you can copy, customize, and deploy.
