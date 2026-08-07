# AzureServiceBusTopic -- Design Research

## The Resource

A Service Bus topic (`Microsoft.ServiceBus/namespaces/topics`) is the
publish-subscribe entity; subscriptions under it are the consuming
views. The component maps onto `azurerm_servicebus_topic` (azurerm
v4.x, `internal/services/servicebus/servicebus_topic_resource.go`),
parity-verified against pulumi-azure v6 (`servicebus.Topic`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `namespace_id` | `namespace_id` | The v4 typed parent id; single authoritative parent FK, ForceNew |
| `name` | `topic_name` | Required, ForceNew, provider regex (≤260, start/end alphanumeric) |
| `max_size_in_megabytes` | same | Azure's fixed ladder as a closed int set |
| `max_message_size_in_kilobytes` | same | 1024-102400; Premium-only (apply-time) |
| `partitioning_enabled` | same | ForceNew; Premium pairing apply-time |
| `requires_duplicate_detection` | same | ForceNew |
| `duplicate_detection_history_time_window` | same | Paired-with-dedup CEL |
| `default_message_ttl` / `auto_delete_on_idle` | same | Unset leaves Azure's defaults |
| `batched_operations_enabled` | same | optional bool, default true |
| `express_enabled` | same | STANDARD only; express-vs-dedup CEL |
| `support_ordering` | same | Pairs with session-aware subscriptions |
| `status` | enum | ACTIVE/DISABLED (the provider's own restricted set for topics) |

Outputs: `topic_id`, `topic_name`, `namespace_name`.

## Design Notes

- **No consumer semantics here** -- lock duration, delivery counts,
  sessions, and dead-lettering are subscription-level concerns and live
  on AzureServiceBusSubscription. This is Azure's own grain, not a
  simplification.
- **The queue's express-on-Premium apply-time block has no topic
  counterpart in the provider** -- express on a Premium topic surfaces
  as ARM's own rejection; documented rather than invented as a CEL.

## Recorded Skips (with reasons)

- **Transitional `status` values** -- topics accept only
  Active/Disabled in azurerm; modeled as-is.

## Operational Behavior Worth Knowing

- **Topics create in seconds**; deleting a topic deletes every
  subscription under it.
- **Partitioned STANDARD topics report 16x the sold size**; the
  provider normalizes the read back.

## Composition

- `namespace_id` → `AzureServiceBusNamespace.status.outputs.namespace_id`
- `topic_id` output ← AzureServiceBusSubscription.topic_id,
  topic-scoped AzureServiceBusAuthorizationRule, data-plane role
  assignments (Azure Service Bus Data Sender)
