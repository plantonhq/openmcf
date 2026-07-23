---
title: "Presets"
description: "Ready-to-deploy configuration presets for UDP Route"
type: "preset-list"
componentSlug: "udp-route"
componentTitle: "UDP Route"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-dns-forwarding"
    rank: "01"
    title: "UDP DNS Forwarding"
    excerpt: "The most common UDPRoute: forward DNS queries arriving on a Gateway's UDP listener to a DNS Service on port 53, split across a weighted pair of backends. A UDP route has no matching at all -- the..."
  - slug: "02-game-server"
    rank: "02"
    title: "UDP Game Server"
    excerpt: "Forward all datagrams arriving on a Gateway's UDP listener to a single game server Service on a custom port. This is the pattern for any single-backend datagram protocol -- game servers, syslog..."
---

# UDP Route Presets

Ready-to-deploy configuration presets for UDP Route. Each preset is a complete manifest you can copy, customize, and deploy.
