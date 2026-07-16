# Delegation-Ready Zone

A public zone tuned for active operations: the SOA contact points at your DNS team and the negative-caching TTL is lowered so newly created records (certificate validation records especially) become visible fast.

## When to Use

- Zones that back automated certificate flows -- Front Door and Container Apps validation records need to resolve quickly after creation
- Organizations whose operational tooling reads the SOA contact
- Frequently changing zones where a 5-minute "name does not exist" cache is too slow

## Key Configuration Choices

- `soaRecord.email` -- the zone contact in SOA host format (dots instead of `@`)
- `soaRecord.minimumTtl: 60` -- resolvers cache negative answers for one minute instead of five; the one SOA field with real operational impact
- Refresh/retry/expire timers stay at Azure's defaults -- they only matter for zone-transfer secondaries, which Azure DNS does not use

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `example.com` | Replace with your domain (appears in `zoneName` AND the SOA email) | Your domain registrar account |
| `<your-resource-group>` | The resource group name | `AzureResourceGroup.status.outputs.resource_group_name` |

## Related Presets

- `01-public-zone` -- the same zone with Azure's SOA defaults
