# AzureDnsZone -- Design Research

## The Resource

An Azure public DNS zone (`Microsoft.Network/dnsZones`) is a global,
internet-facing authoritative zone served by Azure's anycast name-server
fleet. The component maps onto `azurerm_dns_zone` (azurerm v4.x,
`internal/services/dns/dns_zone_resource.go`, DNS API 2018-05-01),
parity-verified against pulumi-azure v6 (`dns.Zone`).

## Field Mapping (azurerm -> spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `zone_name` | ForceNew; renaming replaces the zone, its records, and the assigned name-server set |
| `resource_group_name` | `resource_group` | FK-defaults to `AzureResourceGroup.resource_group_name`. ForceNew |
| `soa_record` | `soa_record` | Optional fold (below) |
| `tags` | `tags` | User tags merged over metadata-derived tags |

SOA block: `email` (required; the provider's segment rules -- 2-34
dot-separated segments of letters/digits/`._-`, each <=63 chars -- are
mirrored as a CEL), `expire_time` (2419200), `minimum_ttl` (300),
`refresh_time` (3600), `retry_time` (300), `serial_number` (1), `ttl`
(3600), `tags`. Defaults mirror what Azure creates when the block is
omitted, so a partial block and Azure's own values deploy identically.
The provider's create/update check that `len(zone_name) +
len(soa_email)` <= 253 is front-loaded as a message CEL.

## Decomposition Decision

The zone previously bundled a `records` list -- DISSOLVED. Records are
many-per-zone with independent lifecycles (app deployments add and
remove them constantly), and record management is exactly what the
standalone `AzureDnsRecord` kind models. A zone-embedded list would
force zone updates for every record change and could not express the
typed per-record shapes (MX preferences, SRV endpoints, alias targets).

## Deliberately Not Modeled (with reasons)

- **`soa_record.host_name`** -- Azure owns the SOA host and rejects
  changes; the provider itself preserves the API-assigned value on
  writes. Modeling it would be contradictable state.
- **`number_of_record_sets`** (computed) -- a point-in-time count that
  is stale the moment a record deploys; no composition value.
  `max_number_of_record_sets` (a capacity fact) IS exported.
- **DNSSEC / DS / TLSA / NAPTR** -- absent from the azurerm v4 surface
  entirely; lands when the provider models it.
- **Apex NS management** -- the zone's own NS record set is
  Azure-managed; child-zone delegation is an `AzureDnsRecord` NS record.

## Operational Behavior Worth Knowing

- Public DNS zone names are NOT globally unique -- the same zone can
  exist in many subscriptions, each with a different name-server set;
  only the registrar-delegated copy answers the internet.
- Azure does not auto-increment the SOA serial on record changes (its
  name servers replicate internally), so `serial_number` is cosmetic
  metadata.
- Zones are global resources: no region, and the resource group governs
  only management-plane lifecycle.
