---
title: "Presets"
description: "Ready-to-deploy configuration presets for Log Metric"
type: "preset-list"
componentSlug: "log-metric"
componentTitle: "Log Metric"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-error-counter"
    rank: "01"
    title: "Error Counter"
    excerpt: "Count a service's error log entries, labeled by status code — the canonical log-based metric and the input to every error-rate alert."
  - slug: "02-latency-distribution"
    rank: "02"
    title: "Latency Distribution"
    excerpt: "Percentile latency charts from access logs alone: extract each request's latency into a histogram — no instrumentation change, no tracing dependency."
---

# Log Metric Presets

Ready-to-deploy configuration presets for Log Metric. Each preset is a complete manifest you can copy, customize, and deploy.
