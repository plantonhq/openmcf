---
title: "Presets"
description: "Ready-to-deploy configuration presets for Route53 Health Check"
type: "preset-list"
componentSlug: "route53-health-check"
componentTitle: "Route53 Health Check"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-https-endpoint"
    rank: "01"
    title: "HTTPS Endpoint Health Check"
    excerpt: "This preset creates the standard internet-facing HTTPS probe: Route 53's global checker fleet requests your application's health endpoint and the check stays healthy while the endpoint answers..."
  - slug: "02-cloudwatch-alarm"
    rank: "02"
    title: "CloudWatch Alarm Health Check"
    excerpt: "This preset creates a health check that mirrors a CloudWatch alarm's state instead of probing an endpoint. It is the pattern for PRIVATE resources the Route 53 checker fleet cannot reach — internal..."
---

# Route53 Health Check Presets

Ready-to-deploy configuration presets for Route53 Health Check. Each preset is a complete manifest you can copy, customize, and deploy.
