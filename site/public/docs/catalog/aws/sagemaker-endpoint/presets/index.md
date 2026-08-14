---
title: "Presets"
description: "Ready-to-deploy configuration presets for SageMaker Endpoint"
type: "preset-list"
componentSlug: "sagemaker-endpoint"
componentTitle: "SageMaker Endpoint"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-serverless-endpoint"
    rank: "01"
    title: "Serverless Endpoint"
    excerpt: "This preset serves one model on serverless compute — SageMaker scales capacity with traffic and bills per inference, so an idle endpoint costs $0. The start-cheap shape for new and spiky workloads."
  - slug: "02-production-canary-endpoint"
    rank: "02"
    title: "Production Canary Endpoint"
    excerpt: "This preset runs two weighted instance-backed variants — 90% of traffic on the stable model, 10% on the candidate — with data capture feeding Model Monitor and every capacity change rolled blue/green..."
---

# SageMaker Endpoint Presets

Ready-to-deploy configuration presets for SageMaker Endpoint. Each preset is a complete manifest you can copy, customize, and deploy.
