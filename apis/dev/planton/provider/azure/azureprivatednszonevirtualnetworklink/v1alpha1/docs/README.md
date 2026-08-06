# AzurePrivateDnsZoneVirtualNetworkLink -- Design Research

## The Resource

A virtual network link
(`Microsoft.Network/privateDnsZones/{zone}/virtualNetworkLinks/{name}`) is
the ARM child resource that makes a private DNS zone resolvable from one
virtual network. The component maps onto
`azurerm_private_dns_zone_virtual_network_link` (azurerm v4.x,
`internal/services/privatedns/private_dns_zone_virtual_network_link_resource.go`),
parity-verified against pulumi-azure v6
(`privatedns.ZoneVirtualNetworkLink`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `name` | Required, ForceNew (case-insensitive in ARM today) |
| `private_dns_zone_name` + `resource_group_name` | `private_dns_zone_id` | DERIVED -- see the parent-reference decision below |
| `virtual_network_id` | `virtual_network_id` | Required, ForceNew; FK → AzureVirtualNetwork |
| `registration_enabled` | `registration_enabled` | azurerm default false → `optional bool` default "false" |
| `resolution_policy` | `resolution_policy` enum | Optional+Computed in azurerm (ARM sets a per-zone-type default); only an explicit value is sent |
| `tags` | `tags` | User tags merged over Planton-derived tags |

## The Parent-Reference Decision

azurerm v4's contract takes the zone by NAME + RESOURCE GROUP; its own
source marks that as the legacy shape it intends to replace with
`private_dns_zone_id`. This component models the target shape directly:
one `private_dns_zone_id` reference, with both modules deriving the zone
name and resource group from the ARM ID
(`.../resourceGroups/{rg}/providers/Microsoft.Network/privateDnsZones/{zone}`).

Rationale: the link is an ARM child of the zone, and the zone's ID already
encodes everything -- modeling name+RG as separate fields would create
redundant state that can contradict the referenced zone. The catalog's
parent-child precedent (the federated identity credential deriving its
resource group from the parent identity's ID) applies unchanged.

## Design Decisions

- **`resolution_policy` as a closed proto enum** with unspecified meaning
  "Azure chooses": ARM computes a default that differs by zone type
  (privatelink vs custom), so forcing a value would fight the platform.
  `NxDomainRedirect` is the fallback pattern for privatelink zones shared
  across environments where some records exist only publicly.
- **`registration_enabled` default false**, with the field comment
  carrying the two operational rules users hit: privatelink zones must
  keep it off (private endpoints own those records), and Azure allows only
  ONE registration-enabled link per network.
- **Link name defaults to nothing** -- it is required and should be named
  after the network it attaches, because a zone accumulates one link per
  network and the name is how they are told apart.

## Operational Behavior Worth Knowing

- Link creation is fast, but resolution through a newly-linked network can
  take a short propagation window.
- Replacing a link (rename, re-zone, re-network) causes a brief resolution
  gap for the affected network only -- records and other links are
  untouched.
- Deleting a zone requires its links to be gone; composed destroy ordering
  handles this through the dependency graph.

## Composition

- `private_dns_zone_id` → `AzurePrivateDnsZone.status.outputs.zone_id`
- `virtual_network_id` → `AzureVirtualNetwork.status.outputs.virtual_network_id`
- Registry prerequisites: `[AzurePrivateDnsZone, AzureVirtualNetwork]`
  (the resource group chains transitively through both).
