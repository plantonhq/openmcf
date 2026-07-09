---
title: "Presets"
description: "Ready-to-deploy configuration presets for Vertex AI Index Endpoint"
type: "preset-list"
componentSlug: "vertex-ai-index-endpoint"
componentTitle: "Vertex AI Index Endpoint"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-public"
    rank: "01"
    title: "Public Index Endpoint"
    excerpt: "The simplest serving surface: a public Vector Search index endpoint in the ambient project."
  - slug: "02-vpc-peered"
    rank: "02"
    title: "VPC-Peered Index Endpoint"
    excerpt: "A private serving surface reachable only inside a peered VPC — vector search that never touches the public internet."
  - slug: "03-psc"
    rank: "03"
    title: "Private Service Connect Index Endpoint"
    excerpt: "The strongest network isolation for vector search: consumers connect through a PSC service attachment, no VPC peering required."
---

# Vertex AI Index Endpoint Presets

Ready-to-deploy configuration presets for Vertex AI Index Endpoint. Each preset is a complete manifest you can copy, customize, and deploy.
