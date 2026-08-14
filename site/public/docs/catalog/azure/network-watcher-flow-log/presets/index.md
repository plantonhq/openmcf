---
title: "Presets"
description: "Ready-to-deploy configuration presets for Network Watcher Flow Log"
type: "preset-list"
componentSlug: "network-watcher-flow-log"
componentTitle: "Network Watcher Flow Log"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-vnet-flow-log"
    rank: "01"
    title: "Virtual Network Flow Log"
    excerpt: "This preset records every flow in one virtual network into a storage account with 30-day retention -- the audit baseline: traffic records exist before the incident that needs them."
  - slug: "02-traffic-analytics"
    rank: "02"
    title: "Flow Log with Traffic Analytics"
    excerpt: "This preset records a virtual network's flows AND processes them into a Log Analytics workspace -- queryable flows, topology maps, and threat detections instead of raw files."
---

# Network Watcher Flow Log Presets

Ready-to-deploy configuration presets for Network Watcher Flow Log. Each preset is a complete manifest you can copy, customize, and deploy.
