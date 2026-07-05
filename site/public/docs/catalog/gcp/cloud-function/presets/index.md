---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cloud Function"
type: "preset-list"
componentSlug: "cloud-function"
componentTitle: "Cloud Function"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-http-api"
    rank: "01"
    title: "HTTP API — basic"
    excerpt: "A public HTTPS endpoint for webhooks, small REST APIs, and glue code. The source lives as a versioned zip in GCS ([GcpGcsBucket](/docs/catalog/gcp/gcpgcsbucket)) — shipping a new deploy is uploading..."
  - slug: "02-pubsub-event"
    rank: "02"
    title: "Pub/Sub event processor"
    excerpt: "An event-driven worker consuming a Pub/Sub topic ([GcpPubSubTopic](/docs/catalog/gcp/gcppubsubtopic)) through Eventarc. Ingress is locked to internal traffic — nothing about an event consumer needs a..."
  - slug: "03-private-vpc-egress"
    rank: "03"
    title: "Private VPC egress — database-backed function"
    excerpt: "The composed private-networking pattern: the function routes egress through a [GcpServerlessVpcConnector](/docs/catalog/gcp/gcpserverlessvpcconnector) to reach a private-IP database (Cloud SQL..."
---

# Cloud Function Presets

Ready-to-deploy configuration presets for Cloud Function. Each preset is a complete manifest you can copy, customize, and deploy.
