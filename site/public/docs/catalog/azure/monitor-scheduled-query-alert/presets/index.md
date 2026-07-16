---
title: "Presets"
description: "Ready-to-deploy configuration presets for Monitor Scheduled Query Alert"
type: "preset-list"
componentSlug: "monitor-scheduled-query-alert"
componentTitle: "Monitor Scheduled Query Alert"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-error-spike"
    rank: "01"
    title: "Error Spike Alert (Row Count)"
    excerpt: "This preset creates the bread-and-butter log alert: count the error rows a KQL query returns and fire when the count crosses a line -- here, more than 10 application exceptions in 10 minutes."
  - slug: "02-latency-threshold"
    rank: "02"
    title: "Latency Threshold Alert (Metric Measurement, Per-Dimension)"
    excerpt: "This preset creates a metric-measurement log alert: the query computes average request duration per service role, and each role alerts independently when its average exceeds 500ms in 3 of the last 5..."
  - slug: "03-missing-heartbeat"
    rank: "03"
    title: "Missing Heartbeat Alert (Absence of Data)"
    excerpt: "This preset creates the absence-of-data alert: fire when heartbeat rows drop BELOW a floor -- the pattern that catches dead agents, stopped VMs, and silently broken ingestion pipelines, which no..."
---

# Monitor Scheduled Query Alert Presets

Ready-to-deploy configuration presets for Monitor Scheduled Query Alert. Each preset is a complete manifest you can copy, customize, and deploy.
