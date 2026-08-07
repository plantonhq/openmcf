# Apex Alias Record

Point the zone apex (your-domain.com itself) at an Azure Public IP with an alias record. Azure keeps the answer in sync with the resource -- when the IP's address changes, the record follows automatically, with no drift window and no stale-IP outage.

## When to Use

- Serving a naked domain from an Azure load balancer or application gateway (DNS forbids CNAME at the apex; the alias A record is the mechanism)
- Any A/AAAA record whose address belongs to an Azure resource -- prefer the alias over copying the literal address

## Key Configuration Choices

- `a.targetResourceId` -- also accepts Traffic Manager profiles, CDN endpoints, and Front Door endpoints; reference the resource's ARM-id output
- Alias and literal `addresses` are mutually exclusive -- an alias record delegates its answer entirely

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-resource-group>` | The zone's resource group | `AzureResourceGroup.status.outputs.resource_group_name` |
| `example.com` | Replace with the zone name | `AzureDnsZone.status.outputs.zone_name` |
| `<your-public-ip-resource>` | The AzurePublicIp resource's metadata name | Your Planton resource inventory |

## Related Presets

- `01-web-app-a-record` -- a literal-address record for non-Azure targets
