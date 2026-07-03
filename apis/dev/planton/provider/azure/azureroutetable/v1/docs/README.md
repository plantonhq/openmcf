# AzureRouteTable -- Design Research

## The Resource

An Azure route table (`Microsoft.Network/routeTables`) holds user-defined
routes (UDRs) that override Azure's system routes for the subnets attached
to it. The component maps onto `azurerm_route_table` (azurerm v4.x,
`internal/services/network/route_table_resource.go`), parity-verified
against pulumi-azure v6 (`network.RouteTable`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `name` | Required, ForceNew |
| `location` | `region` | Required, ForceNew |
| `resource_group_name` | `resource_group` | FK → AzureResourceGroup |
| `route` (set) | `routes` (repeated message) | Inline-managed; see folding decision below |
| `route.*.next_hop_type` | `next_hop_type` enum | VirtualNetworkGateway / VnetLocal / Internet / VirtualAppliance / None |
| `route.*.next_hop_in_ip_address` | `next_hop_in_ip_address` | Message-level CEL: required iff VIRTUAL_APPLIANCE |
| `bgp_route_propagation_enabled` | `bgp_route_propagation_enabled` | azurerm default true → `optional bool` with default "true" |
| `subnets` (computed) | -- | Not modeled; membership is subnet-side (see below) |
| `tags` | `tags` | User tags merged over Planton-derived tags |

## Decomposition Decisions

- **Routes FOLD inside the table.** azurerm offers both inline `route`
  blocks and a standalone `azurerm_route` resource; a route has no
  independent lifecycle, is never FK-referenced, and is one-per-table --
  it fails every split test. Both modules manage routes inline and always
  send the list (an explicitly empty list removes the last route; leaving
  the field unset would make the provider treat existing routes as
  externally managed).
- **Subnet attachment is expressed from the SUBNET side.** Azure's model
  is that a subnet declares its route table
  (`azurerm_subnet_route_table_association`), and one table serves many
  subnets. The table therefore exports `route_table_id` for subnets to
  reference, and never lists its consumers -- the same
  never-mutate-what-you-reference principle the rest of the catalog
  follows. The computed `subnets` attribute is deliberately not an output:
  it is derived state owned by the subnets.

## Design Decisions

- **`next_hop_type` as a closed proto enum** with `defined_only` -- the
  five ARM values are stable API surface, and an enum gives wizards and
  agents the full option set from the spec alone.
- **The appliance-IP pairing enforced in the spec** (message-level CEL):
  ARM rejects `next_hop_in_ip_address` on any hop type except
  VirtualAppliance and requires it there; catching it at validation time
  beats a mid-apply ARM error.
- **`address_prefix` is a free string**, because ARM accepts both CIDR
  blocks and Azure service tags ("AzureBackup") -- a CIDR-only validation
  would reject valid, useful configurations.

## Operational Behavior Worth Knowing

- Route changes apply immediately to every attached subnet -- a route
  table edit is a network-wide change, not a resource-local one.
- The most specific prefix wins; among equal prefixes, a user-defined
  route beats a system route beats a BGP-learned route.
- Disabling `bgp_route_propagation_enabled` is the standard hardening for
  forced tunneling: it stops learned on-premises routes from bypassing the
  user-defined default route.
- 0.0.0.0/0 via VirtualAppliance breaks Azure's default outbound access
  for the attached subnets -- intended (that is the point), but plan the
  appliance's own egress path.

## Composition

- `resource_group` → `AzureResourceGroup.status.outputs.resource_group_name`
- `route_table_id` output is the seam `AzureSubnet` consumes to attach the
  table (subnet-side field).
