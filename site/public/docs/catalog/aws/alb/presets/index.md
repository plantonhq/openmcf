---
title: "Presets"
description: "Ready-to-deploy configuration presets for ALB"
type: "preset-list"
componentSlug: "alb"
componentTitle: "ALB"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-internet-facing"
    rank: "01"
    title: "Internet-Facing ALB"
    excerpt: "This preset creates the standard public entry point: an internet-facing Application Load Balancer across two public subnets, with an explicit security group, deletion protection, and Route53 alias..."
  - slug: "02-internal-hardened"
    rank: "02"
    title: "Internal Hardened ALB"
    excerpt: "This preset creates a VPC-internal Application Load Balancer with the HTTP hardening turned all the way up: invalid request headers are dropped, desync mitigation runs in `strictest` mode, and every..."
  - slug: "03-reserved-capacity"
    rank: "03"
    title: "Reserved Capacity for a Traffic Event"
    excerpt: "This preset pre-provisions an internet-facing ALB for a known traffic surge -- a product launch, a ticket sale, a marketing spike -- instead of waiting for the load balancer to scale organically. It..."
---

# ALB Presets

Ready-to-deploy configuration presets for ALB. Each preset is a complete manifest you can copy, customize, and deploy.
