# Resource Group Events

This preset creates a system topic on a resource group -- exposing resource-lifecycle events (write, delete, action success/failure) for governance and audit automation.

## When to Use

- Audit trails: react to every resource created, changed, or deleted in the group
- Drift and hygiene automation: notify or remediate when unexpected resources appear

## Key Configuration Choices

- **Topic type `Microsoft.Resources.ResourceGroups`** -- the ARM control plane's event catalog for one group
- **Region `Global`** -- resource groups are global sources; Azure rejects a regional placement for this topic type
- **Subscription-shaped events** -- lifecycle events carry the operation and resource ID; pair the topic with an event subscription filtering `included_event_types` (e.g. only `Microsoft.Resources.ResourceWriteSuccess`)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of the `AzureResourceGroup` the TOPIC lives in | Planton console |
| `<your-resource-group-arm-id>` | The ARM ID of the group being WATCHED (often the same group) | `/subscriptions/{sub}/resourceGroups/{name}` -- the group's properties blade |
