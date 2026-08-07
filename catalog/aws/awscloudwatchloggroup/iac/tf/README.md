# AwsCloudwatchLogGroup Terraform Module

Provisions an AWS CloudWatch Logs log group with configurable retention, KMS
encryption, log group class, deletion protection, and the group-scoped
satellites: metric filters, subscription filters, a data protection policy,
and a field index policy.

## Resources Created

- `aws_cloudwatch_log_group.this` — The CloudWatch log group
- `aws_cloudwatch_log_metric_filter.this` — One per named metric filter
- `aws_cloudwatch_log_subscription_filter.this` — One per named subscription filter (max 2)
- `aws_cloudwatch_log_data_protection_policy.this` — When a data protection policy is configured
- `aws_cloudwatch_log_index_policy.this` — When a field index policy is configured

## Inputs

| Variable | Description |
|----------|-------------|
| `metadata` | Resource metadata (name, org, env, id, labels) |
| `spec` | AwsCloudwatchLogGroupSpec — desired configuration |

## Outputs

| Output | Description |
|--------|-------------|
| `log_group_arn` | The ARN of the log group (without the `:*` suffix) |
| `log_group_name` | The name of the log group |

## Provider

Requires `hashicorp/aws` provider version `>= 6.25.0` (the floor where
`deletion_protection_enabled` landed).
