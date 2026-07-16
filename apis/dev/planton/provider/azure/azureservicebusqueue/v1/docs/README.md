# AzureServiceBusQueue -- Design Research

## The Resource

A Service Bus queue (`Microsoft.ServiceBus/namespaces/queues`) is the
point-to-point messaging entity. The component maps onto
`azurerm_servicebus_queue` (azurerm v4.x,
`internal/services/servicebus/servicebus_queue_resource.go`),
parity-verified against pulumi-azure v6 (`servicebus.Queue`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `namespace_id` | `namespace_id` | The v4 typed parent id; single authoritative parent FK, ForceNew |
| `name` | `queue_name` | Required, ForceNew, provider regex (≤260, start/end alphanumeric, `[\w-./~]` interior) |
| `max_size_in_megabytes` | same | Azure's fixed ladder as a closed int set (large sizes Premium-only, apply-time) |
| `max_message_size_in_kilobytes` | same | 1024-102400; Premium-only (apply-time -- the tier lives on the referenced namespace) |
| `partitioning_enabled` | same | ForceNew; Premium partition-layout pairing is apply-time (cross-resource) |
| `requires_duplicate_detection` | same | ForceNew |
| `duplicate_detection_history_time_window` | same | Paired-with-dedup CEL (front-loads a silent no-op) |
| `default_message_ttl` / `lock_duration` / `max_delivery_count` / `auto_delete_on_idle` | same | Unset leaves Azure's defaults (unbounded / PT1M / 10 / never) |
| `dead_lettering_on_message_expiration` | same | Default false |
| `requires_session` | same | ForceNew |
| `batched_operations_enabled` | same | optional bool, default true (Azure's default) |
| `express_enabled` | same | BASIC/STANDARD only (apply-time); the express-vs-dedup conflict is a spec CEL |
| `forward_to` / `forward_dead_lettered_messages_to` | same | Entity NAMES in the same namespace (Azure's own addressing) |
| `status` | enum | ACTIVE/DISABLED/SEND_DISABLED/RECEIVE_DISABLED (see skips) |

Outputs: `queue_id`, `queue_name`, `namespace_name` (parsed from the
resolved parent id with identical anchored parsing on both engines).

## Recorded Skips (with reasons)

- **Transitional `status` states** (Creating, Deleting, Renaming,
  Unknown, and the provider's permissive acceptance of them) -- server
  states, not knobs; the enum models the four administratively settable
  gates.
- **Cross-resource Premium contracts stay apply-time**: express-on-
  Premium, the partition-layout pairing, and large-message tier gating
  depend on the referenced namespace's sku, which a reference cannot see
  at validation time. Documented on the fields; Azure enforces them
  verbatim.

## Operational Behavior Worth Knowing

- **Queues create in seconds** -- the namespace dominates any composed
  deploy.
- **Partitioned multi-tenant queues report 16x the sold size**; the
  provider normalizes the read back, so plans stay clean.
- **Forward targets must exist first** -- Azure rejects forwarding to a
  missing entity; order compositions accordingly.

## Composition

- `namespace_id` → `AzureServiceBusNamespace.status.outputs.namespace_id`
- `queue_id` output ← queue-scoped AzureServiceBusAuthorizationRule,
  data-plane role assignments (Azure Service Bus Data Receiver/Sender)
- `queue_name` + the namespace's `endpoint` ← SDK clients, function
  bindings
