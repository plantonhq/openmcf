# AzurePrivateDnsZone -- Design Research

## The Resource

An Azure Private DNS zone (`Microsoft.Network/privateDnsZones`) is a
global record container providing name resolution inside the virtual
networks linked to it. The component maps onto `azurerm_private_dns_zone`
(azurerm v4.x,
`internal/services/privatedns/private_dns_zone_resource.go`),
parity-verified against pulumi-azure v6 (`privatedns.Zone`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `name` | Required, ForceNew; DNS-domain CEL validation |
| `resource_group_name` | `resource_group` | FK → AzureResourceGroup |
| `soa_record` | `soa_record` | Optional block; email required inside; timers default to Azure's values; ForceNew |
| `tags` | `tags` | User tags merged over Planton-derived tags |
| `number_of_record_sets`, `max_number_of_*` (computed) | -- | Capacity introspection, not configuration; not modeled |

There is no `location` -- private DNS zones are global, which is why this
spec (like the DNS zone family generally) omits the catalog's usual
`region` field.

## The Decomposition

The zone models ONLY the zone. Its network attachments are the standalone
`AzurePrivateDnsZoneVirtualNetworkLink` kind, because links fail every
fold test:

- **Many-per-parent**: hub-and-spoke links one zone to the hub and every
  spoke network.
- **Independent lifecycle**: networks join and leave the resolution
  audience without touching the zone or its records.
- **Independently referenced**: each link binds a specific zone-network
  pair with its own settings (registration, resolution policy).

A zone with zero links resolves for nobody -- deployments pair the zone
with at least one link, and the presets teach that composition.

## Design Decisions

- **SOA modeled as an optional message** mirroring azurerm's block:
  presence = customize, absence = Azure's standard SOA. The only
  operationally interesting field is `minimum_ttl` (negative caching);
  the field comments say so rather than presenting all six timers as
  equals.
- **`resource_group_name` added to outputs**: downstream record/link
  tooling frequently joins on zone name + resource group; exporting it
  saves manifests from parsing ARM IDs.
- **Zone `name` keeps its DNS-domain CEL validation** (lowercase labels,
  RFC-length bounds) -- privatelink names are fixed strings that pass it,
  and custom zones get early validation instead of an ARM error.

## Operational Behavior Worth Knowing

- Renaming the zone replaces it AND every record in it -- the highest
  blast-radius change in the private DNS family.
- Deleting a zone requires its links (and records) to be deleted first;
  composed destroy ordering handles this via the dependency graph.
- The SOA record is created with the zone; ARM allows editing some SOA
  fields on the record-set API, but azurerm models the block ForceNew --
  the spec documents it as written-at-creation.

## Composition

- `resource_group` → `AzureResourceGroup.status.outputs.resource_group_name`
- `zone_id` output consumed by:
  `AzurePrivateDnsZoneVirtualNetworkLink.private_dns_zone_id`,
  `AzurePrivateEndpoint.private_dns_zone_id`, and the flexible-server
  database kinds' `private_dns_zone_id` fields.
