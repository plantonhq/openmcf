---
title: "Presets"
description: "Ready-to-deploy configuration presets for Monitor Metric Alert"
type: "preset-list"
componentSlug: "monitor-metric-alert"
componentTitle: "Monitor Metric Alert"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-static-threshold"
    rank: "01"
    title: "Static Threshold Alert"
    excerpt: "This preset creates the classic metric alert: a fixed threshold on a platform metric (storage availability below 99.9%, averaged over 15 minutes), evaluated every minute, paging the on-call action..."
  - slug: "02-dynamic-anomaly"
    rank: "02"
    title: "Dynamic Anomaly Alert"
    excerpt: "This preset creates a machine-learning metric alert: Azure learns the metric's normal band (including daily and weekly seasonality) and fires when the value deviates -- in either direction. No fixed..."
  - slug: "03-webtest-availability"
    rank: "03"
    title: "Web-Test Availability Alert"
    excerpt: "This preset pages when an Application Insights availability (web) test fails from multiple locations at once -- the outside-in signal that real users cannot reach the site, independent of any..."
---

# Monitor Metric Alert Presets

Ready-to-deploy configuration presets for Monitor Metric Alert. Each preset is a complete manifest you can copy, customize, and deploy.
