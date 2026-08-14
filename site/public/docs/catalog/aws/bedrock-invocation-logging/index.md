---
title: "Bedrock Invocation Logging"
description: "Bedrock Invocation Logging deployment documentation"
icon: "package"
order: 100
componentName: "awsbedrockinvocationlogging"
---

# AWS Bedrock Invocation Logging

The region's Bedrock model invocation logging — the audit trail of
every model call, delivered to CloudWatch Logs and/or S3. A settings
singleton: one configuration per region, at most one instance
deployed per region.

## What Gets Managed

- Which invocation data types are captured (text, image, embedding,
  video).
- CloudWatch delivery (log group + IAM role, with S3 spillover for
  oversized payloads) and/or S3 delivery.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Bedrock and logging permissions.

### AWS Account

- For CloudWatch delivery: a log group
  ([AWS CloudWatch Log Group](/cloud-catalog/aws-cloudwatch-log-group))
  and an IAM role ([AWS IAM Role](/cloud-catalog/aws-iam-role))
  trusting `bedrock.amazonaws.com` with write access to it.
- For S3 delivery: a bucket ([AWS S3 Bucket](/cloud-catalog/aws-s3-bucket))
  whose policy lets `bedrock.amazonaws.com` put objects.

## Deploy

### Console

Create the resource from the AWS catalog, pick the destinations, and
deploy.

### CLI

```bash
planton apply -f bedrock-invocation-logging.yaml
```

## After Deploy

- Every model invocation in the region is logged to the configured
  destinations from the next call on.
- Destroying this component deletes the configuration — logging
  stops region-wide.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
