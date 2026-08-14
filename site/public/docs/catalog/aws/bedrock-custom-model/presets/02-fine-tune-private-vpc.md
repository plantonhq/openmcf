---
title: "Fine-Tune in a Private VPC"
description: "This preset runs the customization job with VPC-scoped data access and a customer-managed KMS key on the resulting model — the compliance posture for sensitive training data: traffic to S3 rides your..."
type: "preset"
rank: "02"
presetSlug: "02-fine-tune-private-vpc"
componentSlug: "bedrock-custom-model"
componentTitle: "Bedrock Custom Model"
provider: "aws"
icon: "package"
order: 2
---

# Fine-Tune in a Private VPC

This preset runs the customization job with VPC-scoped data access and a
customer-managed KMS key on the resulting model — the compliance posture
for sensitive training data: traffic to S3 rides your subnets and security
groups, and the model at rest is under your key.

## When to Use

- Training data classified as confidential or regulated
- Accounts whose security baseline requires customer-managed keys and
  private data paths

## Key Configuration Choices

- **Both VPC members are required together** — subnets for placement,
  security groups for reachability (the S3 route typically needs a
  gateway VPC endpoint in those subnets' route tables).
- **The job role additionally needs KMS permissions** on the model key.
- **An explicit versioned `jobName`** — this shape tends to be re-run as
  data grows, and job names are single-use forever.

## After Deployment

Reference `custom_model_arn` from an `AwsBedrockProvisionedThroughput`
resource to serve the model.
