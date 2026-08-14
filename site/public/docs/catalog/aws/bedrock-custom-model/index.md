---
title: "Bedrock Custom Model"
description: "Bedrock Custom Model deployment documentation"
icon: "package"
order: 100
componentName: "awsbedrockcustommodel"
---

# AWS Bedrock Custom Model

Fine-tune, continue pre-training, or distill a Bedrock foundation model on
your own data — declaratively, with the resulting model wired into the
rest of your AI stack by reference.

## What Gets Created

- A Bedrock model-customization job (the deploy starts it; it runs
  asynchronously and bills per token processed).
- The resulting custom model, optionally encrypted under your KMS key,
  exported as `custom_model_arn` for provisioned-throughput purchases.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Bedrock model-customization permissions.
- An `AwsIamRole` for the job (trusting `bedrock.amazonaws.com`), or a
  role ARN.

### AWS Account

- Model customization available in the region (us-east-1 has the widest
  base-model coverage) and the base model eligible for the chosen
  customization type.
- Training data staged in S3 in the format the customization type
  requires.

## Deploy

### Console

Create the resource from the AWS catalog, pick the base model and training
data locations, and deploy — then track the job through `job_status`.

### CLI

```bash
planton apply -f custom-model.yaml
```

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsBedrockCustomModel
metadata:
  name: support-titan-ft
spec:
  region: us-east-1
  baseModelArn: arn:aws:bedrock:us-east-1::foundation-model/amazon.titan-text-lite-v1
  hyperparameters:
    epochCount: "1"
  roleArn:
    valueFrom:
      kind: AwsIamRole
      name: bedrock-customization
  trainingDataS3Uri: s3://my-bucket/data/train.jsonl
  outputDataS3Uri: s3://my-bucket/output/
```

## Operational Notes

- **Training cost is real** — it scales with data size, epochs, and the
  base model. Start with one epoch on a small dataset to validate the
  pipeline before the full run.
- **The deploy returns while the job runs.** `job_status` reports
  InProgress/Completed/Failed; the custom model is usable only after
  Completed.
- **Job names never come back.** AWS reserves them forever — a
  destroy/recreate needs a fresh `job_name`.
- **Serving needs capacity.** Fine-tuned models serve through Provisioned
  Throughput — buy it with `AwsBedrockProvisionedThroughput` referencing
  this component's `custom_model_arn`.
