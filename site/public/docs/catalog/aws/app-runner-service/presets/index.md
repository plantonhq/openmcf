---
title: "Presets"
description: "Ready-to-deploy configuration presets for App Runner Service"
type: "preset-list"
componentSlug: "app-runner-service"
componentTitle: "App Runner Service"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-basic-public-image"
    rank: "01"
    title: "Basic Public Image"
    excerpt: "This preset creates an App Runner service from a public ECR image with all default settings. It is the simplest possible App Runner deployment -- no IAM roles, no VPC, no encryption configuration...."
  - slug: "02-production-vpc-encrypted"
    rank: "02"
    title: "Production VPC-Connected and Encrypted Service"
    excerpt: "This preset creates a production-grade App Runner service with private ECR image, VPC egress, customer-managed KMS encryption, tuned auto scaling, and HTTP health checks. It represents the..."
  - slug: "03-github-code-source"
    rank: "03"
    title: "GitHub Code Source (Node.js)"
    excerpt: "This preset creates an App Runner service that deploys directly from a GitHub repository using the Node.js 22 managed runtime. App Runner clones the repository, runs the build command, and starts the..."
  - slug: "04-private-vpc-ingress"
    rank: "04"
    title: "Private Internal Service (VPC Ingress)"
    excerpt: "This preset creates an App Runner service with NO public endpoint: incoming traffic is disabled on the internet side, and the service is published into a VPC through an interface VPC endpoint (AWS..."
---

# App Runner Service Presets

Ready-to-deploy configuration presets for App Runner Service. Each preset is a complete manifest you can copy, customize, and deploy.
