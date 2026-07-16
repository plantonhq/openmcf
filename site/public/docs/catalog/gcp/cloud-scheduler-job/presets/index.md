---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cloud Scheduler Job"
type: "preset-list"
componentSlug: "cloud-scheduler-job"
componentTitle: "Cloud Scheduler Job"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-basic-http-job"
    rank: "01"
    title: "Basic HTTP Job"
    excerpt: "The minimal scheduled trigger: an unauthenticated HTTP GET fired on a cron cadence, with GCP-managed defaults for retries and deadlines."
  - slug: "02-pubsub-publisher"
    rank: "02"
    title: "Pub/Sub Publisher"
    excerpt: "Cron-driven event publishing: the job publishes a message to a Pub/Sub topic on schedule, and downstream subscriptions fan the event out to whatever needs to run."
  - slug: "03-secure-cloud-run-trigger"
    rank: "03"
    title: "Secure Cloud Run Trigger"
    excerpt: "The authenticated cron pattern: an HTTP POST to a private Cloud Run service with an OIDC token minted per attempt as a referenced service account — no credentials stored anywhere."
---

# Cloud Scheduler Job Presets

Ready-to-deploy configuration presets for Cloud Scheduler Job. Each preset is a complete manifest you can copy, customize, and deploy.
