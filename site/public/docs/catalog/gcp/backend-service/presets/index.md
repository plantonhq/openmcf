---
title: "Presets"
description: "Ready-to-deploy configuration presets for Backend Service"
type: "preset-list"
componentSlug: "backend-service"
componentTitle: "Backend Service"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-external-web-backend"
    rank: "01"
    title: "External Web Backend"
    excerpt: "The workhorse backend service: an instance-group pool behind the global external Application Load Balancer, health-checked, balanced on CPU utilization, draining gracefully on rollouts, with sampled..."
  - slug: "02-cdn-cached-api"
    rank: "02"
    title: "CDN-Cached API"
    excerpt: "A read-heavy API served through Cloud CDN with the origin in control: `USE_ORIGIN_HEADERS` caches exactly what the application marks cacheable, tracking parameters are stripped from the cache key,..."
  - slug: "03-iap-protected-internal-tool"
    rank: "03"
    title: "IAP-Protected Internal Tool"
    excerpt: "Zero-trust access to an internal tool without a VPN: the backend service sits behind the public global load balancer, but Identity-Aware Proxy authenticates every request against Google identities..."
---

# Backend Service Presets

Ready-to-deploy configuration presets for Backend Service. Each preset is a complete manifest you can copy, customize, and deploy.
