---
title: "Presets"
description: "Ready-to-deploy configuration presets for Listener Set"
type: "preset-list"
componentSlug: "listener-set"
componentTitle: "Listener Set"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-team-https-listeners"
    rank: "01"
    title: "Team HTTPS Listeners"
    excerpt: "Attach a team-owned HTTPS listener to a centrally managed Gateway without editing the Gateway itself. The platform team owns the Gateway (address, load balancer, GatewayClass); each application team..."
  - slug: "02-tls-passthrough-listener"
    rank: "02"
    title: "TLS Passthrough Listener"
    excerpt: "Attach a TLS passthrough listener to a shared Gateway through a ListenerSet. The listener forwards encrypted connections untouched -- the backend terminates TLS itself -- so the team exposes an..."
---

# Listener Set Presets

Ready-to-deploy configuration presets for Listener Set. Each preset is a complete manifest you can copy, customize, and deploy.
