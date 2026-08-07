# AzureVirtualNetworkPeering -- Design Research

## The Resource

An Azure virtual network peering (`Microsoft.Network/virtualNetworks/
virtualNetworkPeerings`) connects two virtual networks privately over the
Microsoft backbone. The component maps 1:1 onto
`azurerm_virtual_network_peering` (azurerm v4.x,
`internal/services/network/virtual_network_peering_resource.go`),
parity-verified against pulumi-azure v6 (`network.VirtualNetworkPeering`).

ARM models each direction as a separate child of the local network. One
Planton resource therefore models one direction -- connectivity only
flows once both directions exist.

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `name` | Required, ForceNew |
| `virtual_network_name` | -- | Derived from `virtual_network_id` ARM ID |
| `resource_group_name` | -- | Derived from `virtual_network_id` ARM ID |
| `remote_virtual_network_id` | `remote_virtual_network_id` | FK → AzureVirtualNetwork; ForceNew |
| `allow_virtual_network_access` | `allow_virtual_network_access` | azurerm default true → `optional bool` with default "true" |
| `allow_forwarded_traffic` | `allow_forwarded_traffic` | azurerm default false → `optional bool` with default "false" |
| `allow_gateway_transit` | `allow_gateway_transit` | azurerm default false → `optional bool` with default "false" |
| `use_remote_gateways` | `use_remote_gateways` | azurerm default false → `optional bool` with default "false" |
| `peer_complete_virtual_networks_enabled` | `peer_complete_virtual_networks_enabled` | azurerm default true → `optional bool` with default "true"; ForceNew |
| `local_subnet_names` | `local_subnet_names` | Meaningful only when complete-network peering is off |
| `remote_subnet_names` | `remote_subnet_names` | Meaningful only when complete-network peering is off |
| `only_ipv6_peering_enabled` | `only_ipv6_peering_enabled` | ForceNew; only sent when explicitly true |
| `tags` | -- | Not modeled; ARM peerings are not taggable |

## Decomposition Decisions

- **One direction per resource.** ARM stores each peering as a child of the
  local network with its own name and policy flags. Folding both directions
  into one resource would fight the provider model and hide the directional
  gateway-transit pairing (`allow_gateway_transit` on one side,
  `use_remote_gateways` on the other).
- **Local placement derived from `virtual_network_id`.** The peering's
  resource group and local network name are parsed from the referenced
  network's ARM ID. Restating them in the spec would invite drift against
  the network resource the chart already owns.
- **Both network references are FKs defaulting to `AzureVirtualNetwork`.**
  The remote side is still a first-class reference even when it lives in
  another subscription or region -- the ARM ID carries everything Azure
  needs for global and cross-subscription peering.
- **Subnet names validated at spec level.** `local_subnet_names` and
  `remote_subnet_names` are rejected unless
  `peer_complete_virtual_networks_enabled` is false, matching when ARM
  actually consumes them.

## Design Decisions

- **No `region` or `resource_group` fields.** Unlike regional standalone
  resources, a peering's placement is fully determined by its parent
  network. The outputs export the derived `virtual_network_name` and
  `resource_group_name` so charts can compose siblings without re-parsing
  the ARM ID.
- **Connectivity flags as optional bools with proto defaults** mirroring
  Azure's defaults, so both IaC engines always send the same effective
  values whether the user omits a flag or sets it explicitly.
- **`only_ipv6_peering_enabled` sent only when true.** ARM treats it as a
  creation-time property; the modules omit it when false to avoid
  unintended ForceNew behavior on engines that distinguish null from false.

## Operational Behavior Worth Knowing

- The reciprocal peering can deploy concurrently; Azure retries internally
  while the far side catches up. A one-direction deploy is valid for
  provisioning proofs even though traffic needs both sides.
- `allow_forwarded_traffic` on the hub-to-spoke direction is what admits
  traffic relayed by an NVA or VPN gateway sitting in the hub -- without
  it, spokes can address hub subnets directly but not traffic forwarded
  through the hub.
- Gateway transit is directional: `allow_gateway_transit` on the hub
  peering plus `use_remote_gateways` on the spoke peering. Only one
  peering per network may set `use_remote_gateways`, the local network
  must have no gateway of its own, and the pair cannot be global
  (cross-region).
- Subnet-scoped peering changes which address ranges are reachable, not
  which NSG or route table applies -- those remain properties of the
  subnets themselves.
- Replacing a peering (name or network change) causes a brief connectivity
  gap for that direction only; the reciprocal direction stays up until it
  is replaced too.

## Composition

- `virtual_network_id` → `AzureVirtualNetwork.status.outputs.virtual_network_id`
- `remote_virtual_network_id` → `AzureVirtualNetwork.status.outputs.virtual_network_id`
- Hub-and-spoke charts stamp two resources per spoke: hub→spoke and
  spoke→hub, with direction-appropriate `allow_forwarded_traffic`,
  `allow_gateway_transit`, and `use_remote_gateways` values
- `peering_id` / `peering_name` identify this direction; `virtual_network_name`
  and `resource_group_name` are convenience exports for sibling composition
