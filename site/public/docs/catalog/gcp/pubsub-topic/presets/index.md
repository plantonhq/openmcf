---
title: "Presets"
description: "Ready-to-deploy configuration presets for Pub/Sub Topic"
type: "preset-list"
componentSlug: "pubsub-topic"
componentTitle: "Pub/Sub Topic"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-basic-topic"
    rank: "01"
    title: "Basic Topic"
    excerpt: "The minimal event channel: a named topic in the ambient project with Google-managed encryption and GCP defaults everywhere."
  - slug: "02-regional-encrypted"
    rank: "02"
    title: "Regional Encrypted Topic"
    excerpt: "The data-residency posture: customer-managed encryption plus hard region pinning for regulated event streams."
  - slug: "03-message-retention"
    rank: "03"
    title: "Topic with Message Retention"
    excerpt: "The replay-ready event stream: topic-level retention so consumers can seek backwards after a bug or a new subscriber needs history."
  - slug: "04-cloud-storage-ingestion"
    rank: "04"
    title: "Cloud Storage Ingestion Topic"
    excerpt: "The no-code file-to-stream bridge: Pub/Sub tails a GCS bucket and turns matching objects into messages — no Dataflow job, no custom loader."
  - slug: "05-schema-validated"
    rank: "05"
    title: "Schema-Validated Topic"
    excerpt: "The contract-enforced event stream: every published message is validated against a shared schema before it enters the topic."
---

# Pub/Sub Topic Presets

Ready-to-deploy configuration presets for Pub/Sub Topic. Each preset is a complete manifest you can copy, customize, and deploy.
