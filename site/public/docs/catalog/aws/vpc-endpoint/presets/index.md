---
title: "Presets"
description: "Ready-to-deploy configuration presets for VPC Endpoint"
type: "preset-list"
componentSlug: "vpc-endpoint"
componentTitle: "VPC Endpoint"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-s3-gateway"
    rank: "01"
    title: "S3 Gateway Endpoint"
    excerpt: "This preset gives a VPC's private subnets a free, private path to S3 by injecting the S3 prefix-list route into their route tables. It is the default cost-and-security move for any VPC that touches..."
  - slug: "02-interface-endpoint"
    rank: "02"
    title: "Interface Endpoint for an AWS Service"
    excerpt: "This preset places an ENI-based PrivateLink endpoint for an AWS service in two private subnets, with private DNS on -- workloads in the VPC reach the service through their default SDK endpoints,..."
  - slug: "03-privatelink-service"
    rank: "03"
    title: "Third-Party PrivateLink Service"
    excerpt: "This preset connects a VPC privately to a vendor's PrivateLink service (a SaaS database, observability platform, or partner API published behind their NLB) -- the vendor's endpoint-service name goes..."
---

# VPC Endpoint Presets

Ready-to-deploy configuration presets for VPC Endpoint. Each preset is a complete manifest you can copy, customize, and deploy.
