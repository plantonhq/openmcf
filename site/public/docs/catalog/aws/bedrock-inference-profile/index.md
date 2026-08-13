---
title: "Bedrock Inference Profile"
description: "Bedrock Inference Profile deployment documentation"
icon: "package"
order: 100
componentName: "awsbedrockinferenceprofile"
---

# AWS Bedrock Inference Profile

Per-application handles over Bedrock models — track, cost-allocate, and
IAM-scope each consumer's model usage through its own ARN.

## What Gets Created

- An application inference profile routing to a foundation model or an
  AWS system-defined cross-region profile. Free to create; invocations
  bill at the underlying model's rates.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Bedrock permissions
  (`bedrock:CreateInferenceProfile` and read/delete siblings).

### AWS Account

- Access to the underlying model (see `AwsBedrockModelAccess` for models
  that need a marketplace agreement).

## Deploy

### Console

Create the resource from the AWS catalog, name it after the consuming
application, point it at the model, and deploy.

### CLI

```bash
planton apply -f inference-profile.yaml
```

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsBedrockInferenceProfile
metadata:
  name: checkout-nova
spec:
  region: us-west-2
  description: Cost tracking for the checkout service
  sourceArn: arn:aws:bedrock:us-west-2:123456789012:inference-profile/us.amazon.nova-micro-v1:0
```

## Operational Notes

- **One profile per consumer** is the pattern — the profile's tags carry
  the cost-allocation identity, and IAM policies grant each application
  `bedrock:InvokeModel` on ITS profile ARN only.
- **Replacement changes the ARN.** Every spec field is
  create-time-immutable; consumers pinning the ARN must move when the
  profile is replaced.
- **Cross-region routing is inherited, not configured** — source the
  profile from AWS's system-defined geo profile and invocations ride
  AWS's cross-region capacity pools.
