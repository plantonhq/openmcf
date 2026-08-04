# OciQueue

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1`

OciQueueSpec defines the specification for an OCI Queue -- a fully
managed, serverless message queue for asynchronous communication
between decoupled services.

A queue provides at-least-once delivery with configurable visibility
timeouts, dead-letter queue support, optional KMS encryption, and
two additive capabilities: large message support (up to 512 KB) and
consumer groups for partitioned consumption.

Key behaviors:
  - retention_in_seconds is ForceNew (changing forces queue recreation)
  - All other fields are updatable after creation
  - Capabilities are additive and modeled as flattened fields
    (is_large_messages_enabled, consumer_group_config) rather than
    the provider's discriminated list format

Excluded from v1:
  - purge_trigger / purge_type -- imperative operational controls,
    not declarative infrastructure
  - primary_consumer_group_filter -- always empty for primary group
  - defined_tags, system_tags -- managed by platform
  - freeform_tags -- auto-populated from metadata labels

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.customEncryptionKeyId` | `string \| valueFrom` |  |  | OciKmsKey (`status.outputs.key_id`) |
| `spec.deadLetterQueueDeliveryCount` | `int32` |  |  |  |
| `spec.retentionInSeconds` | `int32` |  |  |  |
| `spec.timeoutInSeconds` | `int32` |  |  |  |
| `spec.visibilityInSeconds` | `int32` |  |  |  |
| `spec.channelConsumptionLimit` | `int32` |  |  |  |
| `spec.isLargeMessagesEnabled` | `bool` |  |  |  |
| `spec.consumerGroupConfig` | `ConsumerGroupConfig` |  |  |  |
| `spec.consumerGroupConfig.isPrimaryEnabled` | `bool` |  |  |  |
| `spec.consumerGroupConfig.primaryDeadLetterQueueDeliveryCount` | `int32` |  |  |  |
| `spec.consumerGroupConfig.primaryDisplayName` | `string` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the queue will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.customEncryptionKeyId

`string | valueFrom`

OCID of the custom encryption key for encrypting message content.
When omitted, Oracle-managed encryption is used.

- references: OciKmsKey (`status.outputs.key_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_id}} -- a bare string does not parse

### spec.deadLetterQueueDeliveryCount

`int32` · optional (explicit presence)

Number of delivery attempts before a message is moved to the dead
letter queue. A value of 0 disables the DLQ. When omitted, OCI
applies its default.

### spec.retentionInSeconds

`int32` · optional (explicit presence)

Retention period for messages in seconds. ForceNew -- changing this
value forces queue recreation. When omitted, OCI applies its
default (604800 seconds / 7 days).

### spec.timeoutInSeconds

`int32` · optional (explicit presence)

Default polling timeout for GetMessages calls, in seconds.
When omitted, OCI applies its default.

### spec.visibilityInSeconds

`int32` · optional (explicit presence)

Default visibility timeout for consumed messages, in seconds.
A consumed message is invisible to other consumers for this
duration. When omitted, OCI applies its default.

### spec.channelConsumptionLimit

`int32` · optional (explicit presence)

Percentage of allocated queue resources that can be consumed by a
single channel. Both Terraform and Pulumi providers type this as
integer; consult OCI documentation for the valid range and scale.
When omitted, OCI applies its default (no per-channel limit).

### spec.isLargeMessagesEnabled

`bool` · optional (explicit presence)

Enable large message support (up to 512 KB per message).
When omitted or false, standard message size limits apply.
Maps to the LARGE_MESSAGES capability in the OCI API.

### spec.consumerGroupConfig

`ConsumerGroupConfig`

Consumer group configuration. When provided, the CONSUMER_GROUPS
capability is enabled for this queue, allowing partitioned
consumption patterns.

### spec.consumerGroupConfig.isPrimaryEnabled

`bool` · optional (explicit presence)

Whether to automatically enable the primary consumer group
after adding the capability.

### spec.consumerGroupConfig.primaryDeadLetterQueueDeliveryCount

`int32` · optional (explicit presence)

Dead letter queue delivery count for the primary consumer group.
Overrides the queue-level dead_letter_queue_delivery_count.
A value of 0 disables the DLQ for this consumer group.

### spec.consumerGroupConfig.primaryDisplayName

`string`

Display name for the primary consumer group.
When omitted, OCI defaults to "Primary Consumer Group".

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciQueue, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.queue_id` | `string` | OCID of the queue. |
| `status.outputs.messages_endpoint` | `string` | Endpoint URL for consuming or publishing messages in the queue. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |
| `spec.customEncryptionKeyId` | OciKmsKey | `status.outputs.key_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
