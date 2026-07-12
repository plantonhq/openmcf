---
title: "Presets"
description: "Ready-to-deploy configuration presets for Monitor Activity Log Alert"
type: "preset-list"
componentSlug: "monitor-activity-log-alert"
componentTitle: "Monitor Activity Log Alert"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-service-health"
    rank: "01"
    title: "Service-Health Incident Alert"
    excerpt: "This preset creates a subscription-scoped alert that fires when Azure reports a service-health incident or maintenance event. Service-health alerts are the ones you cannot get any other way -- they..."
  - slug: "02-resource-delete"
    rank: "02"
    title: "Administrative Change Alert"
    excerpt: "This preset creates a resource-group-scoped alert on administrative operations at critical/error severity that succeeded -- a broad safety net for \"something significant just changed in this resource..."
---

# Monitor Activity Log Alert Presets

Ready-to-deploy configuration presets for Monitor Activity Log Alert. Each preset is a complete manifest you can copy, customize, and deploy.
