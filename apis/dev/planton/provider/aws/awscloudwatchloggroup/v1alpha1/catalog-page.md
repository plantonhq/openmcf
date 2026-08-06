# AWS CloudWatch Log Group

Deploys an AWS CloudWatch Logs log group with configurable retention policy, optional KMS encryption, log group class selection, deletion protection, and the group-scoped satellites: metric filters (log events → CloudWatch metrics), subscription filters (real-time delivery to Kinesis, Firehose, or Lambda), a data protection policy (PII audit + masking), and a field index policy (faster Logs Insights queries). The log group serves as a centralized destination for application logs, service logs, and operational data, and is referenced by many other AWS components including Step Functions, API Gateway, and OpenSearch.

## What Gets Created

- **CloudWatch Log Group** — a container for log streams with the specified retention, encryption, class, and deletion-protection settings
- **Metric filters** (optional) — one per named filter, publishing metrics extracted from matching log events
- **Subscription filters** (optional, max 2) — real-time delivery of matching events to a Kinesis stream, Firehose delivery stream, or Lambda function
- **Data protection policy** (optional) — audits and masks sensitive data (PII) in log events at ingestion
- **Field index policy** (optional) — indexes selected log fields to accelerate and cheapen Logs Insights queries

## Prerequisites

- An AWS account with credentials configured in the stack input
- An AwsKmsKey resource if enabling customer-managed encryption

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsCloudwatchLogGroup
metadata:
  name: app-logs
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsCloudwatchLogGroup.app-logs
spec:
  region: us-west-2
  retentionInDays: 30
```

```shell
planton apply -f log-group.yaml
```

This creates a STANDARD class log group with 30-day retention using default AWS encryption.

## Configuration Reference

### Required Fields

No fields are strictly required. An empty spec creates a STANDARD class log group with indefinite retention and default encryption.

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `retentionInDays` | int | 0 (never expire) | Days to retain log events. Must be one of: 0, 1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1096, 1827, 2192, 2557, 2922, 3288, 3653. Recommended default: 30. |
| `kmsKeyId` | StringValueOrRef | — | KMS key ARN for encrypting log data at rest. Can reference an AwsKmsKey resource via `valueFrom`. The key policy must allow `logs.<region>.amazonaws.com`. |
| `logGroupClass` | string | STANDARD | Log group class. Valid values: `STANDARD`, `INFREQUENT_ACCESS`, `DELIVERY`. ForceNew — changing requires replacing the log group. |
| `deletionProtectionEnabled` | bool | false | When true, every delete (including IaC destroy) fails until the flag is cleared. |
| `metricFilters[]` | list | — | Metric filters publishing CloudWatch metrics from matching log events. Each has a unique `name`, a `pattern` (empty = all events), and a `transformation` (metric name, namespace, value, optional default value / dimensions / unit). |
| `subscriptionFilters[]` | list | — | Real-time subscriptions delivering matching events to a destination ARN (Kinesis stream, Firehose, or Lambda), with an IAM `roleArn` for Kinesis/Firehose delivery, `distribution` mode, and optional `@aws.account`/`@aws.region` enrichment. Maximum 2 per group. |
| `dataProtectionPolicy` | object | — | CloudWatch Logs data protection policy document (audit + deidentify statements over data identifiers). |
| `fieldIndexPolicy` | object | — | Field index policy document, e.g. `{"Fields": ["requestId"]}` (up to 20 fields). |

**Validation rules:**
- `retentionInDays` must be one of the AWS-allowed discrete values listed above
- `logGroupClass` must be `STANDARD`, `INFREQUENT_ACCESS`, or `DELIVERY` when set
- `retentionInDays` must not be set when `logGroupClass` is `DELIVERY` (AWS manages retention for Delivery log groups)
- Metric filter and subscription filter names must be unique within the group; at most 2 subscription filters
- Metric/subscription filters are rejected on `INFREQUENT_ACCESS` groups (AWS does not support them there)
- A metric transformation cannot combine `defaultValue` with `dimensions`, and supports at most 3 dimensions

## Examples

### Standard 30-Day Retention

A general-purpose log group for application logging:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsCloudwatchLogGroup
metadata:
  name: app-logs
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: acme
    pulumi.planton.dev/project: platform
    pulumi.planton.dev/stack.name: dev.AwsCloudwatchLogGroup.app-logs
spec:
  region: us-west-2
  retentionInDays: 30
```

### Encrypted Production Log Group

A log group with 90-day retention and KMS encryption for compliance workloads:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsCloudwatchLogGroup
metadata:
  name: prod-app-logs
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: acme
    pulumi.planton.dev/project: platform
    pulumi.planton.dev/stack.name: prod.AwsCloudwatchLogGroup.prod-app-logs
spec:
  region: us-west-2
  retentionInDays: 90
  kmsKeyId:
    valueFrom:
      kind: AwsKmsKey
      name: log-encryption-key
      fieldPath: status.outputs.key_arn
```

### Error-Count Metric Filter (Log-Derived Alerting)

A log group whose metric filter counts `ERROR` events, ready to alarm on with an AWS CloudWatch Alarm:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsCloudwatchLogGroup
metadata:
  name: app-logs
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: acme
    pulumi.planton.dev/project: platform
    pulumi.planton.dev/stack.name: prod.AwsCloudwatchLogGroup.app-logs
spec:
  region: us-west-2
  retentionInDays: 30
  metricFilters:
    - name: error-count
      pattern: "ERROR"
      transformation:
        metricName: ErrorCount
        metricNamespace: MyApp/Errors
        metricValue: "1"
        defaultValue: 0
```

### Infrequent Access for High-Volume Logs

A cost-optimized log group for VPC flow logs or CDN access logs:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsCloudwatchLogGroup
metadata:
  name: vpc-flow-logs
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: acme
    pulumi.planton.dev/project: networking
    pulumi.planton.dev/stack.name: prod.AwsCloudwatchLogGroup.vpc-flow-logs
spec:
  region: us-west-2
  retentionInDays: 365
  logGroupClass: INFREQUENT_ACCESS
  kmsKeyId:
    valueFrom:
      kind: AwsKmsKey
      name: infra-key
      fieldPath: status.outputs.key_arn
```

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `log_group_arn` | string | The ARN of the CloudWatch Log Group. Used by downstream resources (Step Functions, API Gateway, OpenSearch) via `valueFrom`. |
| `log_group_name` | string | The name of the CloudWatch Log Group. Used by services that reference log groups by name (ElastiCache, ECS). |

## Related Components

- [AWS CloudWatch Alarm](/docs/catalog/aws/awscloudwatchalarm) — Alarms on metrics published by this group's metric filters
- [AWS KMS Key](/docs/catalog/aws/awskmskey) — Customer-managed encryption key for log data
- [AWS Kinesis Stream](/docs/catalog/aws/awskinesisstream) — Real-time subscription filter destination
- [AWS IAM Role](/docs/catalog/aws/awsiamrole) — Delivery role for Kinesis/Firehose subscription filters
- [AWS Step Function](/docs/catalog/aws/awsstepfunction) — Uses log group ARN for execution logging
- [AWS HTTP API Gateway](/docs/catalog/aws/awshttpapigateway) — Uses log group ARN for access logging
- [AWS OpenSearch Domain](/docs/catalog/aws/awsopensearchdomain) — Uses log group ARN for slow logs, app logs, and audit logs
- [AWS Route53 Zone](/docs/catalog/aws/awsroute53zone) — Uses log group ARN for public-zone query logging
