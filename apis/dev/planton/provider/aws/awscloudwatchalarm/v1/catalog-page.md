# AWS CloudWatch Alarm

Deploys an AWS CloudWatch metric alarm that monitors a metric, a metric math expression, or a PromQL query against an Amazon Managed Service for Prometheus workspace, and triggers actions when the alarm condition holds. Supports M-of-N evaluation to reduce false positives, anomaly detection bands, cross-account metrics, and can target SNS topics, Auto Scaling policies, or EC2 automation actions.

## What Gets Created

When you deploy an AwsCloudwatchAlarm resource, Planton provisions:

- **CloudWatch Metric Alarm** — an `aws_cloudwatch_metric_alarm` resource configured with the specified metric source (single metric or metric math queries), threshold, evaluation window, and actions

No additional sub-resources are created. The alarm is a standalone monitoring resource.

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **A metric to monitor** — the metric must exist in CloudWatch (published by an AWS service or custom application)
- **An SNS topic** if configuring alarm actions (the most common action target)
- **IAM permissions** — `cloudwatch:PutMetricAlarm`, `cloudwatch:DeleteAlarms`, `cloudwatch:DescribeAlarms`, `cloudwatch:TagResource`

## Quick Start

Create a file `alarm.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsCloudwatchAlarm
metadata:
  name: cpu-high
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsCloudwatchAlarm.cpu-high
spec:
  region: us-west-2
  comparisonOperator: GreaterThanThreshold
  evaluationPeriods: 3
  threshold: 80.0
  metricName: CPUUtilization
  namespace: AWS/EC2
  period: 300
  statistic: Average
```

Deploy:

```shell
planton apply -f alarm.yaml
```

This creates an alarm that triggers when EC2 CPU utilization exceeds 80% for 3 consecutive 5-minute periods.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region where the alarm will be created (e.g., `us-west-2`). | Required |
| `comparisonOperator` | `string` | Arithmetic operation comparing the statistic to the threshold. Required for simple-metric and metric-query alarms; must be omitted for PromQL alarms. | One of: `GreaterThanOrEqualToThreshold`, `GreaterThanThreshold`, `LessThanThreshold`, `LessThanOrEqualToThreshold`, `LessThanLowerOrGreaterThanUpperThreshold`, `LessThanLowerThreshold`, `GreaterThanUpperThreshold` |
| `evaluationPeriods` | `int` | Number of consecutive periods over which data is compared to the threshold. Required for simple-metric and metric-query alarms; must be omitted for PromQL alarms. | >= 1 |

Exactly one of the following metric source modes is required:

**Simple Metric Mode** — set `metricName`, `namespace`, `period`, and one of `statistic` or `extendedStatistic`:

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `metricName` | `string` | CloudWatch metric name (e.g., `CPUUtilization`, `RequestCount`). | Max 255 chars. Mutually exclusive with `metricQueries` and `evaluationCriteria`. |
| `namespace` | `string` | Metric namespace (e.g., `AWS/EC2`, `AWS/SQS`). | Required when `metricName` is set. Max 255 chars. |
| `period` | `int` | Period in seconds for statistic evaluation. | Required when `metricName` is set. Valid: 10, 20, 30, or multiple of 60. |
| `statistic` | `string` | Standard statistic. | One of: `SampleCount`, `Average`, `Sum`, `Minimum`, `Maximum`. Mutually exclusive with `extendedStatistic`. |

**Metric Query Mode** — set `metricQueries` for metric math or anomaly detection:

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `metricQueries` | `object[]` | Metric math expressions or multi-metric queries. | Max 20 items. Exactly one query sets `returnData: true`; each query sets exactly one of `expression`/`metric`. Mutually exclusive with `metricName` and `evaluationCriteria`. |

**PromQL Mode** — set `evaluationCriteria` to alarm on an Amazon Managed Service for Prometheus workspace:

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `evaluationCriteria.promqlCriteria.query` | `string` | PromQL query producing a boolean-style firing signal; the query itself is the alarm condition. | Required in this mode. Max 10000 chars. |
| `evaluationCriteria.promqlCriteria.pendingPeriod` | `int` | Seconds the query must continuously fire before ALARM (Prometheus `for:`). | 0–86400. |
| `evaluationCriteria.promqlCriteria.recoveryPeriod` | `int` | Seconds the query must be quiet before returning to OK. | 0–86400. |
| `evaluationInterval` | `int` | Query evaluation cadence in seconds. | 10, 20, 30, or multiple of 60. Only valid with `evaluationCriteria`. |

In PromQL mode the threshold fields (`comparisonOperator`, `evaluationPeriods`, `threshold`, `datapointsToAlarm`, and the simple-metric fields) must be omitted — the query expresses the condition.

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `threshold` | `double` | — | Static threshold value. Mutually exclusive with `thresholdMetricId`. |
| `thresholdMetricId` | `string` | — | ID of the ANOMALY_DETECTION_BAND function for anomaly detection alarms. Mutually exclusive with `threshold`. |
| `datapointsToAlarm` | `int` | Same as `evaluationPeriods` | Number of breaching data points to trigger alarm (M-of-N evaluation). Must be <= `evaluationPeriods`. |
| `treatMissingData` | `string` | `missing` | How missing data is treated. One of: `missing`, `ignore`, `breaching`, `notBreaching`. |
| `actionsEnabled` | `bool` | `true` | Whether actions execute on state transitions. Set to `false` during tuning or maintenance. |
| `extendedStatistic` | `string` | — | Percentile statistic (e.g., `p95`, `p99.9`, `IQM`). Mutually exclusive with `statistic`. |
| `dimensions` | `map<string,string>` | — | Dimensions to narrow the metric to a specific resource. |
| `unit` | `string` | — | Filters data points to matching unit. |
| `alarmActions` | `StringValueOrRef[]` | `[]` | Actions for ALARM transitions. Can reference AwsSnsTopic via `valueFrom`. Max 5. |
| `okActions` | `StringValueOrRef[]` | `[]` | Actions for OK transitions. Can reference AwsSnsTopic via `valueFrom`. Max 5. |
| `insufficientDataActions` | `StringValueOrRef[]` | `[]` | Actions for INSUFFICIENT_DATA transitions. Can reference AwsSnsTopic via `valueFrom`. Max 5. |
| `alarmDescription` | `string` | — | Human-readable description. Max 1024 chars. |
| `evaluateLowSampleCountPercentiles` | `string` | — | Percentile alarm behavior with low sample counts. One of: `evaluate`, `ignore`. |

### Metric Query Fields

Each item in `metricQueries`:

| Field | Type | Description |
|-------|------|-------------|
| `id` | `string` | **Required.** Variable name for expressions (e.g., `m1`, `errors`). Must start with lowercase letter. |
| `expression` | `string` | Metric math expression (e.g., `m1/m2*100`). Mutually exclusive with `metric`. |
| `metric` | `object` | Raw metric definition. Mutually exclusive with `expression`. |
| `metric.metricName` | `string` | **Required.** Metric name. |
| `metric.namespace` | `string` | **Required.** Metric namespace. |
| `metric.period` | `int` | **Required.** Period in seconds. |
| `metric.stat` | `string` | **Required.** Statistic (standard or extended). |
| `metric.dimensions` | `map<string,string>` | Dimensions for the metric. |
| `label` | `string` | Display label for the query. |
| `returnData` | `bool` | Set `true` on exactly one query to use its result as the alarm signal. |
| `accountId` | `string` | AWS account ID for cross-account monitoring. |

## Examples

### Simple Metric with SNS Notification

An EC2 CPU alarm that sends to an SNS topic when CPU exceeds 80%:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsCloudwatchAlarm
metadata:
  name: ec2-cpu-alarm
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsCloudwatchAlarm.ec2-cpu-alarm
spec:
  region: us-west-2
  comparisonOperator: GreaterThanThreshold
  evaluationPeriods: 3
  datapointsToAlarm: 2
  threshold: 80.0
  metricName: CPUUtilization
  namespace: AWS/EC2
  period: 300
  statistic: Average
  dimensions:
    InstanceId: i-0abcdef1234567890
  treatMissingData: breaching
  alarmDescription: "EC2 CPU exceeds 80% for 2 of 3 periods"
  alarmActions:
    - value: arn:aws:sns:us-east-1:123456789012:ops-alerts
```

### Error Rate with Metric Math

Computes ALB 5xx error rate as a percentage and alerts when it exceeds 5%:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsCloudwatchAlarm
metadata:
  name: error-rate-alarm
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsCloudwatchAlarm.error-rate-alarm
spec:
  region: us-west-2
  comparisonOperator: GreaterThanThreshold
  evaluationPeriods: 3
  datapointsToAlarm: 2
  threshold: 5.0
  treatMissingData: notBreaching
  alarmDescription: "ALB 5xx error rate exceeds 5%"
  metricQueries:
    - id: errors
      metric:
        metricName: HTTPCode_Target_5XX_Count
        namespace: AWS/ApplicationELB
        period: 300
        stat: Sum
        dimensions:
          LoadBalancer: app/my-alb/1234567890abcdef
    - id: requests
      metric:
        metricName: RequestCount
        namespace: AWS/ApplicationELB
        period: 300
        stat: Sum
        dimensions:
          LoadBalancer: app/my-alb/1234567890abcdef
    - id: error_rate
      expression: "errors/requests*100"
      label: "Error Rate %"
      returnData: true
  alarmActions:
    - value: arn:aws:sns:us-east-1:123456789012:ops-alerts
```

### Production Multi-Action with Foreign Key References

A production SQS depth alarm using `valueFrom` to reference an Planton-managed SNS topic, with actions on all three state transitions:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsCloudwatchAlarm
metadata:
  name: sqs-depth-alarm
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AwsCloudwatchAlarm.sqs-depth-alarm
spec:
  region: us-west-2
  comparisonOperator: GreaterThanOrEqualToThreshold
  evaluationPeriods: 5
  datapointsToAlarm: 3
  threshold: 1000.0
  metricName: ApproximateNumberOfMessagesVisible
  namespace: AWS/SQS
  period: 60
  statistic: Maximum
  treatMissingData: notBreaching
  alarmDescription: "SQS queue depth exceeds 1000 — consumers may be backed up"
  alarmActions:
    - valueFrom:
        kind: AwsSnsTopic
        name: ops-critical
        fieldPath: status.outputs.topic_arn
  okActions:
    - valueFrom:
        kind: AwsSnsTopic
        name: ops-resolved
        fieldPath: status.outputs.topic_arn
  insufficientDataActions:
    - valueFrom:
        kind: AwsSnsTopic
        name: ops-warnings
        fieldPath: status.outputs.topic_arn
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `alarm_arn` | `string` | ARN of the CloudWatch metric alarm. Used by composite alarms and operational tooling. |
| `alarm_name` | `string` | Name of the alarm, unique within the AWS account and region. |

### PromQL Alarm on a Prometheus Workspace

Alarm on a PromQL query with Prometheus-style pending/recovery semantics:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsCloudwatchAlarm
metadata:
  name: promql-cpu-saturation
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: acme
    pulumi.planton.dev/project: platform
    pulumi.planton.dev/stack.name: prod.AwsCloudwatchAlarm.promql-cpu-saturation
spec:
  region: us-west-2
  evaluationCriteria:
    promqlCriteria:
      query: 'avg(rate(node_cpu_seconds_total{mode!="idle"}[5m])) > 0.8'
      pendingPeriod: 300
      recoveryPeriod: 120
  evaluationInterval: 60
  alarmDescription: "Prometheus-reported CPU saturation above 80% for 5 minutes"
  alarmActions:
    - value: arn:aws:sns:us-east-1:123456789012:ops-alerts
```

## Related Components

- [AwsCloudwatchCompositeAlarm](/docs/catalog/aws/awscloudwatchcompositealarm) — combines this alarm's state with others by name for shared-cause paging and maintenance suppression
- [AwsSnsTopic](/docs/catalog/aws/awssnstopic) — the most common alarm action target for notifications
- [AwsCloudwatchLogGroup](/docs/catalog/aws/awscloudwatchloggroup) — log storage that generates metrics (via metric filters) for alarm evaluation
- [AwsSqsQueue](/docs/catalog/aws/awssqsqueue) — queues commonly monitored by depth and age alarms
- [AwsLambda](/docs/catalog/aws/awslambda) — functions monitored by error and duration alarms
