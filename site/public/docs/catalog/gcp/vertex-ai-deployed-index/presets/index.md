---
title: "Presets"
description: "Ready-to-deploy configuration presets for Vertex AI Deployed Index"
type: "preset-list"
componentSlug: "vertex-ai-deployed-index"
componentTitle: "Vertex AI Deployed Index"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-automatic"
    rank: "01"
    title: "Automatic Serving Compute"
    excerpt: "Deploys an index onto an endpoint with Vertex-managed serving compute — the zero-configuration path. Vertex AI picks the machine types; replicas scale between the bounds you set."
  - slug: "02-dedicated"
    rank: "02"
    title: "Dedicated Serving Compute"
    excerpt: "Deploys an index with a pinned machine type and explicit replica bounds — predictable performance and cost for production serving."
  - slug: "03-peered-reserved-ranges"
    rank: "03"
    title: "Peered Endpoint with Reserved Ranges and JWT Auth"
    excerpt: "Deploys onto a VPC-peered endpoint with the deployment pinned to reserved IP ranges, access logging on, and JWT authentication on the private query endpoint — the full private-serving posture."
---

# Vertex AI Deployed Index Presets

Ready-to-deploy configuration presets for Vertex AI Deployed Index. Each preset is a complete manifest you can copy, customize, and deploy.
