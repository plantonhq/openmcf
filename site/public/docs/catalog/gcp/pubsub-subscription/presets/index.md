---
title: "Presets"
description: "Ready-to-deploy configuration presets for Pub/Sub Subscription"
type: "preset-list"
componentSlug: "pubsub-subscription"
componentTitle: "Pub/Sub Subscription"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-basic-pull"
    rank: "01"
    title: "Basic Pull Subscription"
    excerpt: "The default consumer shape: a pull subscription with GCP defaults everywhere, composed onto a referenced topic."
  - slug: "02-push-with-oidc"
    rank: "02"
    title: "Push with OIDC Authentication"
    excerpt: "The authenticated webhook: Pub/Sub POSTs each message to an HTTPS endpoint with a signed OIDC JWT the receiver can verify."
  - slug: "03-bigquery-delivery"
    rank: "03"
    title: "BigQuery Delivery"
    excerpt: "The zero-ETL analytics sink: Pub/Sub streams every message straight into a BigQuery table — no Dataflow job, no custom consumer."
  - slug: "04-dead-letter"
    rank: "04"
    title: "Reliable Processing with Dead-Lettering"
    excerpt: "The production-money shape: exactly-once delivery, a deep replay window, backoff retries, and a dead-letter queue for poison messages."
---

# Pub/Sub Subscription Presets

Ready-to-deploy configuration presets for Pub/Sub Subscription. Each preset is a complete manifest you can copy, customize, and deploy.
