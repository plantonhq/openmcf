---
title: "Presets"
description: "Ready-to-deploy configuration presets for Event Hub Cluster"
type: "preset-list"
componentSlug: "event-hub-cluster"
componentTitle: "Event Hub Cluster"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-dedicated-cluster"
    rank: "01"
    title: "Dedicated Cluster"
    excerpt: "This preset provisions a single-tenant Event Hubs cluster at the entry size (1 capacity unit) -- the top of the capacity ladder, which namespaces join via their `dedicatedClusterId` reference."
---

# Event Hub Cluster Presets

Ready-to-deploy configuration presets for Event Hub Cluster. Each preset is a complete manifest you can copy, customize, and deploy.
