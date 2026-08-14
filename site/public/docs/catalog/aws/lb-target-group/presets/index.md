---
title: "Presets"
description: "Ready-to-deploy configuration presets for LB Target Group"
type: "preset-list"
componentSlug: "lb-target-group"
componentTitle: "LB Target Group"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-ecs-service-http"
    rank: "01"
    title: "ECS Service HTTP"
    excerpt: "This preset creates an HTTP target group with the `ip` target type -- the shape ECS awsvpc tasks and Kubernetes pod-IP integrations register into. The group stays empty at deploy time; the ECS..."
  - slug: "02-nlb-tcp-passthrough"
    rank: "02"
    title: "NLB TCP Pass-Through"
    excerpt: "This preset creates a TCP target group for a Network Load Balancer fronting a non-HTTP protocol (shown here on 5432, the PostgreSQL port). Traffic passes through at Layer 4 untouched; the application..."
  - slug: "03-lambda-function"
    rank: "03"
    title: "Lambda Function Target"
    excerpt: "This preset creates a Lambda target group: the ALB invokes the function directly for each request, so there is no port, protocol, or VPC to configure. The function is the group's single target,..."
  - slug: "04-quic-passthrough"
    rank: "04"
    title: "QUIC / HTTP/3 Pass-Through"
    excerpt: "This preset creates a `TCP_QUIC` target group for a Network Load Balancer serving HTTP/3 traffic: QUIC connections pass through natively while the same group serves clients that fall back to TCP on..."
---

# LB Target Group Presets

Ready-to-deploy configuration presets for LB Target Group. Each preset is a complete manifest you can copy, customize, and deploy.
