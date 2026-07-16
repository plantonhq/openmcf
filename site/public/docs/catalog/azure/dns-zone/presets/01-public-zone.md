---
title: "Public Zone"
description: "The 30-second public DNS zone: an empty, internet-facing zone for a domain you own, ready for records and registrar delegation. Azure's defaults cover the SOA record; governance tags carry your..."
type: "preset"
rank: "01"
presetSlug: "01-public-zone"
componentSlug: "dns-zone"
componentTitle: "DNS Zone"
provider: "azure"
icon: "package"
order: 1
---

# Public Zone

The 30-second public DNS zone: an empty, internet-facing zone for a domain you own, ready for records and registrar delegation. Azure's defaults cover the SOA record; governance tags carry your ownership conventions.

## When to Use

- Hosting a domain's DNS on Azure (new domain, or migrating from another provider)
- The DNS foundation for custom-domain flows (Front Door, Container Apps) that publish validation records into the zone
- A delegated subdomain zone a team manages independently (set `zoneName` to the subdomain)

## Key Configuration Choices

- `zoneName` -- the domain, with no trailing dot; renaming replaces the zone and its records
- `resourceGroup` -- management-plane scope only; resolution is unaffected
- SOA settings are omitted deliberately -- Azure's defaults are right for nearly everyone

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `example.com` | Replace with the domain the zone hosts | Your domain registrar account |
| `<your-resource-group>` | The resource group name | `AzureResourceGroup.status.outputs.resource_group_name`, or the Azure portal |

## Related Presets

- `02-delegation-ready-zone` -- adds SOA contact customization for operational tooling
