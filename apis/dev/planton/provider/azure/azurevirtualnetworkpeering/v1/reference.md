# AzureVirtualNetworkPeering

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1`

**AzureVirtualNetworkPeeringSpec** defines the configuration for peering
one Azure virtual network to another: private, low-latency connectivity
over the Microsoft backbone, without gateways, public IPs, or encryption
overhead. Peered networks exchange traffic as if they were one network
while remaining separately owned and managed -- the building block of
hub-and-spoke topologies.

One resource models ONE DIRECTION of a peering, exactly as ARM does.
Connectivity only flows once BOTH directions exist, so a working pair is
two of these resources with local and remote swapped -- typically stamped
from the same chart. Cross-subscription and cross-region (global) peering
both work with plain ARM IDs; note that use_remote_gateways cannot be
combined with global peering.

## Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureVirtualNetworkPeering
metadata:
  name: test-peering
spec:
  name: hub-to-spoke1
  virtualNetworkId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet
  remoteVirtualNetworkId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/spoke1-vnet
  allowForwardedTraffic: true
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.name` | `string` | yes |  |  |
| `spec.virtualNetworkId` | `string \| valueFrom` | yes |  | AzureVirtualNetwork (`status.outputs.virtual_network_id`) |
| `spec.remoteVirtualNetworkId` | `string \| valueFrom` | yes |  | AzureVirtualNetwork (`status.outputs.virtual_network_id`) |
| `spec.allowVirtualNetworkAccess` | `bool` |  | `true` |  |
| `spec.allowForwardedTraffic` | `bool` |  | `false` |  |
| `spec.allowGatewayTransit` | `bool` |  | `false` |  |
| `spec.useRemoteGateways` | `bool` |  | `false` |  |
| `spec.peerCompleteVirtualNetworksEnabled` | `bool` |  | `true` |  |
| `spec.localSubnetNames` | `[]string` |  |  |  |
| `spec.remoteSubnetNames` | `[]string` |  |  |  |
| `spec.onlyIpv6PeeringEnabled` | `bool` |  | `false` |  |

## Field Details

### spec.name

`string` · required

The name of the peering, unique within the local virtual network.
1-80 characters (alphanumerics, underscores, periods, and hyphens;
must start with a letter or number and end with a letter, number, or
underscore). Name it after the far side ("spoke1-to-hub", "hub-to-spoke1")
so a network's peering list reads as a topology map. Changing the name
replaces the peering (a brief connectivity gap for this direction).

- rule: Peering names start with a letter or number, end with a letter, number, or underscore, and may contain alphanumerics, underscores, periods, and hyphens
- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.virtualNetworkId

`string | valueFrom` · required

The LOCAL virtual network -- the side this peering is written on -- by
ARM ID. The peering is an ARM child of this network: its ID carries
both the resource group and the network name, and the modules derive
both from it. Changing the local network replaces the peering.
Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/virtualNetworks/{name}

- references: AzureVirtualNetwork (`status.outputs.virtual_network_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVirtualNetwork, name: <that resource's name>, fieldPath: status.outputs.virtual_network_id}} -- a bare string does not parse

### spec.remoteVirtualNetworkId

`string | valueFrom` · required

The REMOTE virtual network to peer with, by ARM ID. Works across
subscriptions and across regions (global peering) unchanged -- the ID
carries everything Azure needs. Changing the remote network replaces
the peering.

- references: AzureVirtualNetwork (`status.outputs.virtual_network_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVirtualNetwork, name: <that resource's name>, fieldPath: status.outputs.virtual_network_id}} -- a bare string does not parse

### spec.allowVirtualNetworkAccess

`bool` · optional (explicit presence)

Whether traffic from the local network can reach the remote network.
Azure's default is true -- the reason peering exists. Setting false
keeps the peering established but blocks direct VM-to-VM traffic in
this direction (used in inspection topologies where traffic must flow
through an appliance instead). Updatable in place.

- default: `true`

### spec.allowForwardedTraffic

`bool` · optional (explicit presence)

Whether traffic FORWARDED by the remote network (originating outside
it, e.g. relayed by an NVA or VPN in the remote network) is accepted.
Azure's default is false. Spokes peered to a hub with a firewall or
VPN set this true on the hub-to-spoke direction so relayed traffic is
admitted. Updatable in place.

- default: `false`

### spec.allowGatewayTransit

`bool` · optional (explicit presence)

Whether the LOCAL network's VPN/ExpressRoute gateway may be used by
the remote network (gateway transit). Azure's default is false. Set
true on the HUB side of a hub-and-spoke so spokes can ride the hub's
gateway; pair it with use_remote_gateways=true on the spoke side.
Updatable in place.

- default: `false`

### spec.useRemoteGateways

`bool` · optional (explicit presence)

Whether the LOCAL network uses the REMOTE network's gateway for
transit (the spoke side of gateway transit; requires
allow_gateway_transit=true on the remote peering). Azure's default is
false. Only one peering per network may set this, the network must
have no gateway of its own, and it cannot be combined with global
(cross-region) peering. Updatable in place.

- default: `false`

### spec.peerCompleteVirtualNetworksEnabled

`bool` · optional (explicit presence)

Whether the peering spans the networks' COMPLETE address spaces.
Azure's default is true -- the standard full-network peering. Set
false for subnet-scoped peering, listing the specific subnets in
local_subnet_names/remote_subnet_names. Changing this replaces the
peering.

- default: `true`

### spec.localSubnetNames

`[]string`

For subnet-scoped peering (peer_complete_virtual_networks_enabled =
false): the LOCAL subnets included in the peering, by name. Updatable
in place.

### spec.remoteSubnetNames

`[]string`

For subnet-scoped peering (peer_complete_virtual_networks_enabled =
false): the REMOTE subnets included in the peering, by name. Updatable
in place.

### spec.onlyIpv6PeeringEnabled

`bool` · optional (explicit presence)

Whether only the networks' IPv6 address space is peered (for
subnet-scoped peering of dual-stack networks). Azure's default is
false. Changing this replaces the peering.

- default: `false`

## Validation Rules

- `subnet_names_require_subnet_scoped_peering`: local_subnet_names and remote_subnet_names apply only to subnet-scoped peering: set peer_complete_virtual_networks_enabled to false to use them

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureVirtualNetworkPeering, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.peering_id` | `string` | The Azure Resource Manager ID of the peering. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/virtualNetworks/{vnet}/virtualNetworkPeerings/{name} |
| `status.outputs.peering_name` | `string` | The name of the peering within its local virtual network. |
| `status.outputs.virtual_network_name` | `string` | The name of the LOCAL virtual network the peering is written on, derived from the referenced network's ARM ID -- exported so charts can compose sibling resources without re-parsing the ID. |
| `status.outputs.resource_group_name` | `string` | The name of the resource group the local network lives in, derived from the referenced network's ARM ID. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.virtualNetworkId` | AzureVirtualNetwork | `status.outputs.virtual_network_id` |
| `spec.remoteVirtualNetworkId` | AzureVirtualNetwork | `status.outputs.virtual_network_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
