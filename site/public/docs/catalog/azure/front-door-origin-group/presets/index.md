---
title: "Presets"
description: "Ready-to-deploy configuration presets for Front Door Origin Group"
type: "preset-list"
componentSlug: "front-door-origin-group"
componentTitle: "Front Door Origin Group"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-single-origin"
    rank: "01"
    title: "Single-Origin Group"
    excerpt: "This preset creates an origin group with no health probe -- the honest shape for a single-backend deployment, where probing adds origin load without buying failover (there is nowhere to fail over to)."
  - slug: "02-multi-region-probed"
    rank: "02"
    title: "Multi-Region Probed Backends"
    excerpt: "This preset creates the production multi-region shape: HTTPS health probes against a dedicated endpoint, a latency window wide enough to spread traffic across regions, and session affinity off for..."
  - slug: "03-stateful-web-app"
    rank: "03"
    title: "Stateful Web App Backends"
    excerpt: "This preset creates a group tuned for session-based web applications: sticky sessions on, gentle HEAD probes against the root path, and a long traffic-restore ramp so recovering origins warm up..."
---

# Front Door Origin Group Presets

Ready-to-deploy configuration presets for Front Door Origin Group. Each preset is a complete manifest you can copy, customize, and deploy.
