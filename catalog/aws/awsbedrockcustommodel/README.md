<p align="center">
  <img src="logo.svg" alt="AWS Bedrock Custom Model" width="80"/>
</p>

# AWS Bedrock Custom Model

Create and manage [Amazon Bedrock custom models](https://docs.aws.amazon.com/bedrock/latest/userguide/custom-models.html) —
foundation models customized with your training data through a
model-customization job (fine-tuning, continued pre-training, or
distillation).

## What Gets Created

- **A model-customization job** — deploying this component STARTS the
  job; it runs asynchronously in AWS (minutes to many hours) and bills per
  token processed.
- **The custom model** the job produces, encrypted under a Bedrock-managed
  key or your own KMS key.

## Everything Is Create-Time-Immutable

A customization job cannot be altered once started: ANY spec change
destroys the job record and custom model and starts a new job. Job names
are unique per account for all time — set `job_name` explicitly when
re-running a customization, because the derived default (metadata.name)
can only be used once.

## Serving the Model

A fine-tuned custom model needs
[Provisioned Throughput](../awsbedrockprovisionedthroughput/) to serve
traffic — reference this component's `custom_model_arn` output from an
`AwsBedrockProvisionedThroughput` resource.

## Prerequisites

- An IAM role trusting `bedrock.amazonaws.com` with read access to the
  training/validation S3 locations and write access to the output
  location (reference an `AwsIamRole` or pass an ARN).
- Training data in S3, formatted per customization type (JSONL
  prompt/completion pairs for fine-tuning).

## Quick Start

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
    batchSize: "1"
    learningRate: "0.00001"
  roleArn:
    valueFrom:
      kind: AwsIamRole
      name: bedrock-customization
  trainingDataS3Uri: s3://my-bucket/data/train.jsonl
  outputDataS3Uri: s3://my-bucket/output/
```

## Spec Reference

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.
