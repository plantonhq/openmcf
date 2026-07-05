---
title: "Zip-Based Lambda Function"
description: "This preset creates a Lambda function deployed from a zip archive in S3. The function name comes from `metadata.name` (create-time immutable in AWS). It uses the Node.js 22.x runtime with 256 MB..."
type: "preset"
rank: "01"
presetSlug: "01-zip-basic"
componentSlug: "lambda"
componentTitle: "Lambda"
provider: "aws"
icon: "package"
order: 1
---

# Zip-Based Lambda Function

This preset creates a Lambda function deployed from a zip archive in S3. The function name comes from `metadata.name` (create-time immutable in AWS). It uses the Node.js 22.x runtime with 256 MB memory and a 30-second timeout — the most common deployment model for event handlers, API endpoints, and automation scripts.

## When to Use

- Event-driven functions triggered by S3, SNS, EventBridge, or API integrations
- Lightweight API endpoints or webhook handlers (pair with `function_url` or API Gateway)
- Functions written in Node.js, Python, Java, Go, .NET, or Ruby using a managed runtime

## Key Configuration Choices

- **S3 code source** (`spec.s3`) — deployment package stored in S3; your CI/CD pipeline uploads the zip and updates the key (or `source_code_hash` for declarative rolls)
- **Node.js 22.x** (`runtime: nodejs22.x`) — change to `python3.13`, `java21`, `provided.al2023`, etc. for other languages
- **256 MB memory** (`memory_size_mb: 256`) — CPU scales with memory; increase for compute-intensive work
- **30-second timeout** (`timeout_seconds: 30`) — suitable for API handlers; increase up to 900 for long-running jobs
- **Composed execution role** (`role_arn.valueFrom`) — references an `AwsIamRole`; the role must trust `lambda.amazonaws.com` and carry `AWSLambdaBasicExecutionRole`
- **No VPC** — runs in Lambda's managed network; add `subnet_ids` and `security_group_ids` together for private resource access

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aws-region>` | AWS region for the function and S3 bucket | Your deployment region |
| `<execution-role-name>` | Name of the `AwsIamRole` resource for execution | `AwsIamRole` metadata.name |
| `<deployment-bucket>` | S3 bucket containing the zip | Your CI/CD pipeline or `AwsS3Bucket` outputs |
| `<deployment-package-key>` | S3 object key (e.g. `functions/my-function/v1.0.0.zip`) | Your CI/CD pipeline |

## Related Presets

- **02-container-basic** — use instead for container images (large dependencies, custom OS packages, image-based pipelines)
- **AwsLambdaEventSourceMapping** — attach SQS, Kinesis, DynamoDB Streams, or MSK as event sources (separate kind, not on the function spec)
