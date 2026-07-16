---
title: "Presets"
description: "Ready-to-deploy configuration presets for HTTP API VPC Link"
type: "preset-list"
componentSlug: "http-api-vpc-link"
componentTitle: "HTTP API VPC Link"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-private-alb-link"
    rank: "01"
    title: "Preset: Private ALB Link"
    excerpt: "Use this preset to front an internal Application Load Balancer (or NLB / Cloud Map service) with HTTP APIs. One link per VPC is typically enough -- every API that needs private backends in that VPC..."
  - slug: "02-minimal-link"
    rank: "02"
    title: "Preset: Minimal Link"
    excerpt: "Use this preset for development or experimentation: a single-subnet, no-security-group link that gets a private integration working with the least ceremony. Reachability is then governed solely by..."
---

# HTTP API VPC Link Presets

Ready-to-deploy configuration presets for HTTP API VPC Link. Each preset is a complete manifest you can copy, customize, and deploy.
