# AzureServiceBusTopic

A topic inside an Azure Service Bus namespace: the publish-subscribe
primitive. Publishers send to the topic; each AzureServiceBusSubscription
under it receives an independent, filtered copy of the stream. Topics
require the STANDARD or PREMIUM tier (BASIC is queue-only).

## When to Use

Use AzureServiceBusTopic when you need:

- **Event broadcast** -- one publish, many independent consumers
  (order-created feeding billing, inventory, and analytics)
- **Filtered routing** -- subscriptions admit only matching messages
  (SQL or correlation filters), so consumers never see noise
- **Ordered publish-subscribe** -- `support_ordering` with
  session-aware subscriptions

## Key Configuration

- `namespace_id` -- the parent namespace, referenced from an
  AzureServiceBusNamespace output (fixed at creation)
- `topic_name` -- unique within the namespace; slashes create logical
  hierarchy (`billing/invoices`)
- `partitioning_enabled`, `requires_duplicate_detection` -- ForceNew:
  fixed at creation
- Consumer-side dials (locks, delivery counts, sessions,
  dead-lettering) live on the SUBSCRIPTION, not here

## Composition

```yaml
namespaceId:
  valueFrom:
    kind: AzureServiceBusNamespace
    name: events-bus
    fieldPath: status.outputs.namespace_id
```

Subscriptions reference `status.outputs.topic_id`; grant publish rights
via a topic-scoped AzureServiceBusAuthorizationRule or a data-plane
role assignment (Azure Service Bus Data Sender) on the topic's ARM id.

## Documentation

- [Design research](docs/README.md) -- field mapping, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)
