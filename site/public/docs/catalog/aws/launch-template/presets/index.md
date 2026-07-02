---
title: "Presets"
description: "Ready-to-deploy configuration presets for Launch Template"
type: "preset-list"
componentSlug: "launch-template"
componentTitle: "Launch Template"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-web-server"
    rank: "01"
    title: "Web Server Fleet"
    excerpt: "This preset creates the launch blueprint for an auto-scaled web fleet: IMDSv2 enforced, an encrypted gp3 root volume, detailed monitoring for responsive scaling, and instance identity via an IAM..."
  - slug: "02-spot-flexible"
    rank: "02"
    title: "Spot-Flexible Workers"
    excerpt: "This preset describes compute by attributes instead of naming an instance type: any current-generation x86 type with 2-8 vCPUs and 4-16 GiB of memory qualifies. Paired with an `AwsAutoScalingGroup`..."
  - slug: "03-hardened"
    rank: "03"
    title: "Hardened Baseline"
    excerpt: "This preset is the strict-posture template for compliance-sensitive workloads: IMDSv2 with the hop limit locked to the host, root volume encrypted with a customer-managed KMS key (revocable,..."
---

# Launch Template Presets

Ready-to-deploy configuration presets for Launch Template. Each preset is a complete manifest you can copy, customize, and deploy.
