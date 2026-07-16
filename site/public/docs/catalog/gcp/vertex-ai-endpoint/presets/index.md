---
title: "Presets"
description: "Ready-to-deploy configuration presets for Vertex AI Endpoint"
type: "preset-list"
componentSlug: "vertex-ai-endpoint"
componentTitle: "Vertex AI Endpoint"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-basic-public"
    rank: "01"
    title: "Basic Public Endpoint"
    excerpt: "The minimal serving surface: a public Vertex AI endpoint in the ambient project with prediction logging sampled into BigQuery."
  - slug: "02-private-vpc-peered"
    rank: "02"
    title: "Private VPC-Peered Endpoint"
    excerpt: "Prediction serving reachable only from inside a peered VPC, encrypted under a customer-managed key, with an isolated dedicated DNS name."
  - slug: "03-private-psc"
    rank: "03"
    title: "Private PSC Endpoint"
    excerpt: "Prediction serving exposed through Private Service Connect: per-project access control and IAM-authorized connections, without VPC peering."
---

# Vertex AI Endpoint Presets

Ready-to-deploy configuration presets for Vertex AI Endpoint. Each preset is a complete manifest you can copy, customize, and deploy.
