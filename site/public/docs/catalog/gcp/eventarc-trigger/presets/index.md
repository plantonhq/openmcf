---
title: "Presets"
description: "Ready-to-deploy configuration presets for Eventarc Trigger"
type: "preset-list"
componentSlug: "eventarc-trigger"
componentTitle: "Eventarc Trigger"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-pubsub-to-cloud-run"
    rank: "01"
    title: "Pub/Sub to Cloud Run"
    excerpt: "The workhorse event route: a message published to a Pub/Sub topic invokes a Cloud Run service — with your own topic as transport so the topic remains shared infrastructure."
  - slug: "02-audit-log-to-workflow"
    rank: "02"
    title: "Audit Log to Workflow"
    excerpt: "React to GCP control-plane changes with an orchestrated response: every `storage.buckets.create` API call starts a Workflow execution that can notify, tag, or remediate."
---

# Eventarc Trigger Presets

Ready-to-deploy configuration presets for Eventarc Trigger. Each preset is a complete manifest you can copy, customize, and deploy.
