---
title: "Presets"
description: "Ready-to-deploy configuration presets for Log Bucket"
type: "preset-list"
componentSlug: "log-bucket"
componentTitle: "Log Bucket"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-compliance-retention"
    rank: "01"
    title: "Compliance Retention"
    excerpt: "A 400-day audit bucket with a status-code index and a PREVENT destroy posture — the storage half of a compliance logging pipeline (pair it with a GcpLoggingSink routing audit entries in)."
  - slug: "02-analytics-and-views"
    rank: "02"
    title: "Analytics And Views"
    excerpt: "The queryable application-log bucket: Log Analytics on, a linked BigQuery dataset for the data team, and a stderr view for on-call — one bucket serving three audiences with different eyes."
  - slug: "03-adopt-default-bucket"
    rank: "03"
    title: "Adopt Default Bucket"
    excerpt: "Bring the project's built-in `_Default` bucket under declarative management and raise its retention from GCP's 30-day default — the smallest change that stops silent log expiry in its tracks."
---

# Log Bucket Presets

Ready-to-deploy configuration presets for Log Bucket. Each preset is a complete manifest you can copy, customize, and deploy.
