---
title: "Presets"
description: "Ready-to-deploy configuration presets for Application Insights"
type: "preset-list"
componentSlug: "application-insights"
componentTitle: "Application Insights"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard"
    rank: "01"
    title: "Standard Web Application Insights"
    excerpt: "This preset creates a workspace-based Application Insights resource for a web application at full telemetry fidelity -- the default APM shape. Its `connection_string` output is what Function Apps,..."
  - slug: "02-production-sampled"
    rank: "02"
    title: "Production Application Insights (Sampled, Cost-Controlled)"
    excerpt: "This preset creates a production APM resource with the cost levers engaged: 50% sampling, a 10 GB daily cap with notification, one-year retention, and Entra-only ingestion."
---

# Application Insights Presets

Ready-to-deploy configuration presets for Application Insights. Each preset is a complete manifest you can copy, customize, and deploy.
