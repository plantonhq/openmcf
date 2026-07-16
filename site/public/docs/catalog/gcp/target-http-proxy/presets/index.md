---
title: "Presets"
description: "Ready-to-deploy configuration presets for Target HTTP Proxy"
type: "preset-list"
componentSlug: "target-http-proxy"
componentTitle: "Target HTTP Proxy"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-https-redirect-frontend"
    rank: "01"
    title: "HTTPS Redirect Frontend"
    excerpt: "The production-standard role for a target HTTP proxy: port 80 exists only to bounce clients to HTTPS. The proxy points at a redirect-only URL map (its `defaultUrlRedirect` sets `httpsRedirect: true`..."
  - slug: "02-plain-http-frontend"
    rank: "02"
    title: "Plain HTTP Frontend"
    excerpt: "A target HTTP proxy that serves the application itself over plain HTTP — the right shape for internal test environments, health-probe endpoints, or services that terminate TLS elsewhere."
  - slug: "03-traffic-director-mesh"
    rank: "03"
    title: "Traffic Director Mesh Frontend"
    excerpt: "A target HTTP proxy for Traffic Director (service mesh): `proxyBind: true` binds the proxy to the mesh's private IPs instead of Google's edge, and the forwarding rule that references it uses the..."
---

# Target HTTP Proxy Presets

Ready-to-deploy configuration presets for Target HTTP Proxy. Each preset is a complete manifest you can copy, customize, and deploy.
