# AzureVirtualNetworkGatewayConnection

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureVirtualNetworkGatewayConnection
metadata:
  name: test-gateway-connection
spec:
  region: eastus
  resourceGroup:
    value: test-rg
  name: hq-to-azure
  # IPSEC (site-to-site) joins the gateway to an on-premises device
  # described by an AzureLocalNetworkGateway.
  type: IPSEC
  virtualNetworkGatewayId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworkGateways/hub-vpn-gateway
  localNetworkGatewayId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/localNetworkGateways/hq-datacenter
  # The pre-shared key both tunnel ends must agree on. Reference a secret
  # in real manifests; omit to let Azure generate one.
  sharedKey:
    value: replace-with-a-strong-pre-shared-key
  dpdTimeoutSeconds: 45
  # A pinned IPsec/IKE proposal for devices that need exact algorithms;
  # omit to use Azure's default proposal set.
  ipsecPolicy:
    dhGroup: DHGroup14
    ikeEncryption: AES256
    ikeIntegrity: SHA256
    ipsecEncryption: AES256
    ipsecIntegrity: SHA256
    pfsGroup: PFS2048
  tags:
    purpose: hack-test
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.type` | `enum` |  |  |  |
| `spec.virtualNetworkGatewayId` | `string \| valueFrom` | yes |  | AzureVirtualNetworkGateway (`status.outputs.virtual_network_gateway_id`) |
| `spec.localNetworkGatewayId` | `string \| valueFrom` |  |  | AzureLocalNetworkGateway (`status.outputs.local_network_gateway_id`) |
| `spec.peerVirtualNetworkGatewayId` | `string \| valueFrom` |  |  | AzureVirtualNetworkGateway (`status.outputs.virtual_network_gateway_id`) |
| `spec.expressRouteCircuitId` | `string \| valueFrom` |  |  |  |
| `spec.sharedKey` | `string \| valueFrom` (sensitive) |  |  |  |
| `spec.authorizationKey` | `string \| valueFrom` (sensitive) |  |  |  |
| `spec.bgpEnabled` | `bool` |  |  |  |
| `spec.customBgpAddresses` | `AzureVirtualNetworkGatewayConnectionCustomBgpAddresses` |  |  |  |
| `spec.customBgpAddresses.primary` | `string` | yes |  |  |
| `spec.customBgpAddresses.secondary` | `string` |  |  |  |
| `spec.dpdTimeoutSeconds` | `int32` |  |  |  |
| `spec.connectionProtocol` | `enum` |  |  |  |
| `spec.connectionMode` | `enum` |  |  |  |
| `spec.routingWeight` | `int32` |  |  |  |
| `spec.egressNatRuleIds` | `[]string \| valueFrom` |  |  |  |
| `spec.ingressNatRuleIds` | `[]string \| valueFrom` |  |  |  |
| `spec.usePolicyBasedTrafficSelectors` | `bool` |  |  |  |
| `spec.expressRouteGatewayBypass` | `bool` |  |  |  |
| `spec.privateLinkFastPathEnabled` | `bool` |  |  |  |
| `spec.localAzureIpAddressEnabled` | `bool` |  |  |  |
| `spec.trafficSelectorPolicies` | `[]AzureVirtualNetworkGatewayConnectionTrafficSelectorPolicy` |  |  |  |
| `spec.trafficSelectorPolicies[].localAddressCidrs` | `[]string` | yes |  |  |
| `spec.trafficSelectorPolicies[].remoteAddressCidrs` | `[]string` | yes |  |  |
| `spec.ipsecPolicy` | `AzureVirtualNetworkGatewayConnectionIpsecPolicy` |  |  |  |
| `spec.ipsecPolicy.dhGroup` | `string` |  |  |  |
| `spec.ipsecPolicy.ikeEncryption` | `string` |  |  |  |
| `spec.ipsecPolicy.ikeIntegrity` | `string` |  |  |  |
| `spec.ipsecPolicy.ipsecEncryption` | `string` |  |  |  |
| `spec.ipsecPolicy.ipsecIntegrity` | `string` |  |  |  |
| `spec.ipsecPolicy.pfsGroup` | `string` |  |  |  |
| `spec.ipsecPolicy.saDatasize` | `int32` |  |  |  |
| `spec.ipsecPolicy.saLifetime` | `int32` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.resourceGroup

`string | valueFrom` · required

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.type

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_virtual_network_gateway_connection_type_unspecified`
- `IPSEC`
- `VNET_TO_VNET`
- `EXPRESS_ROUTE`

### spec.virtualNetworkGatewayId

`string | valueFrom` · required

- references: AzureVirtualNetworkGateway (`status.outputs.virtual_network_gateway_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVirtualNetworkGateway, name: <that resource's name>, fieldPath: status.outputs.virtual_network_gateway_id}} -- a bare string does not parse

### spec.localNetworkGatewayId

`string | valueFrom`

- references: AzureLocalNetworkGateway (`status.outputs.local_network_gateway_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureLocalNetworkGateway, name: <that resource's name>, fieldPath: status.outputs.local_network_gateway_id}} -- a bare string does not parse

### spec.peerVirtualNetworkGatewayId

`string | valueFrom`

- references: AzureVirtualNetworkGateway (`status.outputs.virtual_network_gateway_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVirtualNetworkGateway, name: <that resource's name>, fieldPath: status.outputs.virtual_network_gateway_id}} -- a bare string does not parse

### spec.expressRouteCircuitId

`string | valueFrom`

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.sharedKey

`string | valueFrom` · sensitive

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.authorizationKey

`string | valueFrom` · sensitive

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.bgpEnabled

`bool`

### spec.customBgpAddresses

`AzureVirtualNetworkGatewayConnectionCustomBgpAddresses`

### spec.customBgpAddresses.primary

`string` · required

- rule: {"required":true,"string":{"ipv4":true}}

### spec.customBgpAddresses.secondary

`string`

- rule: secondary must be an IPv4 address

### spec.dpdTimeoutSeconds

`int32` · optional (explicit presence)

- rule: {"int32":{"gte":9}}

### spec.connectionProtocol

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_virtual_network_gateway_connection_protocol_unspecified`
- `IKE_V1`
- `IKE_V2`

### spec.connectionMode

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_virtual_network_gateway_connection_mode_unspecified`
- `DEFAULT`
- `INITIATOR_ONLY`
- `RESPONDER_ONLY`

### spec.routingWeight

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":32000,"gte":0}}

### spec.egressNatRuleIds

`[]string | valueFrom`

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.ingressNatRuleIds

`[]string | valueFrom`

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.usePolicyBasedTrafficSelectors

`bool`

### spec.expressRouteGatewayBypass

`bool`

### spec.privateLinkFastPathEnabled

`bool`

### spec.localAzureIpAddressEnabled

`bool`

### spec.trafficSelectorPolicies

`[]AzureVirtualNetworkGatewayConnectionTrafficSelectorPolicy`

### spec.trafficSelectorPolicies[].localAddressCidrs

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.trafficSelectorPolicies[].remoteAddressCidrs

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.ipsecPolicy

`AzureVirtualNetworkGatewayConnectionIpsecPolicy`

### spec.ipsecPolicy.dhGroup

`string`

- rule: {"string":{"in":["DHGroup1","DHGroup14","DHGroup2","DHGroup2048","DHGroup24","ECP256","ECP384","None"]}}

### spec.ipsecPolicy.ikeEncryption

`string`

- rule: {"string":{"in":["AES128","AES192","AES256","DES","DES3","GCMAES128","GCMAES256"]}}

### spec.ipsecPolicy.ikeIntegrity

`string`

- rule: {"string":{"in":["GCMAES128","GCMAES256","MD5","SHA1","SHA256","SHA384"]}}

### spec.ipsecPolicy.ipsecEncryption

`string`

- rule: {"string":{"in":["AES128","AES192","AES256","DES","DES3","GCMAES128","GCMAES192","GCMAES256","None"]}}

### spec.ipsecPolicy.ipsecIntegrity

`string`

- rule: {"string":{"in":["GCMAES128","GCMAES192","GCMAES256","MD5","SHA1","SHA256"]}}

### spec.ipsecPolicy.pfsGroup

`string`

- rule: {"string":{"in":["ECP256","ECP384","None","PFS1","PFS14","PFS2","PFS2048","PFS24","PFSMM"]}}

### spec.ipsecPolicy.saDatasize

`int32` · optional (explicit presence)

- rule: {"int32":{"gte":1024}}

### spec.ipsecPolicy.saLifetime

`int32` · optional (explicit presence)

- rule: {"int32":{"gte":300}}

### spec.tags

`map<string, string>`

## Validation Rules

- `type_required`: Choose the connection type explicitly: IPSEC (to an on-premises device), VNET_TO_VNET (to another gateway), or EXPRESS_ROUTE (to a circuit)
- `ipsec_requires_local_network_gateway`: An IPSEC (site-to-site) connection needs local_network_gateway_id -- the AzureLocalNetworkGateway describing the on-premises side
- `vnet_to_vnet_requires_peer_gateway`: A VNET_TO_VNET connection needs peer_virtual_network_gateway_id (and the far gateway needs a mirror connection with the same shared_key)
- `express_route_requires_circuit`: An EXPRESS_ROUTE connection needs express_route_circuit_id
- `custom_bgp_is_ipsec_only`: custom_bgp_addresses applies to IPSEC connections only
- `custom_bgp_requires_bgp`: custom_bgp_addresses requires bgp_enabled
- `private_link_fast_path_requires_bypass`: private_link_fast_path_enabled requires express_route_gateway_bypass (Azure's FastPath contract)
- `policy_based_selectors_require_ipsec_policy`: use_policy_based_traffic_selectors requires a custom ipsec_policy (Azure rejects the flag without one)
- `shared_key_not_for_express_route`: ExpressRoute connections carry no pre-shared key -- remove shared_key (use authorization_key for cross-subscription circuits)
- `authorization_key_is_express_route_only`: authorization_key authorizes an ExpressRoute circuit in another subscription -- it has no meaning on IPSEC or VNET_TO_VNET connections

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureVirtualNetworkGatewayConnection, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.connection_id` | `string` |  |
| `status.outputs.connection_name` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.virtualNetworkGatewayId` | AzureVirtualNetworkGateway | `status.outputs.virtual_network_gateway_id` |
| `spec.localNetworkGatewayId` | AzureLocalNetworkGateway | `status.outputs.local_network_gateway_id` |
| `spec.peerVirtualNetworkGatewayId` | AzureVirtualNetworkGateway | `status.outputs.virtual_network_gateway_id` |

## See Also

- [Overview](../README.md)
