---
title: "Bedrock Provisioned Throughput"
description: "Bedrock Provisioned Throughput deployment documentation"
icon: "package"
order: 100
componentName: "awsbedrockprovisionedthroughput"
---

# AWS Bedrock Provisioned Throughput

Dedicated model serving capacity — the required path for fine-tuned
custom models, purchased in model units with optional commitment terms.

## What Gets Created

- A provisioned model: reserved throughput for the referenced custom
  model (or a provisionable foundation model), addressable by its own
  ARN.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Bedrock permissions
  (`bedrock:CreateProvisionedModelThroughput` and siblings).
- Typically: an `AwsBedrockCustomModel` whose output ARN this purchase
  serves.

### AWS Account

- A no-commitment model-unit quota above zero (Service Quotas; the
  default is often 0), or a commitment-term purchase.
- A clear cost sign-off — capacity bills from creation, and committed
  terms bill in full.

## Deploy

### Console

Create the resource from the AWS catalog, reference the model, set the
units, and deploy.

### CLI

```bash
planton apply -f provisioned-throughput.yaml
```

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsBedrockProvisionedThroughput
metadata:
  name: support-model-capacity
spec:
  region: us-east-1
  modelArn:
    valueFrom:
      kind: AwsBedrockCustomModel
      name: support-titan-ft
  modelUnits: 1
```

## Operational Notes

- **No-commitment first.** Validate integration on hourly billing before
  committing to a term; committed purchases cannot be canceled OR
  destroyed until the term ends.
- **A model unit is model-specific** — AWS quotes each model's
  tokens-per-minute per unit in the console.
- **Everything except tags is create-time-immutable** — resizing units or
  changing the model is a replacement (and a committed purchase blocks
  that replacement until its term lapses).
