# AwsCloudwatchLogGroup

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

AwsCloudwatchLogGroupSpec defines the desired configuration for an AWS
CloudWatch Logs log group.

A CloudWatch Logs log group is a container for log streams that share the same
retention, monitoring, and access control settings. Log groups are the primary
destination for application logs, service logs, and operational data across
AWS services.

CloudWatch Log Groups are referenced by many AWS services:
- Step Functions (execution logging)
- API Gateway (access logging)
- OpenSearch (slow logs, application logs, audit logs)
- Lambda (function execution logs — auto-created, but can be pre-created for policy control)
- ECS (container logs via awslogs driver)
- Route 53 (public hosted zone query logging)
- WAF (web request logging — one of several destination options)
- ElastiCache (slow logs, engine logs — one of several destination options)

Beyond the container itself, the spec folds in the log-group-scoped
satellites that share the group's lifecycle:
- `metric_filters` — turn matching log events into CloudWatch metrics
  (the raw material for alarms on error counts, latencies parsed from logs, etc.).
- `subscription_filters` — stream matching events in real time to a Kinesis
  stream, Firehose delivery stream, or Lambda function.
- `data_protection_policy` — mask/audit sensitive data (PII) in log events.
- `field_index_policy` — index selected log fields to accelerate and cheapen
  Logs Insights queries.

For the common use case, only `retention_in_days` is needed. KMS encryption
and log group class are optional settings for compliance and cost optimization.

Credentials, region, and deployment workflow live outside this spec in stack inputs.

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsCloudwatchLogGroup
metadata:
  name: test-log-group
  org: test-org
  env: dev
  id: test-log-group-dev
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: test-org
    pulumi.planton.dev/project: test-project
    pulumi.planton.dev/stack.name: dev.AwsCloudwatchLogGroup.test-log-group
spec:
  region: us-west-2
  retentionInDays: 30
  metricFilters:
    - name: error-count
      pattern: "ERROR"
      transformation:
        metricName: ErrorCount
        metricNamespace: TestApp/Errors
        metricValue: "1"
        defaultValue: 0
  fieldIndexPolicy:
    Fields:
      - requestId
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.retentionInDays` | `int32` |  | `30` |  |
| `spec.kmsKeyId` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.logGroupClass` | `string` |  |  |  |
| `spec.deletionProtectionEnabled` | `bool` |  |  |  |
| `spec.metricFilters` | `[]AwsCloudwatchLogGroupMetricFilter` |  |  |  |
| `spec.metricFilters[].name` | `string` | yes |  |  |
| `spec.metricFilters[].pattern` | `string` |  |  |  |
| `spec.metricFilters[].applyOnTransformedLogs` | `bool` |  |  |  |
| `spec.metricFilters[].transformation` | `AwsCloudwatchLogGroupMetricTransformation` | yes |  |  |
| `spec.metricFilters[].transformation.metricName` | `string` | yes |  |  |
| `spec.metricFilters[].transformation.metricNamespace` | `string` | yes |  |  |
| `spec.metricFilters[].transformation.metricValue` | `string` | yes |  |  |
| `spec.metricFilters[].transformation.defaultValue` | `double` |  |  |  |
| `spec.metricFilters[].transformation.dimensions` | `map<string, string>` |  |  |  |
| `spec.metricFilters[].transformation.unit` | `string` |  |  |  |
| `spec.subscriptionFilters` | `[]AwsCloudwatchLogGroupSubscriptionFilter` |  |  |  |
| `spec.subscriptionFilters[].name` | `string` | yes |  |  |
| `spec.subscriptionFilters[].destinationArn` | `string \| valueFrom` | yes |  |  |
| `spec.subscriptionFilters[].filterPattern` | `string` |  |  |  |
| `spec.subscriptionFilters[].roleArn` | `string \| valueFrom` |  |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.subscriptionFilters[].distribution` | `string` |  |  |  |
| `spec.subscriptionFilters[].emitSystemFields` | `[]string` |  |  |  |
| `spec.subscriptionFilters[].applyOnTransformedLogs` | `bool` |  |  |  |
| `spec.dataProtectionPolicy` | `object` |  |  |  |
| `spec.fieldIndexPolicy` | `object` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region where the resource will be created.
Example: "us-west-2", "eu-west-1"

- rule: {"string":{"minLen":"1"}}

### spec.retentionInDays

`int32`

Number of days to retain log events. After this period, log events are
automatically deleted. Set to 0 (the default) to retain log events
indefinitely — they will never expire.

AWS only accepts specific values for this field. Any value not in the
allowed set will be rejected by the AWS API.

Allowed values: 0 (never expire), 1, 3, 5, 7, 14, 30, 60, 90, 120, 150,
180, 365, 400, 545, 731, 1096, 1827, 2192, 2557, 2922, 3288, 3653.

Recommended: Set an explicit retention for cost control. Indefinite retention
(0) accumulates storage costs over time.

- default: `30`

### spec.kmsKeyId

`string | valueFrom`

ARN of the KMS key to use for encrypting log data at rest. When omitted,
CloudWatch Logs uses its default server-side encryption (SSE-CWL).

Customer-managed KMS keys provide:
- Key rotation control
- Cross-account access via key policy
- CloudTrail audit trail of log data access
- Compliance with regulations requiring customer-controlled encryption keys

The KMS key must be in the same region as the log group, and its key policy
must allow the CloudWatch Logs service principal
(`logs.<region>.amazonaws.com`) to use the key. Associating or
disassociating the key updates the log group in place (no replacement).

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.logGroupClass

`string`

The class of the log group, which determines pricing, feature availability,
and data durability characteristics.

Valid values:
- "STANDARD" — Full feature set including metric filters, subscription
  filters, Logs Insights, and Contributor Insights. Default when omitted.
- "INFREQUENT_ACCESS" — Reduced cost (~50% cheaper storage) with a subset
  of features. Supports Logs Insights and managed ingestion to S3, but does
  NOT support metric filters, subscription filters, or Contributor Insights.
  Best for high-volume logs accessed infrequently (VPC flow logs, CDN access
  logs, compliance archives).
- "DELIVERY" — Purpose-built for AWS service log delivery (VPC Flow Logs,
  CloudTrail, Route53 Resolver). Lowest cost. Retention is managed by AWS
  and retention_in_days must not be set.

This field is ForceNew: changing it requires replacing the log group.

### spec.deletionProtectionEnabled

`bool`

When true, the log group is protected from deletion. Any attempt to delete
the log group (including via IaC destroy) will fail until this flag is set
to false. Useful for protecting production log groups from accidental
deletion.

### spec.metricFilters

`[]AwsCloudwatchLogGroupMetricFilter`

Metric filters that extract CloudWatch metrics from log events flowing
through this group. Each filter matches events against a pattern and
publishes a metric value for every match — the standard way to alarm on
"ERROR" counts, parsed latencies, or any signal that lives only in logs.

Filters are keyed by `name` within the group; names must be unique.
Not supported on INFREQUENT_ACCESS log groups (AWS rejects the call).

### spec.metricFilters[].name

`string` · required

Name of the metric filter, unique within the log group. Changing the
name replaces the filter.

- rule: {"required":true,"string":{"maxLen":"512"}}

### spec.metricFilters[].pattern

`string`

Filter pattern that selects which log events produce metric values.
An empty pattern matches ALL log events.

Pattern syntax supports plain terms ("ERROR"), JSON field matching
({ $.statusCode = 500 }), and space-delimited column matching
([ip, user, ts, request, status=5*, size]).
Maximum 1024 characters.

- rule: {"string":{"maxLen":"1024"}}

### spec.metricFilters[].applyOnTransformedLogs

`bool`

When true, this filter is also applied to log events AFTER they pass
through the log group's transformer (if one is configured). Leave unset
for the standard behavior of filtering raw ingested events.

### spec.metricFilters[].transformation

`AwsCloudwatchLogGroupMetricTransformation` · required

The metric to publish for each matching log event. Required — a filter
without a transformation does nothing.

- rule: {"required":true}
- rule: default_value cannot be set together with dimensions — AWS does not support a default value on dimensional metric filters
- rule: a metric filter supports at most 3 dimensions
- rule: unit must be a valid CloudWatch unit (e.g. Seconds, Milliseconds, Bytes, Percent, Count, Bytes/Second, Count/Second, None) when set

### spec.metricFilters[].transformation.metricName

`string` · required

Name of the CloudWatch metric to publish (e.g. "ErrorCount").

- rule: {"required":true,"string":{"maxLen":"255"}}

### spec.metricFilters[].transformation.metricNamespace

`string` · required

Namespace for the metric (e.g. "MyApp/Errors"). Custom namespaces keep
log-derived metrics separate from AWS service namespaces.

- rule: {"required":true,"string":{"maxLen":"255"}}

### spec.metricFilters[].transformation.metricValue

`string` · required

Value to publish for each matching event. Either a literal number
("1" to count occurrences) or a field reference that extracts a numeric
value from the matched event (e.g. "$.latencyMs" for JSON patterns, or
"$size" for a named column in space-delimited patterns).
Maximum 100 characters.

- rule: {"required":true,"string":{"maxLen":"100"}}

### spec.metricFilters[].transformation.defaultValue

`double` · optional (explicit presence)

Value to publish for periods when NO log events match the pattern.
Typically "0" so count metrics report zero instead of missing data —
which lets alarms use standard missing-data handling. When unset, no
value is published for non-matching periods.

AWS does not allow a default value on filters that publish dimensions.

### spec.metricFilters[].transformation.dimensions

`map<string, string>`

Dimensions to publish with the metric, mapping dimension names to field
references in the matched event (e.g. {"ErrorCode": "$.errorCode"}).
Maximum 3 dimensions per AWS limit.

Every unique dimension-value combination creates a distinct custom
metric (billed separately) — keep dimension cardinality low.

### spec.metricFilters[].transformation.unit

`string`

Unit for the metric (e.g. "Count", "Milliseconds", "Bytes"). Defaults to
"None" when unset. Must be a valid CloudWatch StandardUnit.

### spec.subscriptionFilters

`[]AwsCloudwatchLogGroupSubscriptionFilter`

Real-time subscription filters that stream matching log events to a
destination: a Kinesis data stream, a Kinesis Data Firehose delivery
stream, or a Lambda function. Use these to fan logs out to analytics
pipelines, SIEM tooling, or custom processing.

AWS allows at most TWO subscription filters per log group. Names must be
unique within the group. Not supported on INFREQUENT_ACCESS log groups.

- rule: distribution must be 'ByLogStream' or 'Random' when set
- rule: emit_system_fields entries must be '@aws.account' or '@aws.region'

### spec.subscriptionFilters[].name

`string` · required

Name of the subscription filter, unique within the log group. Changing
the name replaces the filter.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"512"}}

### spec.subscriptionFilters[].destinationArn

`string | valueFrom` · required

ARN of the destination that receives matching log events. Supported
destinations:
- Kinesis data stream (same-account) — reference an AwsKinesisStream
- Kinesis Data Firehose delivery stream — reference an AwsKinesisFirehose
- Lambda function — reference an AwsLambda
- Cross-account CloudWatch Logs destination (by literal ARN)

No default kind is set because all destination types are equally common —
reference the specific resource kind's ARN output, or provide a literal ARN.

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.subscriptionFilters[].filterPattern

`string`

Filter pattern that selects which log events are delivered. An empty
pattern delivers ALL log events. Same syntax as metric filter patterns.
Maximum 1024 characters.

- rule: {"string":{"maxLen":"1024"}}

### spec.subscriptionFilters[].roleArn

`string | valueFrom`

ARN of the IAM role that grants CloudWatch Logs permission to put records
to the destination. REQUIRED for Kinesis stream and Firehose destinations
(the role must trust `logs.amazonaws.com` and allow `kinesis:PutRecord` /
`firehose:PutRecord` on the destination). NOT used for Lambda destinations —
those authorize via a Lambda resource-based permission instead.

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.subscriptionFilters[].distribution

`string`

How log data is distributed to Kinesis stream destinations:
- "ByLogStream" — (default) events from the same log stream go to the
  same shard, preserving per-stream ordering.
- "Random" — events are spread across shards for maximum throughput,
  at the cost of ordering.

Only meaningful for Kinesis stream destinations.

### spec.subscriptionFilters[].emitSystemFields

`[]string`

System fields to include with each delivered log event. Supported values:
"@aws.account" (the source account ID) and "@aws.region" (the source
region). Useful when a central destination aggregates logs from many
accounts/regions.

### spec.subscriptionFilters[].applyOnTransformedLogs

`bool`

When true, this filter is also applied to log events AFTER they pass
through the log group's transformer (if one is configured).

### spec.dataProtectionPolicy

`object`

CloudWatch Logs data protection policy for this log group, as a JSON
policy document. Data protection audits and masks sensitive data (PII
such as email addresses, credit card numbers, or custom identifiers) in
log events as they are ingested.

The document follows the CloudWatch Logs data protection policy schema:
a `Name`, `Version: "2021-06-01"`, and a `Statement` list carrying one
"Audit" operation statement and one "Deidentify" (masking) statement over
the same data identifiers. Audit findings destinations (CloudWatch Logs,
S3, Firehose) are optional — masking works without them.

One policy per log group. Masked data is visible only to principals with
the `logs:Unmask` permission.

### spec.fieldIndexPolicy

`object`

CloudWatch Logs field index policy for this log group, as a JSON policy
document with a `Fields` array of log-event field names to index
(e.g. {"Fields": ["requestId", "userId"]}).

Field indexes make Logs Insights queries that filter on the indexed
fields faster and cheaper: matching scans skip events where the indexed
value cannot match. Index up to 20 fields per policy. One policy per log
group; an account-level policy (managed outside this resource) can also
apply — the two merge at query time.

## Validation Rules

- `retention_in_days_valid_values`: retention_in_days must be one of: 0, 1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1096, 1827, 2192, 2557, 2922, 3288, 3653
- `log_group_class_valid_values`: log_group_class must be 'STANDARD', 'INFREQUENT_ACCESS', or 'DELIVERY' when set
- `delivery_class_no_retention`: retention_in_days must not be set (must be 0) when log_group_class is 'DELIVERY' — AWS manages retention for Delivery log groups
- `subscription_filters_max_2`: a log group supports at most 2 subscription filters — AWS enforces this limit per log group
- `metric_filter_names_unique`: each metric filter must have a unique name within the log group
- `subscription_filter_names_unique`: each subscription filter must have a unique name within the log group
- `infrequent_access_no_filters`: metric_filters and subscription_filters are not supported on INFREQUENT_ACCESS log groups — use a STANDARD class log group for filtering features

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsCloudwatchLogGroup, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.log_group_arn` | `string` | The Amazon Resource Name (ARN) of the log group. This is the primary identifier used to reference the log group in other AWS resources such as Step Functions logging, API Gateway access logs, and OpenSearch log publishing. |
| `status.outputs.log_group_name` | `string` | The name of the log group. Some AWS services (e.g., ElastiCache log delivery, ECS awslogs driver) reference log groups by name rather than ARN. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.kmsKeyId` | AwsKmsKey | `status.outputs.key_arn` |
| `spec.subscriptionFilters[].roleArn` | AwsIamRole | `status.outputs.role_arn` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AwsClientVpn | `spec.connectionLog.cloudwatchLogGroup` | `status.outputs.log_group_name` |
| AwsCodeBuildProject | `spec.logsConfig.cloudwatchLogs.groupName` | `status.outputs.log_group_name` |
| AwsCognitoUserPool | `spec.logConfigurations[].cloudwatchLogGroupArn` | `status.outputs.log_group_arn` |
| AwsEcsTaskDefinition | `spec.logging.logGroup` | `status.outputs.log_group_name` |
| AwsFsxLustreFileSystem | `spec.logConfiguration.destination` | `status.outputs.log_group_arn` |
| AwsFsxWindowsFileSystem | `spec.auditLogConfiguration.auditLogDestination` | `status.outputs.log_group_arn` |
| AwsHttpApiGateway | `spec.stage.accessLog.destinationArn` | `status.outputs.log_group_arn` |
| AwsLambda | `spec.loggingConfig.logGroup` | `status.outputs.log_group_name` |
| AwsMskCluster | `spec.logging.cloudwatchLogs.logGroup` | `status.outputs.log_group_name` |
| AwsOpenSearchDomain | `spec.logPublishingOptions[].cloudwatchLogGroupArn` | `status.outputs.log_group_arn` |
| AwsRoute53Zone | `spec.queryLogging.cloudwatchLogGroupArn` | `status.outputs.log_group_arn` |
| AwsStepFunction | `spec.logging.logDestination` | `status.outputs.log_group_arn` |

## See Also

- [Overview](../README.md)
