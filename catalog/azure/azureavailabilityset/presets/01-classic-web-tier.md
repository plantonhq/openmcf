# Classic Web Tier

This preset creates an availability set with the provider defaults -- 5 update domains, 3 fault domains, managed-disk alignment on. Put two or more load-balanced VMs in it and the tier carries Azure's classic 99.95% multi-VM SLA.

## When to Use

- A multi-VM tier in a region WITHOUT availability zones
- Classic lift-and-shift topologies whose design predates zonal placement

## Key Configuration Choices

- **Defaults everywhere** -- the provider's 5/3/managed shape is right for almost every tier; everything except tags is create-only, so the set's shape is decided before the first VM exists
- **Two VMs minimum** -- a single-VM availability set provides nothing; the SLA starts at two
- **Zones beat sets** -- in zoned regions, pin VMs to zones instead (a VM uses a zone OR a set, never both)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of your `AzureResourceGroup` resource | Planton console (or replace `valueFrom` with `value:` and a literal group name) |
| `web-avset` | The set's name (up to 80 letters, numbers, dots, dashes, underscores) | Your naming convention |
| `eastus` | The Azure region -- VMs joining the set must be in the same region | Your region strategy |

## Related Presets

- None yet -- proximity-placement pairing arrives when the P2 proximity placement group kind lands.
