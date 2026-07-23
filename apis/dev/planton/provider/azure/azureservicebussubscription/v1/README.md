# AzureServiceBusSubscription

A subscription under an Azure Service Bus topic: an independent,
optionally filtered view of the topic's stream, with its own consumer
semantics (lock duration, delivery counts, sessions, dead-lettering).
Filter rules fold inside the subscription -- they combine with OR
semantics alongside Azure's auto-created `$Default` catch-all (removed
once, out-of-band, for restrictive delivery).

## When to Use

Use AzureServiceBusSubscription when you need:

- **A consumer's view of a topic** -- one subscription per consuming
  application, each reading at its own pace
- **Filtered delivery** -- SQL or correlation rules admit only the
  messages this consumer cares about
- **Fan-out-then-collect routing** -- filter here, `forward_to` a work
  queue for the actual processing fleet

## Key Configuration

- `topic_id` -- the parent topic, referenced from an
  AzureServiceBusTopic output (fixed at creation)
- `max_delivery_count` -- required; Azure has no server default for a
  subscription's delivery tolerance
- `rules` -- folded filter rules, ADDITIVE alongside the auto-created
  catch-all; for restrictive delivery remove the catch-all once:
  `az servicebus topic subscription rule delete --name '$Default' ...`
  (`$Default` cannot be declared -- the service-created rule cannot be
  adopted)
- `requires_session` -- ForceNew; pairs with the topic's
  `support_ordering` for ordered publish-subscribe

## Composition

```yaml
topicId:
  valueFrom:
    kind: AzureServiceBusTopic
    name: order-events
    fieldPath: status.outputs.topic_id
```

## Documentation

- [Design research](docs/README.md) -- field mapping, the $Default
  contract, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
