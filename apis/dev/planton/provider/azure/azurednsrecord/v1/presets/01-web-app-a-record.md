# Web App A Record

The 30-second DNS record: point a hostname (www) at an IPv4 address. Add more addresses to the list for round-robin distribution.

## When to Use

- Pointing a website or API hostname at a load balancer's or VM's public address
- Round-robin across several addresses (add each to `addresses`)

## Key Configuration Choices

- `name: www` -- relative to the zone; use `@` for the apex (your-domain.com itself)
- `ttlSeconds: 300` -- how long resolvers cache the answer; lower it before planned moves
- For an address that belongs to an Azure resource, prefer the alias form (`a.targetResourceId` referencing an `AzurePublicIp`) -- Azure then tracks address changes with no drift window

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-resource-group>` | The zone's resource group | `AzureResourceGroup.status.outputs.resource_group_name` |
| `example.com` | Replace with the zone name | `AzureDnsZone.status.outputs.zone_name` |
| `203.0.113.10` | Replace with the address to answer with | Your load balancer / VM / gateway's public IP |

## Related Presets

- `02-apex-alias-record` -- the zone apex tracking an Azure Public IP
- `03-mail-mx-records` -- mail-exchange records with preferences
