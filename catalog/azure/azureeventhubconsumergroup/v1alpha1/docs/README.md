# AzureEventHubConsumerGroup -- Design Research

## The Resource

An Event Hubs consumer group
(`Microsoft.EventHub/namespaces/eventhubs/consumergroups`) is one
application's independent read cursor over a hub's partitions. The
component maps onto `azurerm_eventhub_consumer_group` (azurerm v4.x,
`internal/services/eventhub/eventhub_consumer_group_resource.go`),
parity-verified against pulumi-azure v6
(`eventhub.EventHubConsumerGroup`).

Consumer groups are many-per-hub with independent lifecycles -- each
consuming team owns its own group -- which is why they are a
first-class kind referencing the hub rather than a list folded into the
hub's spec.

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `resource_group_name` / `namespace_name` / `eventhub_name` | `event_hub_id` | azurerm addresses the parent by discrete names; the modules derive all three from the spec's single ARM-id reference (see below) |
| `name` | `consumer_group_name` | Required, ForceNew; 1-50 chars, alphanumeric ends, `[-._]` interior; `$Default` reserved by a spec CEL |
| `user_metadata` | same | Optional, max 1024 chars; ownership breadcrumbs operators see in the portal |

Outputs: `consumer_group_id`, `consumer_group_name`.

**Parent addressing**: the spec carries one authoritative parent
reference (`event_hub_id`, an ARM id from the AzureEventHub output);
both modules parse the resource group, namespace, and hub names out of
it with an identical anchored regex that fails loudly on a malformed
id. This keeps the spec on the ARM-id grain with no redundant name
fields that could contradict each other.

## Recorded Skips (with reasons)

- **The `$Default` group** -- Azure creates it on every hub, and the
  providers refuse to adopt service-created resources; the spec
  reserves the name with a CEL rule rather than letting the apply fail
  opaquely. SDKs use it implicitly; real applications get their own
  groups.
- **Tier quotas stay apply-time** -- BASIC (no additional groups) and
  STANDARD (20 per hub) limits depend on the referenced namespace's
  live tier, which validation cannot see. Documented on the spec; Azure
  enforces them verbatim.
- **No tags** -- ARM does not support tags on Event Hubs entities
  (hubs/consumer groups/schema groups); the platform's identity tags
  live on the parent namespace.

## Operational Behavior Worth Knowing

- **Consumer groups create in seconds** -- the namespace dominates any
  composed deploy.
- **The group is its name**: renaming replaces it, and the replacement
  starts with no stored offsets -- consumers restart from their
  configured default position.
- **One group per application** is the consumption discipline; two
  applications sharing a group steal each other's events.

## Composition

- `event_hub_id` → `AzureEventHub.status.outputs.event_hub_id`
- `consumer_group_name` output ← consumer SDK clients, function
  bindings (passed alongside the hub name)
