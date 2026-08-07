# AzureServiceBusTopic

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureServiceBusTopicSpec** defines the configuration for creating a
topic inside an Azure Service Bus namespace: the publish-subscribe
primitive. Publishers send to the topic; each
AzureServiceBusSubscription under it receives an independent, filtered
copy of the stream.

Topics are many-per-namespace with independent lifecycles, which is why
they are a first-class kind referencing the namespace rather than a list
folded into the namespace's spec. Consumer-side concerns -- lock
duration, delivery counts, sessions, dead-lettering -- live on the
subscription, not the topic. Topics require the STANDARD or PREMIUM
tier (BASIC is queue-only).

**Contracts enforced by Azure at apply time** (they depend on the parent
namespace's tier, which a reference cannot see at validation time):
- In a PARTITIONED premium namespace (premium_messaging_partitions of
  2 or 4) every topic must set partitioning_enabled true; in a
  non-partitioned one it must stay false. On STANDARD, per-topic
  partitioning is a free choice.
- max_message_size_in_kilobytes is PREMIUM-only (multi-tenant tiers are
  fixed at 256 KB).

**ForceNew fields** (changing these replaces the topic and every
subscription under it): `topic_name`, `partitioning_enabled`,
`requires_duplicate_detection`.

## Example

```yaml
# Offline-plan manifest: a fully-dialed topic exercising the size
# ladder, duplicate detection with its window, publish ordering, and the
# gate-state enum mapping.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureServiceBusTopic
metadata:
  name: test-sb-topic
spec:
  namespaceId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.ServiceBus/namespaces/hack-servicebus-ns
  topicName: billing/invoices
  maxSizeInMegabytes: 4096
  requiresDuplicateDetection: true
  duplicateDetectionHistoryTimeWindow: PT1M
  defaultMessageTtl: P30D
  autoDeleteOnIdle: PT10M
  batchedOperationsEnabled: true
  supportOrdering: true
  status: DISABLED
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespaceId` | `string \| valueFrom` | yes |  | AzureServiceBusNamespace (`status.outputs.namespace_id`) |
| `spec.topicName` | `string` | yes |  |  |
| `spec.maxSizeInMegabytes` | `int32` |  |  |  |
| `spec.maxMessageSizeInKilobytes` | `int32` |  |  |  |
| `spec.partitioningEnabled` | `bool` |  |  |  |
| `spec.requiresDuplicateDetection` | `bool` |  |  |  |
| `spec.defaultMessageTtl` | `string` |  |  |  |
| `spec.duplicateDetectionHistoryTimeWindow` | `string` |  |  |  |
| `spec.autoDeleteOnIdle` | `string` |  |  |  |
| `spec.batchedOperationsEnabled` | `bool` |  | `true` |  |
| `spec.expressEnabled` | `bool` |  |  |  |
| `spec.supportOrdering` | `bool` |  |  |  |
| `spec.status` | `enum` |  |  |  |

## Field Details

### spec.namespaceId

`string | valueFrom` · required

The Service Bus namespace the topic lives in, by ARM ID. References an
AzureServiceBusNamespace's namespace_id output so the namespace and
its topics compose in one manifest set. Fixed at creation.

- references: AzureServiceBusNamespace (`status.outputs.namespace_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureServiceBusNamespace, name: <that resource's name>, fieldPath: status.outputs.namespace_id}} -- a bare string does not parse

### spec.topicName

`string` · required

The topic's name -- unique within the namespace, up to 260 characters.
Starts and ends with a letter or number; letters, numbers, periods,
hyphens, underscores, tildes, and forward slashes in between (slashes
create a logical hierarchy, e.g. "billing/invoices").

**ForceNew**: changing the name replaces the topic and every
subscription under it.

- rule: topic_name must start and end with a letter or number and may contain letters, numbers, periods, hyphens, underscores, tildes, and forward slashes (max 260 characters)
- rule: {"required":true,"string":{"minLen":"1","maxLen":"260"}}

### spec.maxSizeInMegabytes

`int32` · optional (explicit presence)

Maximum topic size in megabytes -- how much message data the topic
(across all its subscriptions) holds before publishers get rejected.
Azure only sells fixed sizes: 1024, 2048, 3072, 4096, or 5120 MB on
any tier; 10240, 20480, 40960, or 81920 MB on PREMIUM. Unset lets
Azure default the size for the namespace's tier. For partitioned
STANDARD topics Azure multiplies the chosen size by the 16 internal
partitions -- the effective capacity is 16x.

- rule: {"int32":{"in":[1024,2048,3072,4096,5120,10240,20480,40960,81920]}}

### spec.maxMessageSizeInKilobytes

`int32` · optional (explicit presence)

The largest message the topic accepts, in kilobytes. PREMIUM only --
multi-tenant tiers are fixed at 256 KB. Range: 1024 (1 MB) to 102400
(100 MB). Unset keeps Azure's default for the tier.

- rule: {"int32":{"lte":102400,"gte":1024}}

### spec.partitioningEnabled

`bool` · optional (explicit presence)

Whether the topic is spread across multiple message stores for higher
throughput. On STANDARD this is a free per-topic choice. On PREMIUM it
is dictated by the namespace: a partitioned namespace requires true, a
non-partitioned one requires false.

**ForceNew**: fixed at creation.
Default: false

### spec.requiresDuplicateDetection

`bool` · optional (explicit presence)

Whether Service Bus tracks MessageId values and silently drops
duplicates arriving within duplicate_detection_history_time_window.
Enables idempotent publishers that retry sends safely.

**ForceNew**: fixed at creation.
Default: false

### spec.defaultMessageTtl

`string` · optional (explicit presence)

Default time-to-live for messages, as an ISO 8601 duration (e.g.
"P14D", "PT1H"). Messages older than this are removed from every
subscription that has not consumed them (subscriptions may set a
shorter TTL of their own). Unset means unbounded.

### spec.duplicateDetectionHistoryTimeWindow

`string` · optional (explicit presence)

How long the duplicate-detection history remembers MessageIds, as an
ISO 8601 duration. Only meaningful with requires_duplicate_detection.
Azure's default: PT10M (10 minutes).

### spec.autoDeleteOnIdle

`string` · optional (explicit presence)

Whether the topic is automatically deleted after sitting idle (no
sends or receives) for this ISO 8601 duration. Minimum PT5M. Unset
means never auto-delete -- the right posture for durable
infrastructure.

### spec.batchedOperationsEnabled

`bool` · optional (explicit presence)

Whether clients may batch multiple operations into one broker call
for higher throughput. Azure's default is true.
Default: true

- default: `true`

### spec.expressEnabled

`bool` · optional (explicit presence)

Whether Express Entities are enabled: messages are held in memory
before being written to storage, trading durability for latency.
STANDARD only -- PREMIUM does not support express topics.
Incompatible with duplicate detection.
Default: false

### spec.supportOrdering

`bool` · optional (explicit presence)

Whether the topic preserves publish order when delivering to
subscriptions. Pair with session-aware subscriptions for
strictly-ordered publish-subscribe.
Default: false

### spec.status

`enum`

The topic's gate state: ACTIVE (normal) or DISABLED (publishes
rejected; subscriptions retained). Unspecified deploys ACTIVE. Topics
do not support the queue's one-directional gate states -- direction
gating happens per subscription.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_service_bus_topic_status_unspecified` -- Not specified -- deploys ACTIVE.
- `ACTIVE` -- Normal operation: publishes flow to subscriptions.
- `DISABLED` -- Publishes are rejected; the topic and its subscriptions are retained.

## Validation Rules

- `service_bus_topic_express_conflicts_with_duplicate_detection`: express_enabled holds messages in memory before storage, which is incompatible with duplicate detection -- disable one of the two
- `service_bus_topic_dedup_window_requires_dedup`: duplicate_detection_history_time_window only applies when requires_duplicate_detection is true

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureServiceBusTopic, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.topic_id` | `string` | The Azure Resource Manager ID of the topic. The parent reference for AzureServiceBusSubscription and topic-scoped SAS rules, and the scope for topic-level data-plane role assignments (Azure Service Bus Data Sender). Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.ServiceBus/namespaces/{ns}/topics/{name} |
| `status.outputs.topic_name` | `string` | The topic's name -- what SDKs and function bindings reference within the namespace. |
| `status.outputs.namespace_name` | `string` | The name of the Service Bus namespace the topic lives in, parsed from the resolved namespace ID -- saves consumers a second reference when they need the namespace/topic pair. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespaceId` | AzureServiceBusNamespace | `status.outputs.namespace_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureServiceBusAuthorizationRule | `spec.topicId` | `status.outputs.topic_id` |
| AzureServiceBusSubscription | `spec.topicId` | `status.outputs.topic_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
