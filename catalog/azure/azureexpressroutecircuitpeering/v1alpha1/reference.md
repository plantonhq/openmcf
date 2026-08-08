# AzureExpressRouteCircuitPeering

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureExpressRouteCircuitPeering
metadata:
  name: test-express-route-circuit-peering
spec:
  resourceGroup:
    value: test-rg
  # ARM addresses peerings as children of the circuit NAME.
  expressRouteCircuitName:
    value: hq-circuit
  # Private peering: VNet connectivity -- the type virtually every
  # deployment needs. The type is the ARM identity (one per circuit).
  peeringType: AZURE_PRIVATE_PEERING
  # Provider-assigned; unique on the circuit.
  vlanId: 100
  # One /30 per physical link: your router takes the first usable
  # address, Microsoft's edge the second.
  primaryPeerAddressPrefix: "192.168.16.0/30"
  secondaryPeerAddressPrefix: "192.168.16.4/30"
  peerAsn: 65010
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.expressRouteCircuitName` | `string \| valueFrom` | yes |  | AzureExpressRouteCircuit (`status.outputs.express_route_circuit_name`) |
| `spec.peeringType` | `enum` |  |  |  |
| `spec.vlanId` | `int32` |  |  |  |
| `spec.primaryPeerAddressPrefix` | `string` |  |  |  |
| `spec.secondaryPeerAddressPrefix` | `string` |  |  |  |
| `spec.ipv4Enabled` | `bool` |  | `true` |  |
| `spec.peerAsn` | `int64` |  |  |  |
| `spec.sharedKey` | `string \| valueFrom` (sensitive) |  |  |  |
| `spec.microsoftPeeringConfig` | `AzureExpressRouteCircuitPeeringMicrosoftConfig` |  |  |  |
| `spec.microsoftPeeringConfig.advertisedPublicPrefixes` | `[]string` | yes |  |  |
| `spec.microsoftPeeringConfig.customerAsn` | `int64` |  |  |  |
| `spec.microsoftPeeringConfig.routingRegistryName` | `string` |  | `NONE` |  |
| `spec.microsoftPeeringConfig.advertisedCommunities` | `[]string` |  |  |  |
| `spec.ipv6` | `AzureExpressRouteCircuitPeeringIpv6` |  |  |  |
| `spec.ipv6.primaryPeerAddressPrefix` | `string` | yes |  |  |
| `spec.ipv6.secondaryPeerAddressPrefix` | `string` | yes |  |  |
| `spec.ipv6.enabled` | `bool` |  | `true` |  |
| `spec.ipv6.routeFilterId` | `string` |  |  |  |
| `spec.ipv6.microsoftPeering` | `AzureExpressRouteCircuitPeeringMicrosoftConfig` |  |  |  |
| `spec.ipv6.microsoftPeering.advertisedPublicPrefixes` | `[]string` | yes |  |  |
| `spec.ipv6.microsoftPeering.customerAsn` | `int64` |  |  |  |
| `spec.ipv6.microsoftPeering.routingRegistryName` | `string` |  | `NONE` |  |
| `spec.ipv6.microsoftPeering.advertisedCommunities` | `[]string` |  |  |  |
| `spec.routeFilterId` | `string` |  |  |  |
| `spec.connections` | `[]AzureExpressRouteCircuitPeeringConnection` |  |  |  |
| `spec.connections[].name` | `string` | yes |  |  |
| `spec.connections[].peerPeeringId` | `string \| valueFrom` | yes |  | AzureExpressRouteCircuitPeering (`status.outputs.express_route_circuit_peering_id`) |
| `spec.connections[].addressPrefixIpv4` | `string` | yes |  |  |
| `spec.connections[].addressPrefixIpv6` | `string` |  |  |  |
| `spec.connections[].authorizationKey` | `string \| valueFrom` (sensitive) |  |  |  |

## Field Details

### spec.resourceGroup

`string | valueFrom` · required

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.expressRouteCircuitName

`string | valueFrom` · required

- references: AzureExpressRouteCircuit (`status.outputs.express_route_circuit_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureExpressRouteCircuit, name: <that resource's name>, fieldPath: status.outputs.express_route_circuit_name}} -- a bare string does not parse

### spec.peeringType

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_express_route_circuit_peering_type_unspecified`
- `AZURE_PRIVATE_PEERING`
- `AZURE_PUBLIC_PEERING`
- `MICROSOFT_PEERING`

### spec.vlanId

`int32`

- rule: {"int32":{"lte":4094,"gte":1}}

### spec.primaryPeerAddressPrefix

`string`

### spec.secondaryPeerAddressPrefix

`string`

### spec.ipv4Enabled

`bool` · optional (explicit presence)

- default: `true`

### spec.peerAsn

`int64`

- rule: {"int64":{"gte":"0"}}

### spec.sharedKey

`string | valueFrom` · sensitive

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.microsoftPeeringConfig

`AzureExpressRouteCircuitPeeringMicrosoftConfig`

### spec.microsoftPeeringConfig.advertisedPublicPrefixes

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.microsoftPeeringConfig.customerAsn

`int64`

- rule: {"int64":{"gte":"0"}}

### spec.microsoftPeeringConfig.routingRegistryName

`string` · optional (explicit presence)

- default: `NONE`

### spec.microsoftPeeringConfig.advertisedCommunities

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.ipv6

`AzureExpressRouteCircuitPeeringIpv6`

### spec.ipv6.primaryPeerAddressPrefix

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.ipv6.secondaryPeerAddressPrefix

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.ipv6.enabled

`bool` · optional (explicit presence)

- default: `true`

### spec.ipv6.routeFilterId

`string`

### spec.ipv6.microsoftPeering

`AzureExpressRouteCircuitPeeringMicrosoftConfig`

### spec.ipv6.microsoftPeering.advertisedPublicPrefixes

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.ipv6.microsoftPeering.customerAsn

`int64`

- rule: {"int64":{"gte":"0"}}

### spec.ipv6.microsoftPeering.routingRegistryName

`string` · optional (explicit presence)

- default: `NONE`

### spec.ipv6.microsoftPeering.advertisedCommunities

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.routeFilterId

`string`

### spec.connections

`[]AzureExpressRouteCircuitPeeringConnection`

### spec.connections[].name

`string` · required

- rule: Connection names start with a letter or number, end with a letter, number, or underscore, and may contain alphanumerics, underscores, periods, and hyphens
- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.connections[].peerPeeringId

`string | valueFrom` · required

- references: AzureExpressRouteCircuitPeering (`status.outputs.express_route_circuit_peering_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureExpressRouteCircuitPeering, name: <that resource's name>, fieldPath: status.outputs.express_route_circuit_peering_id}} -- a bare string does not parse

### spec.connections[].addressPrefixIpv4

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.connections[].addressPrefixIpv6

`string`

### spec.connections[].authorizationKey

`string | valueFrom` · sensitive

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

## Validation Rules

- `peering_type_required`: Choose the peering type explicitly -- AZURE_PRIVATE_PEERING for VNet connectivity (the common case) or MICROSOFT_PEERING for Microsoft public services
- `address_prefix_pair_travels_together`: primary_peer_address_prefix and secondary_peer_address_prefix are set together -- one /30 per physical link
- `route_filter_is_microsoft_peering_only`: route_filter_id applies to MICROSOFT_PEERING only -- private peering carries no service communities to filter
- `microsoft_config_is_microsoft_peering_only`: microsoft_peering_config applies to MICROSOFT_PEERING only
- `microsoft_ipv4_requires_config`: Microsoft peering with IPv4 prefixes needs microsoft_peering_config (the advertised-prefix contract)
- `microsoft_config_requires_ipv4_prefixes`: microsoft_peering_config configures the IPv4 session -- set primary_peer_address_prefix and secondary_peer_address_prefix with it
- `ipv6_not_on_public_peering`: IPv6 configuration applies to AZURE_PRIVATE_PEERING and MICROSOFT_PEERING only (public peering is deprecated and IPv4-only)
- `connections_are_private_peering_only`: Global Reach connections link PRIVATE peerings -- set peering_type to AZURE_PRIVATE_PEERING or remove connections
- `connection_names_unique`: Global Reach connection names must be unique on the peering

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureExpressRouteCircuitPeering, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.express_route_circuit_peering_id` | `string` |  |
| `status.outputs.azure_asn` | `int64` |  |
| `status.outputs.primary_azure_port` | `string` |  |
| `status.outputs.secondary_azure_port` | `string` |  |
| `status.outputs.connection_ids` | `map<string, string>` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.expressRouteCircuitName` | AzureExpressRouteCircuit | `status.outputs.express_route_circuit_name` |
| `spec.connections[].peerPeeringId` | AzureExpressRouteCircuitPeering | `status.outputs.express_route_circuit_peering_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureExpressRouteCircuitPeering | `spec.connections[].peerPeeringId` | `status.outputs.express_route_circuit_peering_id` |

## See Also

- [Overview](../README.md)
