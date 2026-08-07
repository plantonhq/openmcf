# AwsSnsSubscription

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

AwsSnsSubscriptionSpec defines a subscription that delivers messages from an
SNS topic to a target endpoint — an SQS queue, a Lambda function, an HTTP/S
endpoint, an email address, an SMS number, a Kinesis Data Firehose stream,
or a mobile platform application endpoint.

The subscription is its own node in the resource graph, not a setting of the
topic: it has independent AWS identity (a subscription ARN), a topic may
have many subscriptions, and each subscription owns its own delivery
lifecycle — filtering, raw delivery, dead-lettering, archived-message
replay, and (for HTTP/S and email) an asynchronous confirmation handshake.
This shape also supports subscribing an endpoint you own to a topic owned
by another team or account: reference the topic by literal ARN.

Delivery couplings to know about:
- "sqs" delivery requires the QUEUE's own resource policy to grant
  sqs:SendMessage to sns.amazonaws.com (scoped by aws:SourceArn to this
  topic). Creating the subscription succeeds without it, but every delivery
  is silently dropped. Set the policy on the AwsSqsQueue resource.
- FIFO topics can only deliver to SQS (standard or FIFO) endpoints.
- "lambda" delivery additionally needs a Lambda invoke permission for
  sns.amazonaws.com (the AwsLambda kind's invoke_permissions fold).

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsSnsSubscription
metadata:
  name: test-sqs-subscription
spec:
  region: us-west-2
  topicArn:
    value: "<sns-topic-arn>"
  protocol: sqs
  endpoint:
    value: "<sqs-queue-arn>"
  rawMessageDelivery: true
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.topicArn` | `string \| valueFrom` | yes |  | AwsSnsTopic (`status.outputs.topic_arn`) |
| `spec.protocol` | `string` | yes |  |  |
| `spec.endpoint` | `string \| valueFrom` | yes |  |  |
| `spec.filterPolicy` | `object` |  |  |  |
| `spec.filterPolicyScope` | `string` |  |  |  |
| `spec.rawMessageDelivery` | `bool` |  |  |  |
| `spec.deadLetterConfig` | `AwsSnsSubscriptionDeadLetterConfig` |  |  |  |
| `spec.deadLetterConfig.deadLetterTargetArn` | `string \| valueFrom` | yes |  | AwsSqsQueue (`status.outputs.queue_arn`) |
| `spec.deliveryPolicy` | `string` |  |  |  |
| `spec.replayPolicy` | `object` |  |  |  |
| `spec.subscriptionRoleArn` | `string \| valueFrom` |  |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.endpointAutoConfirms` | `bool` |  |  |  |
| `spec.confirmationTimeoutMinutes` | `int32` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region of the subscription — must be the topic's region.
Example: "us-west-2", "eu-west-1"

- rule: {"string":{"minLen":"1"}}

### spec.topicArn

`string | valueFrom` · required

The SNS topic to subscribe to. Accepts a reference to an AwsSnsTopic
resource or a literal topic ARN (including a cross-account topic that has
granted this account sns:Subscribe). Create-time immutable — a different
topic is a different subscription.

- references: AwsSnsTopic (`status.outputs.topic_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSnsTopic, name: <that resource's name>, fieldPath: status.outputs.topic_arn}} -- a bare string does not parse

### spec.protocol

`string` · required

Protocol for message delivery. Determines how SNS delivers messages and
what `endpoint` must contain. Create-time immutable.
- "sqs": SQS queue ARN
- "lambda": Lambda function ARN
- "http" / "https": URL endpoint (requires endpoint-side confirmation)
- "email" / "email-json": email address (always requires manual confirmation)
- "sms": phone number in E.164 format
- "firehose": Kinesis Data Firehose delivery stream ARN (requires
  subscription_role_arn)
- "application": mobile platform endpoint ARN

- rule: {"required":true}

### spec.endpoint

`string | valueFrom` · required

Endpoint to deliver messages to. The format depends on the protocol (see
`protocol` field documentation). Accepts a direct value or a reference to
another resource's output — e.g. an SQS subscription references an
AwsSqsQueue's queue_arn, a Lambda subscription references an AwsLambda's
function_arn. Create-time immutable.

No `default_kind` is set because the target resource varies by protocol
(SQS queue, Lambda function, Firehose stream, plain URL, ...).

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.filterPolicy

`object`

Filter policy to select which messages this subscription receives.
Messages that match no filter are not delivered. The filter can be applied
to message attributes (default) or the message body (see
`filter_policy_scope`). Expressed as a JSON structure in YAML.

### spec.filterPolicyScope

`string`

Scope for the filter policy. Controls whether the filter is evaluated
against message attributes or the message body.
Valid values: "MessageAttributes" (default), "MessageBody".
Only relevant when `filter_policy` is set.

### spec.rawMessageDelivery

`bool`

When true, messages are delivered as-is without JSON wrapping. Supported
for SQS, HTTP/S, and Firehose protocols. When false (default), SNS wraps
the message in a JSON envelope containing metadata (MessageId, TopicArn,
Timestamp, etc.).

### spec.deadLetterConfig

`AwsSnsSubscriptionDeadLetterConfig`

Dead letter queue for this subscription's delivery failures. When SNS
cannot deliver a message after all retry attempts, the message is routed
to the specified SQS queue. This is separate from any DLQ the endpoint
itself has — it catches SNS-to-subscriber delivery failures.

### spec.deadLetterConfig.deadLetterTargetArn

`string | valueFrom` · required

ARN of the SQS dead letter queue for failed message deliveries. Accepts a
direct ARN or a reference to an AwsSqsQueue resource. The queue must
reside in the same AWS account and region as the SNS topic, and its
resource policy must allow sns.amazonaws.com to SendMessage.

- references: AwsSqsQueue (`status.outputs.queue_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSqsQueue, name: <that resource's name>, fieldPath: status.outputs.queue_arn}} -- a bare string does not parse

### spec.deliveryPolicy

`string`

HTTP/HTTPS delivery retry policy override for this subscription. Expressed
as a JSON string matching the SNS delivery policy format. Overrides the
topic-level delivery policy for this subscription only. Most users do not
need to customize this.

### spec.replayPolicy

`object`

Replay policy for archived messages. When the topic (FIFO with
`archive_policy`) retains messages, a new subscription can replay the
archive from a starting point before receiving live traffic — the
mechanism for backfilling a new consumer. Expressed as the SNS replay
policy JSON document, e.g. {"PointType": "Timestamp",
"StartingPoint": "2026-07-01T00:00:00Z"}.

### spec.subscriptionRoleArn

`string | valueFrom`

IAM role ARN for Firehose delivery. Required when protocol is "firehose".
The role must grant SNS permission to write to the Firehose delivery
stream. Accepts a direct ARN or a reference to an AwsIamRole resource.

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.endpointAutoConfirms

`bool`

Set to true when the HTTP/S endpoint confirms subscriptions on its own
(responds to the SubscriptionConfirmation callback without a human).
Deployment then waits for the confirmation to complete instead of leaving
the subscription pending. Only meaningful for "http"/"https" protocols —
SQS, Lambda, Firehose, and application subscriptions confirm automatically,
and email subscriptions always require a manual click.

### spec.confirmationTimeoutMinutes

`int32`

Minutes to wait for an HTTP/S endpoint to confirm the subscription before
deployment fails. Leave at 0 for the AWS provider default (1 minute).
Only meaningful for "http"/"https" protocols.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"gte":1}}

## Validation Rules

- `protocol_valid`: protocol must be one of: sqs, lambda, http, https, email, email-json, sms, firehose, application
- `filter_policy_scope_valid`: filter_policy_scope must be 'MessageAttributes' or 'MessageBody' when set
- `filter_policy_scope_requires_filter_policy`: filter_policy_scope requires filter_policy to be set
- `subscription_role_arn_required_for_firehose`: subscription_role_arn is required when protocol is 'firehose'
- `auto_confirms_http_only`: endpoint_auto_confirms is only meaningful for 'http' or 'https' protocols
- `confirmation_timeout_http_only`: confirmation_timeout_minutes is only meaningful for 'http' or 'https' protocols

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsSnsSubscription, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.subscription_arn` | `string` | The Amazon Resource Name (ARN) of the subscription. The identifier used for sns:Unsubscribe/SetSubscriptionAttributes permissions and for auditing which subscriptions exist on a topic. |
| `status.outputs.owner_id` | `string` | The AWS account ID that owns the subscription. Useful when subscribing an owned endpoint to a cross-account topic. |
| `status.outputs.pending_confirmation` | `bool` | True while the subscription is awaiting endpoint confirmation (HTTP/S endpoints that have not completed the handshake, or email recipients who have not clicked the confirmation link). Deliveries do not start until confirmation completes. |
| `status.outputs.confirmation_was_authenticated` | `bool` | True when the confirmation request was authenticated (signed). Always true for protocols that confirm automatically (SQS, Lambda, Firehose). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.topicArn` | AwsSnsTopic | `status.outputs.topic_arn` |
| `spec.deadLetterConfig.deadLetterTargetArn` | AwsSqsQueue | `status.outputs.queue_arn` |
| `spec.subscriptionRoleArn` | AwsIamRole | `status.outputs.role_arn` |

## See Also

- [Overview](../README.md)
