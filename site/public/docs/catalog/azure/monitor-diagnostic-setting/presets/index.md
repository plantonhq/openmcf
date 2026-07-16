---
title: "Presets"
description: "Ready-to-deploy configuration presets for Monitor Diagnostic Setting"
type: "preset-list"
componentSlug: "monitor-diagnostic-setting"
componentTitle: "Monitor Diagnostic Setting"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-logs-to-workspace"
    rank: "01"
    title: "All Logs and Metrics to a Workspace"
    excerpt: "This preset routes everything a resource emits -- every log category (via the `allLogs` group) and all metrics -- to a Log Analytics Workspace in the modern resource-specific table layout. It is the..."
  - slug: "02-archive-to-storage"
    rank: "02"
    title: "Audit Trail to Storage Archival"
    excerpt: "This preset routes a resource's audit-class logs (via the `audit` category group) to a storage account -- the pennies-per-GB destination for compliance retention measured in years, where the..."
  - slug: "03-stream-to-siem"
    rank: "03"
    title: "Security Stream to an External SIEM"
    excerpt: "This preset streams a resource's audit logs and metrics to an Event Hub -- the standard hand-off point for external SIEMs and streaming analytics pipelines that consume Azure telemetry outside Azure."
---

# Monitor Diagnostic Setting Presets

Ready-to-deploy configuration presets for Monitor Diagnostic Setting. Each preset is a complete manifest you can copy, customize, and deploy.
