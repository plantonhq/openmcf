# AzureEventHubConsumerGroup

A consumer group on an Azure event hub: one application's independent
read cursor over the hub's partitions. Each group tracks its own
offsets, so a real-time processor, a batch loader, and an anomaly
detector can consume the same stream at their own pace without
interfering with each other.

## When to Use

Use AzureEventHubConsumerGroup when you need:

- **One group per consuming application** -- the core discipline of
  Event Hubs consumption; applications sharing a group steal each
  other's events because their offsets collide
- **Independent replay** -- a new consumer rewinds through retained
  events from its own cursor without disturbing existing readers
- **Operational attribution** -- `user_metadata` records the owning
  application or team, so operators can tell whose cursor each group is

Know the reservation: Azure creates a group named `$Default` on every
hub, and service-created resources cannot be adopted declaratively --
the spec reserves the name. SDK quick-starts use `$Default` implicitly;
production applications should never share it.

## Key Configuration

- `event_hub_id` -- the parent hub, referenced from an AzureEventHub
  output (fixed at creation)
- `consumer_group_name` -- unique within the hub, 1-50 characters;
  ForceNew: the group is its name, so renaming replaces it and resets
  its consumers' stored offsets
- `user_metadata` -- free-form ownership breadcrumbs (max 1024
  characters) visible to operators in the portal

Tier limits are enforced by Azure at apply time: BASIC hubs allow no
additional groups beyond the service-created `$Default`; STANDARD
allows 20 per hub; PREMIUM and dedicated allow more.

## Composition

```yaml
eventHubId:
  valueFrom:
    kind: AzureEventHub
    name: telemetry-stream
    fieldPath: status.outputs.event_hub_id
```

Consumers pass `status.outputs.consumer_group_name` to their SDK client
alongside the hub name.

## Documentation

- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
