---
title: "Lambda"
description: "Lambda deployment documentation"
icon: "package"
order: 100
componentName: "awslambda"
---

# AWS Lambda

Deploys an AWS Lambda function — zip archive from S3 or container image from ECR — with execution environment, VPC attachment, integrations, and function-scoped satellites (aliases, function URL, invoke permissions, async invoke config, recursion detection, runtime management) folded into the same spec. The function name is `metadata.name` (create-time immutable). Event sources attach through the separate `AwsLambdaEventSourceMapping` kind.

## What Gets Created

When you deploy an AwsLambda resource, Planton provisions:

- **Lambda function** — an `aws_lambda_function` with exactly one code source (`spec.s3` zip or `spec.image_uri` container), the referenced execution role, and any configured satellites materialized as their own AWS resources (aliases, function URL, resource-policy permissions, async invoke config, provisioned concurrency on aliases)
- **CloudWatch Logs** — when `logging_config` is unset, AWS creates and writes to `/aws/lambda/<function-name>` on first invoke (Planton does not pre-create a log group). Set `logging_config.log_group` to write into an `AwsCloudwatchLogGroup` you manage (retention, encryption, subscription filters)

Planton never creates an execution IAM role or invoke permissions unless you declare them: `role_arn` references an existing `AwsIamRole`, and `invoke_permissions` lists the resource-policy statements you need.

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **An IAM execution role** (`AwsIamRole`) trusting `lambda.amazonaws.com` with at minimum `AWSLambdaBasicExecutionRole` (add `AWSLambdaVPCAccessExecutionRole` for VPC attachment)
- **A deployment artifact** — zip in S3 (`spec.s3`) or container image in ECR (`spec.image_uri`); exactly one, never both
- **VPC subnets and security groups** (`AwsSubnet`, `AwsSecurityGroup`) only when the function needs private network access — both lists are required together

## Quick Start

Create a file `lambda.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsLambda
metadata:
  name: my-function
spec:
  region: us-west-2
  role_arn:
    valueFrom:
      kind: AwsIamRole
      name: lambda-exec
      fieldPath: status.outputs.role_arn
  runtime: nodejs22.x
  handler: index.handler
  s3:
    bucket:
      value: my-deploy-bucket
    key: functions/my-function.zip
```

Deploy:

```shell
planton apply -f lambda.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region (must match S3 bucket, VPC subnets, EFS access point). | Required; non-empty |
| `role_arn` | `StringValueOrRef` | Execution role. Defaults to `AwsIamRole` → `status.outputs.role_arn`. | Required |
| `s3` or `image_uri` | `object` / `string` | Exactly one code source. | CEL-enforced |
| `runtime` | `string` | Managed runtime (e.g. `nodejs22.x`, `python3.13`). Required for zip; empty for container. | Conditional |
| `handler` | `string` | Entrypoint for zip (e.g. `index.handler`). Required for zip; empty for container. | Conditional |
| `s3.bucket` | `StringValueOrRef` | S3 bucket for the zip. Defaults to `AwsS3Bucket` → `status.outputs.bucket_id`. | Required with `s3` |
| `s3.key` | `string` | Object key for the zip archive. | Required with `s3` |

### Common Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `description` | `string` | — | Console description (max 256 characters). |
| `memory_size_mb` | `int32` | 128 | Memory 128–10240 MB; CPU scales linearly. |
| `timeout_seconds` | `int32` | 3 | Max invocation time 1–900 seconds. |
| `architecture` | `string` | `x86_64` | `x86_64` or `arm64`. |
| `environment` | `map<string,string>` | — | Plain config vars (not for secrets). |
| `kms_key_arn` | `StringValueOrRef` | AWS-managed | Encrypts environment variables at rest. |
| `subnet_ids` | `StringValueOrRef[]` | — | VPC subnets for ENIs; requires `security_group_ids`. |
| `security_group_ids` | `StringValueOrRef[]` | — | Security groups for ENIs; requires `subnet_ids`. |
| `layer_arns` | `StringValueOrRef[]` | — | Up to five layer version ARNs (zip only). |
| `publish` | `bool` | `false` | Publish immutable version on each change (required for aliases, SnapStart). |
| `reserved_concurrent_executions` | `int32` | unreserved pool | `0` throttles all invocations; positive reserves capacity. |
| `logging_config` | `object` | AWS default on invoke | Log format, levels, and optional `log_group` ref. |
| `aliases` | `object[]` | — | Named version pointers with optional canary routing and provisioned concurrency. |
| `function_url` | `object` | — | Built-in HTTPS endpoint (`authorization_type` required). |
| `invoke_permissions` | `object[]` | — | Resource-policy statements for services/accounts to invoke. |
| `async_invoke_config` | `object` | — | Retry, event age, and on-success/on-failure destinations. |
| `recursive_loop` | `string` | `Terminate` | `Allow` or `Terminate` for self-invocation loops. |
| `runtime_management` | `object` | — | How runtime patches roll out (`Auto`, `FunctionUpdate`, `Manual`). |

Event sources (SQS, Kinesis, DynamoDB Streams, MSK, self-managed Kafka) are **not** on this spec — use [AwsLambdaEventSourceMapping](/docs/catalog/aws/lambda-event-source-mapping).

## Examples

### Container Image

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsLambda
metadata:
  name: image-function
spec:
  region: us-west-2
  role_arn:
    value: arn:aws:iam::123456789012:role/lambda-exec
  image_uri: 123456789012.dkr.ecr.us-west-2.amazonaws.com/my-func:latest
  architecture: arm64
  memory_size_mb: 512
  timeout_seconds: 30
```

### VPC-Connected with Managed Log Group

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsLambda
metadata:
  name: vpc-function
spec:
  region: us-west-2
  role_arn:
    valueFrom:
      kind: AwsIamRole
      name: lambda-vpc
      fieldPath: status.outputs.role_arn
  runtime: python3.13
  handler: app.handler
  s3:
    bucket:
      value: deploy-artifacts
    key: functions/vpc-function.zip
  subnet_ids:
    - valueFrom:
        kind: AwsSubnet
        name: private-a
        fieldPath: status.outputs.subnet_id
    - valueFrom:
        kind: AwsSubnet
        name: private-b
        fieldPath: status.outputs.subnet_id
  security_group_ids:
    - valueFrom:
        kind: AwsSecurityGroup
        name: lambda-sg
        fieldPath: status.outputs.security_group_id
  logging_config:
    log_format: JSON
    log_group:
      valueFrom:
        kind: AwsCloudwatchLogGroup
        name: lambda-logs
        fieldPath: status.outputs.log_group_name
```

### Aliases and Function URL

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsLambda
metadata:
  name: api-function
spec:
  region: us-west-2
  role_arn:
    valueFrom:
      kind: AwsIamRole
      name: lambda-exec
      fieldPath: status.outputs.role_arn
  runtime: nodejs22.x
  handler: index.handler
  s3:
    bucket:
      value: deploy-bucket
    key: api.zip
  publish: true
  aliases:
    - name: live
      function_version: "1"
  function_url:
    authorization_type: AWS_IAM
  invoke_permissions:
    - statement_id: allow-events
      principal: events.amazonaws.com
      source_arn: arn:aws:events:us-west-2:123456789012:rule/my-rule
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `function_arn` | `string` | Function ARN — join key for event-source mappings and IAM policies |
| `function_name` | `string` | Function name (matches `metadata.name`) |
| `invoke_arn` | `string` | API Gateway–shaped invocation ARN |
| `qualified_arn` | `string` | Qualified ARN of the latest published version (empty when `publish` is false) |
| `version` | `string` | Latest published version number (empty when `publish` is false) |
| `function_url` | `string` | HTTPS endpoint (empty when no function URL configured) |
| `alias_arns` | `map<string,string>` | Alias name → alias ARN |
| `log_group_name` | `string` | Log group receiving function logs (AWS default or `logging_config.log_group`) |

## Related Components

- [AwsIamRole](/docs/catalog/aws/iam-role) — execution role (required prerequisite)
- [AwsLambdaEventSourceMapping](/docs/catalog/aws/lambda-event-source-mapping) — SQS, Kinesis, DynamoDB Streams, MSK event sources
- [AwsS3Bucket](/docs/catalog/aws/s3-bucket) — hosts zip deployment packages
- [AwsCloudwatchLogGroup](/docs/catalog/aws/cloudwatch-log-group) — managed log destination when you need retention or encryption
- [AwsKmsKey](/docs/catalog/aws/kms-key) — encrypts environment variables at rest
- [AwsSubnet](/docs/catalog/aws/subnet) / [AwsSecurityGroup](/docs/catalog/aws/security-group) — VPC attachment for private resources
