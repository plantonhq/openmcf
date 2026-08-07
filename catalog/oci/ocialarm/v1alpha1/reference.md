# OciAlarm

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciAlarmSpec defines the specification for an OCI Monitoring Alarm --
a rule that evaluates metrics via Monitoring Query Language (MQL)
expressions and triggers notifications to ONS topics or Streaming
endpoints when thresholds are breached.

Key behaviors:
  - All fields are updatable after creation (no ForceNew attributes)
  - Severity is required; alarms without severity cannot fire
  - Destinations must contain at least one ONS topic or Stream OCID
  - Overrides enable multi-threshold evaluation (e.g., warn at 80%,
    critical at 95%) with per-override query/severity/body/duration

Excluded from v1:
  - suppression -- requires hardcoded RFC3339 timestamps, an
    anti-pattern in declarative IaC (timestamps go stale); manage
    via OCI Console or CLI
  - resolution -- only supported value is "1m", no user choice
  - defined_tags, system_tags -- managed by platform
  - freeform_tags -- auto-populated from metadata labels

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.metricCompartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.namespace` | `string` | yes |  |  |
| `spec.query` | `string` | yes |  |  |
| `spec.severity` | `enum` |  |  |  |
| `spec.destinations` | `[]string` | yes |  |  |
| `spec.isEnabled` | `bool` |  |  |  |
| `spec.body` | `string` |  |  |  |
| `spec.alarmSummary` | `string` |  |  |  |
| `spec.notificationTitle` | `string` |  |  |  |
| `spec.pendingDuration` | `string` |  |  |  |
| `spec.evaluationSlackDuration` | `string` |  |  |  |
| `spec.repeatNotificationDuration` | `string` |  |  |  |
| `spec.messageFormat` | `enum` |  |  |  |
| `spec.metricCompartmentIdInSubtree` | `bool` |  |  |  |
| `spec.isNotificationsPerMetricDimensionEnabled` | `bool` |  |  |  |
| `spec.resourceGroup` | `string` |  |  |  |
| `spec.notificationVersion` | `string` |  |  |  |
| `spec.ruleName` | `string` |  |  |  |
| `spec.overrides` | `[]AlarmOverride` |  |  |  |
| `spec.overrides[].ruleName` | `string` | yes |  |  |
| `spec.overrides[].query` | `string` |  |  |  |
| `spec.overrides[].severity` | `enum` |  |  |  |
| `spec.overrides[].body` | `string` |  |  |  |
| `spec.overrides[].pendingDuration` | `string` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the alarm will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.metricCompartmentId

`string | valueFrom` · required

OCID of the compartment containing the metric being evaluated.
Often the same as compartment_id, but can differ when monitoring
metrics from another compartment.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.namespace

`string` · required

Source service or application emitting the metric.
Examples: "oci_computeagent", "oci_blockstore",
"oci_autonomous_database", "oci_vcn".

- rule: {"string":{"minLen":"1"}}

### spec.query

`string` · required

Monitoring Query Language (MQL) expression to evaluate.
Must specify metric, statistic, interval, and trigger rule.
Example: "CpuUtilization[5m].mean() > 80"

- rule: {"string":{"minLen":"1"}}

### spec.severity

`enum`

Perceived severity when the alarm is in the FIRING state.
Must be explicitly set (unspecified is rejected).

- rule: severity must be explicitly set (critical, error, warning, or info)

Allowed values (use exactly as shown):

- `unspecified`
- `critical`
- `error`
- `warning`
- `info`

### spec.destinations

`[]string` · required

OCIDs of notification destinations. Each OCID must reference an
ONS Notification Topic or a Streaming stream. At least one
destination is required.

- rule: {"repeated":{"minItems":"1"}}

### spec.isEnabled

`bool`

Whether the alarm is enabled. When false, the alarm does not
evaluate metrics or send notifications.
Proto3 default is false -- alarms start disabled unless
explicitly set to true.

### spec.body

`string`

Human-readable content of the delivered alarm notification.
Supports dynamic variables: {{severity}}, {{query}},
{{metricValues}}, {{resourceId}}, {{timestamp}}.

### spec.alarmSummary

`string`

Customizable alarm summary. Appears in the alarm message body
and API responses. Supports the same dynamic variables as body.

### spec.notificationTitle

`string`

Custom notification title. Appears as the email subject line
or Slack notification title. Supports dynamic variables.

### spec.pendingDuration

`string`

Period of time the condition must persist before the alarm
transitions from OK to FIRING. ISO 8601 duration format.
Minimum: PT1M. Maximum: PT1H. OCI default: PT1M.

### spec.evaluationSlackDuration

`string`

Slack period to wait for metric ingestion before evaluating
the alarm. Accounts for delayed metric delivery. ISO 8601
duration format. Minimum: PT3M. Maximum: PT2H. OCI default: PT3M.

### spec.repeatNotificationDuration

`string`

Frequency for re-submitting alarm notifications while the
alarm remains in FIRING state. ISO 8601 duration format.
Minimum: PT1M. Maximum: P30D. When omitted, notifications
are not re-submitted.

### spec.messageFormat

`enum`

Format for alarm notification delivery.
When omitted (zero-value), OCI applies RAW format.

Allowed values (use exactly as shown):

- `raw` -- raw is the default OCI format. Works with both Notifications and Streaming destinations.
- `pretty_json` -- pretty_json delivers a human-readable JSON payload. Only works with Notifications destinations.
- `ons_optimized` -- ons_optimized delivers a compact payload optimized for email. Only works with Notifications destinations.

### spec.metricCompartmentIdInSubtree

`bool` · optional (explicit presence)

When true, evaluates metrics from the specified compartment
and all of its sub-compartments. Can only be true when
metric_compartment_id is a tenancy OCID.

### spec.isNotificationsPerMetricDimensionEnabled

`bool` · optional (explicit presence)

When true, splits alarm notifications per metric stream
dimension. When false (default), groups notifications across
all metric streams matching the query.

### spec.resourceGroup

`string`

Resource group to match when filtering metric data.
When omitted, only metric data with no resource group is
returned.

### spec.notificationVersion

`string`

Version of the alarm notification to be delivered.
Format: a number (up to 4 digits) followed by a period
and uppercase X (e.g., "1.X").

### spec.ruleName

`string`

Identifier for the alarm's base values for evaluation.
Default value is "BASE". Only meaningful when overrides
are configured. Must be unique across all rule_name values
(including those in overrides).

### spec.overrides

`[]AlarmOverride`

Overrides controlling alarm evaluation at different thresholds.
Each override can specify its own query, severity, body, and
pending duration. Evaluated in order before the base rule.

### spec.overrides[].ruleName

`string` · required

User-friendly identifier for this override. Must be unique
across all rule_name values for the alarm.

- rule: {"string":{"minLen":"1"}}

### spec.overrides[].query

`string`

Override MQL query. When omitted, the base query is used.

### spec.overrides[].severity

`enum`

Override severity. When unspecified (zero-value), the base
severity is used.

Allowed values (use exactly as shown):

- `unspecified`
- `critical`
- `error`
- `warning`
- `info`

### spec.overrides[].body

`string`

Override notification body. When omitted, the base body
is used. Supports the same dynamic variables.

### spec.overrides[].pendingDuration

`string`

Override pending duration. When omitted, the base pending
duration is used. ISO 8601 format.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciAlarm, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.alarm_id` | `string` | OCID of the alarm. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |
| `spec.metricCompartmentId` | OciCompartment | `status.outputs.compartment_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
