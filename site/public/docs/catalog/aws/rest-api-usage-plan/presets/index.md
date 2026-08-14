---
title: "Presets"
description: "Ready-to-deploy configuration presets for REST API Usage Plan"
type: "preset-list"
componentSlug: "rest-api-usage-plan"
componentTitle: "REST API Usage Plan"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-metered-api-keys"
    rank: "01"
    title: "Metered API Keys"
    excerpt: "This preset covers one REST API stage with a 1000-request daily quota and one enabled API key — the starting shape for a metered consumer."
  - slug: "02-throttled-stages"
    rank: "02"
    title: "Throttled Stages"
    excerpt: "This preset covers one REST API stage with plan-wide throttle ceilings and a tighter cap on `GET /search` — the shape for protecting an expensive method without starving the rest of the API."
---

# REST API Usage Plan Presets

Ready-to-deploy configuration presets for REST API Usage Plan. Each preset is a complete manifest you can copy, customize, and deploy.
