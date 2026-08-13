---
title: "Presets"
description: "Ready-to-deploy configuration presets for Bedrock Custom Model"
type: "preset-list"
componentSlug: "bedrock-custom-model"
componentTitle: "Bedrock Custom Model"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-fine-tune-minimal"
    rank: "01"
    title: "Fine-Tune Minimal"
    excerpt: "This preset fine-tunes the cheapest fine-tunable Amazon base model (Titan Text Lite) for a single epoch — the pipeline-validation shape that proves your data format, role permissions, and S3 wiring..."
  - slug: "02-fine-tune-private-vpc"
    rank: "02"
    title: "Fine-Tune in a Private VPC"
    excerpt: "This preset runs the customization job with VPC-scoped data access and a customer-managed KMS key on the resulting model — the compliance posture for sensitive training data: traffic to S3 rides your..."
---

# Bedrock Custom Model Presets

Ready-to-deploy configuration presets for Bedrock Custom Model. Each preset is a complete manifest you can copy, customize, and deploy.
