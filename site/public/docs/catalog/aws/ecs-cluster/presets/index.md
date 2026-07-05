---
title: "Presets"
description: "Ready-to-deploy configuration presets for ECS Cluster"
type: "preset-list"
componentSlug: "ecs-cluster"
componentTitle: "ECS Cluster"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-fargate-standard"
    rank: "01"
    title: "Standard Fargate Cluster"
    excerpt: "This preset creates an ECS cluster on AWS Fargate with enhanced Container Insights and audited ECS Exec. Fargate eliminates instance management entirely -- AWS runs the compute -- and the cluster..."
  - slug: "02-fargate-cost-optimized"
    rank: "02"
    title: "Fargate Cost-Optimized Cluster"
    excerpt: "This preset creates an ECS cluster with both Fargate and Fargate Spot capacity providers, using a weighted strategy that runs approximately 80% of scaled tasks on Spot for significant cost savings..."
  - slug: "03-ec2-capacity"
    rank: "03"
    title: "EC2-Backed Cluster"
    excerpt: "This preset adds EC2 capacity to an ECS cluster by wrapping a referenced auto-scaling group as a capacity provider. ECS's managed scaling drives the group's instance count from task demand -- you..."
---

# ECS Cluster Presets

Ready-to-deploy configuration presets for ECS Cluster. Each preset is a complete manifest you can copy, customize, and deploy.
