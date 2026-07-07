---
title: "CloudWatch Composite Alarm"
description: "CloudWatch Composite Alarm deployment documentation"
icon: "package"
order: 100
componentName: "awscloudwatchcompositealarm"
---

# AWS CloudWatch Composite Alarm

Deploys an AWS CloudWatch composite alarm — a boolean combination of other alarms' states that pages once for shared-cause outages, expresses alerting dependencies, and gates actions behind maintenance suppression. The composite has no metrics of its own: it re-evaluates its rule whenever any referenced alarm changes state.

## What Gets Created

- **CloudWatch Composite Alarm** — an `aws_cloudwatch_composite_alarm` with the specified rule expression, actions, and optional actions suppressor

## Prerequisites

- **Existing alarms** in the same account and region — the rule references them by name (compose from each `AwsCloudwatchAlarm`'s exported `alarm_name` output)
- **AWS credentials** configured via environment variables or Planton provider config
- **An SNS topic** if configuring actions
- **IAM permissions** — `cloudwatch:PutCompositeAlarm`, `cloudwatch:DeleteAlarms`, `cloudwatch:DescribeAlarms`, `cloudwatch:TagResource`

## Quick Start

Create a file `composite.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsCloudwatchCompositeAlarm
metadata:
  name: shared-cause
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsCloudwatchCompositeAlarm.shared-cause
spec:
  region: us-west-2
  alarmRule: 'ALARM("cpu-high") AND ALARM("error-rate-high")'
  alarmActions:
    - value: arn:aws:sns:us-west-2:123456789012:ops-critical
```

Deploy:

```shell
planton apply -f composite.yaml
```

This pages the critical channel once when BOTH constituent alarms are in ALARM.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region (must match the referenced alarms). | Required |
| `alarmRule` | `string` | Boolean expression over other alarms' states: `ALARM("name")` / `OK("name")` / `INSUFFICIENT_DATA("name")` with `AND`/`OR`/`NOT`, parentheses, and the `TRUE`/`FALSE` constants. | Required; max 10240 chars |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `alarmDescription` | `string` | — | What this composite represents and what to do when it fires. Max 1024 chars. |
| `actionsEnabled` | `bool` | `true` (AWS default) | Whether actions execute on state transitions. ForceNew — changing it replaces the alarm. |
| `alarmActions` | `StringValueOrRef[]` | `[]` | Actions on ALARM (SNS topics, SSM OpsItems). Max 5. |
| `okActions` | `StringValueOrRef[]` | `[]` | Actions on OK. Max 5. |
| `insufficientDataActions` | `StringValueOrRef[]` | `[]` | Actions on INSUFFICIENT_DATA. Max 5. |
| `actionsSuppressor` | `object` | — | Suppresses actions while a designated alarm is in ALARM: `alarm` (by name — reference an alarm's `alarm_name` output), `waitPeriod` and `extensionPeriod` in seconds. |

## Examples

### Maintenance-Suppressed Paging

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsCloudwatchCompositeAlarm
metadata:
  name: prod-api-paging
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: acme
    pulumi.planton.dev/project: platform
    pulumi.planton.dev/stack.name: prod.AwsCloudwatchCompositeAlarm.prod-api-paging
spec:
  region: us-west-2
  alarmRule: 'ALARM("api-5xx") OR ALARM("api-latency")'
  alarmActions:
    - value: arn:aws:sns:us-west-2:123456789012:ops-critical
  actionsSuppressor:
    alarm:
      valueFrom:
        kind: AwsCloudwatchAlarm
        name: maintenance-flag
        fieldPath: status.outputs.alarm_name
    waitPeriod: 60
    extensionPeriod: 120
```

While the `maintenance-flag` alarm is in ALARM (raised by the deploy pipeline), this composite still evaluates and records transitions — it just withholds the page.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `alarm_arn` | `string` | ARN of the composite alarm. |
| `alarm_name` | `string` | Name of the composite alarm — the join key parent composites use in their own rule expressions. |

## Related Components

- [AwsCloudwatchAlarm](/docs/catalog/aws/cloudwatch-alarm) — the constituent alarms this composite combines
- [AwsSnsTopic](/docs/catalog/aws/sns-topic) — the most common action target for notifications
- [AwsCloudwatchLogGroup](/docs/catalog/aws/cloudwatch-log-group) — source of log-derived metrics that feed constituent alarms
