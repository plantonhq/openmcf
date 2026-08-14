---
title: "Presets"
description: "Ready-to-deploy configuration presets for Storage Bucket"
type: "preset-list"
componentSlug: "storage-bucket"
componentTitle: "Storage Bucket"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-private-standard"
    rank: "01"
    title: "Private Standard Bucket"
    excerpt: "The default posture for application data: IAM-only access control, public access impossible, versioned objects with bounded history, and an additive grant for the workload's service account."
  - slug: "02-static-website"
    rank: "02"
    title: "Static Website Bucket"
    excerpt: "A publicly readable multi-region bucket serving a static site: index and 404 pages configured, CORS opened for the application origin, and public access granted the additive way (one explicit..."
  - slug: "03-data-lake-autoclass"
    rank: "03"
    title: "Dual-Region Data Lake with Autoclass"
    excerpt: "Analytics storage that manages its own cost: a custom dual-region bucket where Autoclass moves every object to the cheapest justified storage class, hygiene rules reclaim hidden and temporary..."
  - slug: "04-event-driven-pipeline"
    rank: "04"
    title: "Event-Driven Pipeline Bucket"
    excerpt: "An intake bucket that announces its own changes: every new object under `uploads/` (and every delete) publishes an event to a Pub/Sub topic, which is where the processing pipeline — Cloud Run..."
---

# Storage Bucket Presets

Ready-to-deploy configuration presets for Storage Bucket. Each preset is a complete manifest you can copy, customize, and deploy.
