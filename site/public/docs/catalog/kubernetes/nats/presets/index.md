---
title: "Presets"
description: "Ready-to-deploy configuration presets for NATS"
type: "preset-list"
componentSlug: "nats"
componentTitle: "NATS"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-dev"
    rank: "01"
    title: "Dev preset"
    excerpt: "The smallest useful NATS: one server, JetStream on. A single server is a complete JetStream deployment for dev — streams, consumers, KV and object stores all work, and the file-store volume means a..."
  - slug: "02-production"
    rank: "02"
    title: "Production preset"
    excerpt: "A 3-server NATS cluster with JetStream on 20Gi file stores per server, authenticated clients, and Prometheus metrics. Three servers is the smallest count that keeps replicated (R3) streams available..."
---

# NATS Presets

Ready-to-deploy configuration presets for NATS. Each preset is a complete manifest you can copy, customize, and deploy.
