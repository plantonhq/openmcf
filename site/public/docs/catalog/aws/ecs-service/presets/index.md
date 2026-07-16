---
title: "Presets"
description: "Ready-to-deploy configuration presets for ECS Service"
type: "preset-list"
componentSlug: "ecs-service"
componentTitle: "ECS Service"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-web-service"
    rank: "01"
    title: "Web Service"
    excerpt: "This preset runs a load-balanced web service: a referenced task-definition revision scheduled into a referenced cluster, registered into a first-class `AwsLbTargetGroup` that an `AwsLbListener` (or..."
  - slug: "02-cost-optimized-spot"
    rank: "02"
    title: "Cost-Optimized Spot Blend"
    excerpt: "This preset blends Fargate capacity instead of naming a launch type: one guaranteed on-demand task as the base, then a 1:4 on-demand:Spot split for everything above it -- roughly 70% Spot at a ~70%..."
  - slug: "03-blue-green"
    rank: "03"
    title: "Blue/Green with Canary"
    excerpt: "This preset turns on ECS-native blue/green deployments: new revisions stand up in the green target group, take 5% canary traffic for 5 minutes, then all traffic, and bake 10 minutes before the blue..."
---

# ECS Service Presets

Ready-to-deploy configuration presets for ECS Service. Each preset is a complete manifest you can copy, customize, and deploy.
