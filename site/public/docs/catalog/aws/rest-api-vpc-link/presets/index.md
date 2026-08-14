---
title: "Presets"
description: "Ready-to-deploy configuration presets for REST API VPC Link"
type: "preset-list"
componentSlug: "rest-api-vpc-link"
componentTitle: "REST API VPC Link"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-nlb-vpc-link"
    rank: "01"
    title: "NLB VPC Link"
    excerpt: "This preset fronts an internal Network Load Balancer with a REST API VPC link — the only backend type API Gateway v1 links accept."
  - slug: "02-shared-backend-link"
    rank: "02"
    title: "Shared Backend Link"
    excerpt: "This preset is the production shape: one VPC link per NLB, shared by every REST API in the environment that needs that private backend."
---

# REST API VPC Link Presets

Ready-to-deploy configuration presets for REST API VPC Link. Each preset is a complete manifest you can copy, customize, and deploy.
