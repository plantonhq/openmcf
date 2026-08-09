---
title: "Presets"
description: "Ready-to-deploy configuration presets for Logging Sink"
type: "preset-list"
componentSlug: "logging-sink"
componentTitle: "Logging Sink"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-error-archive-to-gcs"
    rank: "01"
    title: "Error Archive to GCS"
    excerpt: "The cheapest compliance archive: every ERROR-and-above entry in the project lands as hourly JSON batches in a Cloud Storage bucket."
  - slug: "02-audit-logs-to-bigquery"
    rank: "02"
    title: "Audit Logs to BigQuery"
    excerpt: "Audit-log forensics in SQL: every Cloud Audit Logs entry lands in partitioned BigQuery tables within seconds, queryable by principal, method, and resource."
  - slug: "03-log-stream-to-pubsub"
    rank: "03"
    title: "Log Stream to Pub/Sub"
    excerpt: "The front door to third-party log pipelines: entries stream to a Pub/Sub topic in near real time, where Datadog/Splunk-class collectors (or your own subscribers) consume them."
---

# Logging Sink Presets

Ready-to-deploy configuration presets for Logging Sink. Each preset is a complete manifest you can copy, customize, and deploy.
