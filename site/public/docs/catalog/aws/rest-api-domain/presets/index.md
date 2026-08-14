---
title: "Presets"
description: "Ready-to-deploy configuration presets for REST API Domain"
type: "preset-list"
componentSlug: "rest-api-domain"
componentTitle: "REST API Domain"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-regional-mapped-domain"
    rank: "01"
    title: "Regional Mapped Domain"
    excerpt: "This preset fronts a REST API stage at `api.example.com` with a REGIONAL custom domain and a root base-path mapping."
  - slug: "02-edge-domain"
    rank: "02"
    title: "Edge Custom Domain"
    excerpt: "This preset fronts a REST API at a CloudFront-distributed hostname. The ACM certificate must live in `us-east-1` — that is CloudFront's region, regardless of where the API lives."
---

# REST API Domain Presets

Ready-to-deploy configuration presets for REST API Domain. Each preset is a complete manifest you can copy, customize, and deploy.
