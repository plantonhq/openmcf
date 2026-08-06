# AWS Lambda

AWS Lambda runs application code in response to events without provisioning servers. This document describes how Planton models Lambda functions: what the spec covers, how code is packaged, and how the function composes with other resources.

## Function identity

The AWS function name is **`metadata.name`**. It is create-time immutable — there is no separate `function_name` field in the spec. Other resources and clients join to the function through stack outputs (`function_arn`, `function_name`, `invoke_arn`, `alias_arns`).

## Code source

Provide **exactly one** of:

| Source | Fields | When to use |
|--------|--------|-------------|
| **S3 zip** | `s3.bucket`, `s3.key`, plus `runtime` and `handler` | Default for most functions — fast cold starts, managed runtimes |
| **Container image** | `image_uri` (no `runtime`/`handler`) | Large dependencies (up to 10 GB), custom OS packages, image pipelines |

Optional zip refinements:

- `source_code_hash` — declarative code updates from your build pipeline
- `s3.object_version` — pin a versioned S3 object
- `source_kms_key_arn` — encrypt the deployment package at rest
- `code_signing_config_arn` — enforce signed packages

## Execution environment

Core tuning fields:

- `memory_size_mb` (128–10240) — also scales CPU; raising memory often reduces total cost for CPU-bound code
- `timeout_seconds` (1–900)
- `architecture` — `arm64` (Graviton, typically cheaper) or `x86_64`
- `ephemeral_storage_mb` — `/tmp` scratch space
- `environment` — plain configuration (resolve secrets at runtime via SSM/Secrets Manager, not here)
- `kms_key_arn` — customer-managed encryption for environment variables
- `layer_arns` — shared dependencies (zip only, up to five)

## VPC attachment

To reach private resources (RDS, ElastiCache, internal APIs), set **`subnet_ids` and `security_group_ids` together**. Attaching to a VPC removes default internet access — route outbound traffic through a NAT gateway or VPC endpoints as needed. The execution role needs `AWSLambdaVPCAccessExecutionRole`.

## Logging

Planton does **not** pre-create a CloudWatch log group. Behavior:

- **Unset `logging_config`** — AWS creates `/aws/lambda/<function-name>` on first invoke and writes plain-text logs
- **Set `logging_config.log_group`** — write to an `AwsCloudwatchLogGroup` you manage (retention, KMS encryption, subscription filters). The execution role needs `logs:CreateLogStream` and `logs:PutLogEvents` on that group
- **Structured logs** — set `logging_config.log_format: JSON` and optional `application_log_level` / `system_log_level`

## IAM and invoke permissions

- **`role_arn`** (required) — references an `AwsIamRole` the function assumes. Planton never creates or modifies role policies; the role must already trust `lambda.amazonaws.com` and carry the permissions your code needs
- **`invoke_permissions`** — resource-policy statements granting services (S3, SNS, EventBridge) or other principals permission to invoke. Each entry materializes as its own permission keyed by `statement_id`

## Satellites folded into the function spec

These are separate AWS resources but honestly part of function configuration — they live on the same spec:

| Field | Purpose |
|-------|---------|
| `aliases` | Stable invocation targets; canary routing; provisioned concurrency (requires `publish: true`) |
| `function_url` | Built-in HTTPS endpoint without API Gateway |
| `async_invoke_config` | Async retry, max event age, success/failure destinations |
| `runtime_management` | Control how AWS applies runtime patches |
| `recursive_loop` | Terminate or allow self-invocation loops |

## Event sources

SQS queues, Kinesis streams, DynamoDB streams, MSK topics, and other poller-driven sources attach through the separate **`AwsLambdaEventSourceMapping`** kind — not on the function spec. A function may have many mappings; pausing or retuning batching does not require changing the function.

## Versioning and concurrency

- **`publish: true`** — creates an immutable version on each code/config change; required for aliases and SnapStart
- **`reserved_concurrent_executions`** — unset draws from the account pool; `0` is a kill switch; a positive value reserves that many concurrent executions
- **`snap_start: true`** — JVM and slow-init runtimes; requires `publish` and invocation through a version or alias

## Prerequisites checklist

1. `AwsIamRole` with `lambda.amazonaws.com` trust and `AWSLambdaBasicExecutionRole`
2. Deployment artifact in S3 or ECR
3. For VPC: `AwsSubnet` + `AwsSecurityGroup` references
4. For managed logging: `AwsCloudwatchLogGroup` reference in `logging_config.log_group`
5. For event-driven workloads: one or more `AwsLambdaEventSourceMapping` resources

## Stack outputs

| Output | Use |
|--------|-----|
| `function_arn` | Event-source mappings, IAM policies, trigger configs |
| `function_name` | SDK `Invoke` calls |
| `invoke_arn` | API Gateway integrations |
| `qualified_arn`, `version` | Published version tracking |
| `function_url` | Direct HTTPS access |
| `alias_arns` | Stable routing (`status.outputs.alias_arns.live`) |
| `log_group_name` | Observability wiring |

## IaC backends

Planton renders this spec to Pulumi (Go) and Terraform modules under `iac/`. Provider credentials are supplied via stack input, not the resource spec.
