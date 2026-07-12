---
title: "Presets"
description: "Ready-to-deploy configuration presets for Event Hub Namespace"
type: "preset-list"
componentSlug: "event-hub-namespace"
componentTitle: "Event Hub Namespace"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard-streaming"
    rank: "01"
    title: "Standard Streaming Namespace"
    excerpt: "This preset creates a STANDARD-tier Event Hubs namespace with elastic throughput -- the full-featured multi-tenant tier (Kafka endpoint, 20 consumer groups per hub, 7-day retention) that fits most..."
  - slug: "02-locked-down-keyless"
    rank: "02"
    title: "Locked-Down Keyless Namespace"
    excerpt: "This preset creates a STANDARD namespace in the production security posture: SAS authentication disabled (Entra-only data plane) and a DENY firewall admitting only named sources."
  - slug: "03-premium-isolated"
    rank: "03"
    title: "Premium Isolated Namespace"
    excerpt: "This preset creates a PREMIUM-tier namespace: reserved processing units with predictable latency, dynamic partition scale-up, and extended retention -- the isolation tier below a whole dedicated..."
---

# Event Hub Namespace Presets

Ready-to-deploy configuration presets for Event Hub Namespace. Each preset is a complete manifest you can copy, customize, and deploy.
