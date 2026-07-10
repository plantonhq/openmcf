# AzureServiceBusQueue

A queue inside an Azure Service Bus namespace: reliable point-to-point
messaging with FIFO delivery, at-least-once semantics, sessions,
duplicate detection, and a built-in dead-letter sub-queue. The queue's
ARM ID is a data-plane RBAC scope, and queue-scoped SAS credentials are
minted with AzureServiceBusAuthorizationRule.

## When to Use

Use AzureServiceBusQueue when you need:

- **Point-to-point work distribution** -- competing consumers draining
  one queue (orders, jobs, commands)
- **Ordered, stateful processing** -- `requires_session` for strict
  FIFO per session with exclusive consumption
- **Idempotent producers** -- `requires_duplicate_detection` drops
  retried sends within the detection window
- **Poison-message handling** -- `max_delivery_count` +
  `forward_dead_lettered_messages_to` centralize failures

## Key Configuration

- `namespace_id` -- the parent namespace, referenced from an
  AzureServiceBusNamespace output (fixed at creation)
- `queue_name` -- unique within the namespace; slashes create logical
  hierarchy (`orders/priority`)
- `lock_duration` / `max_delivery_count` -- sized to the consumer's
  processing profile
- `partitioning_enabled`, `requires_duplicate_detection`,
  `requires_session` -- ForceNew: fixed at creation
- Premium-gated (apply-time, tier lives on the namespace): large
  messages via `max_message_size_in_kilobytes`; express is
  BASIC/STANDARD-only

## Composition

```yaml
namespaceId:
  valueFrom:
    kind: AzureServiceBusNamespace
    name: orders-bus
    fieldPath: status.outputs.namespace_id
```

Grant receive rights on just this queue via AzureRoleAssignment on
`status.outputs.queue_id` (Azure Service Bus Data Receiver), or mint a
queue-scoped SAS rule.

## Documentation

- [Design research](docs/README.md) -- field mapping, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)
