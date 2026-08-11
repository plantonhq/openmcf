---
title: "Presets"
description: "Ready-to-deploy configuration presets for Monitoring SLO"
type: "preset-list"
componentSlug: "monitoring-slo"
componentTitle: "Monitoring SLO"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-availability-slo"
    rank: "01"
    title: "Availability SLO"
    excerpt: "Three nines of good requests over a rolling 30 days, on a custom service created in the same apply — the standard first SLO for a service whose truth lives in its own metrics."
  - slug: "02-latency-slo"
    rank: "02"
    title: "Latency SLO"
    excerpt: "95% of requests under 500ms, measured over a calendar month — the \"fast enough\" objective, from the load balancer's own latency distribution."
---

# Monitoring SLO Presets

Ready-to-deploy configuration presets for Monitoring SLO. Each preset is a complete manifest you can copy, customize, and deploy.
