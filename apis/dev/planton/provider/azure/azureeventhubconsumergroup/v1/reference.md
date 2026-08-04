# AzureEventHubConsumerGroup

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1`

**AzureEventHubConsumerGroupSpec** defines the configuration for creating
a consumer group on an event hub: one independent read cursor over the
hub's partitions.

A consumer group is a named view of the entire event stream. Each group
tracks its own offsets, so multiple applications -- a real-time
processor, a batch loader, an anomaly detector -- consume the same
events independently at their own pace. Consumer groups are
many-per-hub with independent lifecycles (each consuming team owns its
own group), which is why they are a first-class kind referencing the
hub rather than a list folded into the hub's spec.

**Tier limits enforced by Azure at apply time**: BASIC hubs allow no
additional consumer groups (only the service-created $Default);
STANDARD allows 20 per hub; PREMIUM/dedicated allow more.

**The $Default group**: Azure creates a consumer group named "$Default"
on every hub automatically. It cannot be declared here -- Azure's
providers refuse to adopt service-created resources -- and SDK
quick-starts use it implicitly. Give every real consumer application
its OWN group; sharing $Default across applications makes their
offsets collide.

**ForceNew fields**: `event_hub_id`, `consumer_group_name` (the group
is its name; renaming replaces it and resets stored offsets).

## Example

```yaml
# Offline-plan manifest: a consumer group on an existing hub,
# exercising the ARM-id parse into discrete names and the user_metadata
# passthrough.
apiVersion: azure.planton.dev/v1
kind: AzureEventHubConsumerGroup
metadata:
  name: test-eh-consumer-group
spec:
  eventHubId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.EventHub/namespaces/hack-eventhub-ns/eventhubs/orders-stream
  consumerGroupName: analytics-loader
  userMetadata: "owner=analytics-team purpose=hourly-batch-load"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.eventHubId` | `string \| valueFrom` | yes |  | AzureEventHub (`status.outputs.event_hub_id`) |
| `spec.consumerGroupName` | `string` | yes |  |  |
| `spec.userMetadata` | `string` |  |  |  |

## Field Details

### spec.eventHubId

`string | valueFrom` · required

The event hub the group reads, by ARM ID. References an
AzureEventHub's event_hub_id output so the hub and its consumer
groups compose in one manifest set. Fixed at creation.

- references: AzureEventHub (`status.outputs.event_hub_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureEventHub, name: <that resource's name>, fieldPath: status.outputs.event_hub_id}} -- a bare string does not parse

### spec.consumerGroupName

`string` · required

The group's name -- unique within the hub, 1-50 characters. Starts
and ends with a letter or number; letters, numbers, periods, hyphens,
and underscores in between.

"$Default" is reserved: Azure creates that group on every hub, and
an existing service-created group cannot be adopted declaratively.

**ForceNew**: renaming replaces the group and resets its consumers'
stored offsets.

- rule: consumer_group_name must start and end with a letter or number and may contain letters, numbers, periods, hyphens, and underscores (max 50 characters)
- rule: "$Default" is the service-created catch-all group Azure adds to every hub -- it cannot be declared or adopted; SDKs use it implicitly, and every real application should get its own group
- rule: {"required":true,"string":{"minLen":"1","maxLen":"50"}}

### spec.userMetadata

`string` · optional (explicit presence)

Free-form metadata stored on the group (max 1024 characters) --
record the owning application, team, or purpose so operators can
tell whose cursor this is.

- rule: {"string":{"maxLen":"1024"}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureEventHubConsumerGroup, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.consumer_group_id` | `string` | The Azure Resource Manager ID of the consumer group. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.EventHub/namespaces/{ns}/eventhubs/{hub}/consumergroups/{name} |
| `status.outputs.consumer_group_name` | `string` | The group's name -- what consumer applications pass to their SDK client alongside the hub name. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.eventHubId` | AzureEventHub | `status.outputs.event_hub_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
