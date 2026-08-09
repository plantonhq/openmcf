---
title: "Presets"
description: "Ready-to-deploy configuration presets for Monitoring Alert Policy"
type: "preset-list"
componentSlug: "monitoring-alert-policy"
componentTitle: "Monitoring Alert Policy"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-cpu-threshold"
    rank: "01"
    title: "CPU Threshold"
    excerpt: "The canonical infrastructure alert: any Compute Engine instance whose mean CPU stays above 80% for five sustained minutes opens a WARNING incident."
  - slug: "02-uptime-check-failure"
    rank: "02"
    title: "Uptime Check Failure"
    excerpt: "The other half of every availability monitor: pages CRITICAL when a GcpMonitoringUptimeCheck stops passing. This is Google's own recommended shape for uptime alerting — count the FALSE probe results..."
  - slug: "03-error-log-match"
    rank: "03"
    title: "Error Log Match"
    excerpt: "Pages on matching log entries — the application-truth tier of alerting, where a logged panic pages before any metric moves."
---

# Monitoring Alert Policy Presets

Ready-to-deploy configuration presets for Monitoring Alert Policy. Each preset is a complete manifest you can copy, customize, and deploy.
