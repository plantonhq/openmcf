# CloudEvents Hub

This preset creates a pure-CloudEvents namespace -- the shared hub one environment's services publish through, with streams added later as AzureEventgridNamespaceTopic resources.

## When to Use

- One event hub per environment, many services each owning a namespace topic inside it
- Replacing per-service classic topics with a single capacity-managed hub

## Key Configuration Choices

- **No MQTT block** -- this namespace never serves MQTT clients; that choice is create-only, so pick the MQTT preset instead if devices might connect later
- **Capacity 1 TU** -- the floor; raise it in place when throttled-request metrics climb
- **Public access on, no IP rules** -- the default posture; add `inboundIpRules` or disable public access for locked-down networks

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of your `AzureResourceGroup` resource | Planton console (or replace `valueFrom` with `value:` and a literal group name) |
| `my-org-events-hub` | The namespace's name (appears in its hostnames) | Your naming convention -- org-prefixed |
| `eastus` | The Azure region | Your region strategy |

## Related Presets

- **02 MQTT Broker** -- the namespace with the MQTT broker enabled for device fleets
