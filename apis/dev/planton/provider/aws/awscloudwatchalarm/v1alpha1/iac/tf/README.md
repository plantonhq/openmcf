# AwsCloudwatchAlarm Terraform Module

Provisions an AWS CloudWatch metric alarm supporting all three metric definition modes — simple metric, metric math / anomaly detection queries, and PromQL against an Amazon Managed Service for Prometheus workspace — with configurable M-of-N evaluation and multi-action notifications.

## Resources Created

- `aws_cloudwatch_metric_alarm.this` — The CloudWatch metric alarm

## Inputs

| Variable | Description |
|----------|-------------|
| `metadata` | Resource metadata (name, org, env, id, labels) |
| `spec` | AwsCloudwatchAlarmSpec — desired configuration |

## Outputs

| Output | Description |
|--------|-------------|
| `alarm_arn` | The ARN of the metric alarm |
| `alarm_name` | The name of the metric alarm |

## Provider

Requires `hashicorp/aws` provider version `>= 6.43.0` (the floor where the
PromQL `evaluation_criteria` surface stabilized — 6.42.0 introduced it with a
plan-time regression that 6.43.0 fixed).
