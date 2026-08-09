# AzureVpnGatewayConnection

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureVpnGatewayConnection
metadata:
  name: test-vpn-connection
spec:
  name: branch-london
  vpnGatewayId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/vpnGateways/hub-vpn-gateway
  remoteVpnSiteId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/vpnSites/branch-london
  internetSecurityEnabled: false
  routing:
    associatedRouteTableId:
      value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/virtualHubs/hub-eastus/hubRouteTables/defaultRouteTable
    propagatedRouteTable:
      routeTableIds:
        - value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/virtualHubs/hub-eastus/hubRouteTables/defaultRouteTable
      labels:
        - default
  vpnLinks:
    - name: primary-isp
      vpnSiteLinkId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/vpnSites/branch-london/vpnSiteLinks/primary-isp
      bandwidthMbps: 50
      dpdTimeoutSeconds: 45
      sharedKey:
        value: test-pre-shared-key
      ipsecPolicies:
        - saLifetimeSec: 3600
          saDataSizeKb: 102400000
          encryptionAlgorithm: AES256
          integrityAlgorithm: SHA256
          ikeEncryptionAlgorithm: AES256
          ikeIntegrityAlgorithm: SHA256
          dhGroup: DHGroup14
          pfsGroup: PFS2048
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.name` | `string` | yes |  |  |
| `spec.vpnGatewayId` | `string \| valueFrom` | yes |  | AzureVpnGateway (`status.outputs.vpn_gateway_id`) |
| `spec.remoteVpnSiteId` | `string \| valueFrom` | yes |  | AzureVpnSite (`status.outputs.vpn_site_id`) |
| `spec.internetSecurityEnabled` | `bool` |  |  |  |
| `spec.routing` | `AzureVpnGatewayConnectionRouting` |  |  |  |
| `spec.routing.associatedRouteTableId` | `string \| valueFrom` | yes |  | AzureVirtualHub (`status.outputs.default_route_table_id`) |
| `spec.routing.inboundRouteMapId` | `string \| valueFrom` |  |  | AzureVirtualHub (`status.outputs.route_map_ids`) |
| `spec.routing.outboundRouteMapId` | `string \| valueFrom` |  |  | AzureVirtualHub (`status.outputs.route_map_ids`) |
| `spec.routing.propagatedRouteTable` | `AzureVpnGatewayConnectionPropagatedRouteTable` |  |  |  |
| `spec.routing.propagatedRouteTable.routeTableIds` | `[]string \| valueFrom` | yes |  | AzureVirtualHub (`status.outputs.default_route_table_id`) |
| `spec.routing.propagatedRouteTable.labels` | `[]string` |  |  |  |
| `spec.vpnLinks` | `[]AzureVpnGatewayConnectionLink` | yes |  |  |
| `spec.vpnLinks[].name` | `string` | yes |  |  |
| `spec.vpnLinks[].vpnSiteLinkId` | `string \| valueFrom` | yes |  | AzureVpnSite (`status.outputs.link_ids`) |
| `spec.vpnLinks[].bandwidthMbps` | `int32` |  | `10` |  |
| `spec.vpnLinks[].protocol` | `enum` |  |  |  |
| `spec.vpnLinks[].connectionMode` | `enum` |  |  |  |
| `spec.vpnLinks[].routeWeight` | `int32` |  |  |  |
| `spec.vpnLinks[].dpdTimeoutSeconds` | `int32` |  |  |  |
| `spec.vpnLinks[].sharedKey` | `string \| valueFrom` (sensitive) |  |  |  |
| `spec.vpnLinks[].bgpEnabled` | `bool` |  |  |  |
| `spec.vpnLinks[].ratelimitEnabled` | `bool` |  |  |  |
| `spec.vpnLinks[].localAzureIpAddressEnabled` | `bool` |  |  |  |
| `spec.vpnLinks[].policyBasedTrafficSelectorEnabled` | `bool` |  |  |  |
| `spec.vpnLinks[].egressNatRuleIds` | `[]string \| valueFrom` |  |  | AzureVpnGateway (`status.outputs.nat_rule_ids`) |
| `spec.vpnLinks[].ingressNatRuleIds` | `[]string \| valueFrom` |  |  | AzureVpnGateway (`status.outputs.nat_rule_ids`) |
| `spec.vpnLinks[].ipsecPolicies` | `[]AzureVpnGatewayConnectionIpsecPolicy` |  |  |  |
| `spec.vpnLinks[].ipsecPolicies[].saLifetimeSec` | `int32` |  |  |  |
| `spec.vpnLinks[].ipsecPolicies[].saDataSizeKb` | `int32` |  |  |  |
| `spec.vpnLinks[].ipsecPolicies[].encryptionAlgorithm` | `string` |  |  |  |
| `spec.vpnLinks[].ipsecPolicies[].integrityAlgorithm` | `string` |  |  |  |
| `spec.vpnLinks[].ipsecPolicies[].ikeEncryptionAlgorithm` | `string` |  |  |  |
| `spec.vpnLinks[].ipsecPolicies[].ikeIntegrityAlgorithm` | `string` |  |  |  |
| `spec.vpnLinks[].ipsecPolicies[].dhGroup` | `string` |  |  |  |
| `spec.vpnLinks[].ipsecPolicies[].pfsGroup` | `string` |  |  |  |
| `spec.vpnLinks[].customBgpAddresses` | `[]AzureVpnGatewayConnectionCustomBgpAddress` |  |  |  |
| `spec.vpnLinks[].customBgpAddresses[].ipAddress` | `string` | yes |  |  |
| `spec.vpnLinks[].customBgpAddresses[].ipConfigurationId` | `string` |  |  |  |
| `spec.trafficSelectorPolicies` | `[]AzureVpnGatewayConnectionTrafficSelectorPolicy` |  |  |  |
| `spec.trafficSelectorPolicies[].localAddressCidrs` | `[]string` | yes |  |  |
| `spec.trafficSelectorPolicies[].remoteAddressCidrs` | `[]string` | yes |  |  |

## Field Details

### spec.name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vpnGatewayId

`string | valueFrom` · required

- references: AzureVpnGateway (`status.outputs.vpn_gateway_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVpnGateway, name: <that resource's name>, fieldPath: status.outputs.vpn_gateway_id}} -- a bare string does not parse

### spec.remoteVpnSiteId

`string | valueFrom` · required

- references: AzureVpnSite (`status.outputs.vpn_site_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVpnSite, name: <that resource's name>, fieldPath: status.outputs.vpn_site_id}} -- a bare string does not parse

### spec.internetSecurityEnabled

`bool`

### spec.routing

`AzureVpnGatewayConnectionRouting`

### spec.routing.associatedRouteTableId

`string | valueFrom` · required

- references: AzureVirtualHub (`status.outputs.default_route_table_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVirtualHub, name: <that resource's name>, fieldPath: status.outputs.default_route_table_id}} -- a bare string does not parse

### spec.routing.inboundRouteMapId

`string | valueFrom`

- references: AzureVirtualHub (`status.outputs.route_map_ids`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVirtualHub, name: <that resource's name>, fieldPath: status.outputs.route_map_ids}} -- a bare string does not parse

### spec.routing.outboundRouteMapId

`string | valueFrom`

- references: AzureVirtualHub (`status.outputs.route_map_ids`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVirtualHub, name: <that resource's name>, fieldPath: status.outputs.route_map_ids}} -- a bare string does not parse

### spec.routing.propagatedRouteTable

`AzureVpnGatewayConnectionPropagatedRouteTable`

### spec.routing.propagatedRouteTable.routeTableIds

`[]string | valueFrom` · required

- references: AzureVirtualHub (`status.outputs.default_route_table_id`)
- rule: {"repeated":{"minItems":"1"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVirtualHub, name: <that resource's name>, fieldPath: status.outputs.default_route_table_id}} -- a bare string does not parse

### spec.routing.propagatedRouteTable.labels

`[]string`

### spec.vpnLinks

`[]AzureVpnGatewayConnectionLink` · required

- rule: {"repeated":{"minItems":"1"}}
- rule: policy_based_traffic_selector_enabled requires a custom ipsec_policies entry (Azure rejects the flag without one)

### spec.vpnLinks[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vpnLinks[].vpnSiteLinkId

`string | valueFrom` · required

- references: AzureVpnSite (`status.outputs.link_ids`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVpnSite, name: <that resource's name>, fieldPath: status.outputs.link_ids}} -- a bare string does not parse

### spec.vpnLinks[].bandwidthMbps

`int32` · optional (explicit presence)

- default: `10`
- rule: {"int32":{"gte":1}}

### spec.vpnLinks[].protocol

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_vpn_gateway_connection_protocol_unspecified`
- `IKE_V1`
- `IKE_V2`

### spec.vpnLinks[].connectionMode

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_vpn_gateway_connection_mode_unspecified`
- `DEFAULT`
- `INITIATOR_ONLY`
- `RESPONDER_ONLY`

### spec.vpnLinks[].routeWeight

`int32`

- rule: {"int32":{"gte":0}}

### spec.vpnLinks[].dpdTimeoutSeconds

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":3600,"gte":9}}

### spec.vpnLinks[].sharedKey

`string | valueFrom` · sensitive

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.vpnLinks[].bgpEnabled

`bool`

### spec.vpnLinks[].ratelimitEnabled

`bool`

### spec.vpnLinks[].localAzureIpAddressEnabled

`bool`

### spec.vpnLinks[].policyBasedTrafficSelectorEnabled

`bool`

### spec.vpnLinks[].egressNatRuleIds

`[]string | valueFrom`

- references: AzureVpnGateway (`status.outputs.nat_rule_ids`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVpnGateway, name: <that resource's name>, fieldPath: status.outputs.nat_rule_ids}} -- a bare string does not parse

### spec.vpnLinks[].ingressNatRuleIds

`[]string | valueFrom`

- references: AzureVpnGateway (`status.outputs.nat_rule_ids`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVpnGateway, name: <that resource's name>, fieldPath: status.outputs.nat_rule_ids}} -- a bare string does not parse

### spec.vpnLinks[].ipsecPolicies

`[]AzureVpnGatewayConnectionIpsecPolicy`

### spec.vpnLinks[].ipsecPolicies[].saLifetimeSec

`int32`

- rule: {"int32":{"lte":172799,"gte":300}}

### spec.vpnLinks[].ipsecPolicies[].saDataSizeKb

`int32`

- rule: {"int32":{"gte":0}}

### spec.vpnLinks[].ipsecPolicies[].encryptionAlgorithm

`string`

- rule: {"string":{"in":["AES128","AES192","AES256","DES","DES3","GCMAES128","GCMAES192","GCMAES256","None"]}}

### spec.vpnLinks[].ipsecPolicies[].integrityAlgorithm

`string`

- rule: {"string":{"in":["GCMAES128","GCMAES192","GCMAES256","MD5","SHA1","SHA256"]}}

### spec.vpnLinks[].ipsecPolicies[].ikeEncryptionAlgorithm

`string`

- rule: {"string":{"in":["AES128","AES192","AES256","DES","DES3","GCMAES128","GCMAES256"]}}

### spec.vpnLinks[].ipsecPolicies[].ikeIntegrityAlgorithm

`string`

- rule: {"string":{"in":["GCMAES128","GCMAES256","MD5","SHA1","SHA256","SHA384"]}}

### spec.vpnLinks[].ipsecPolicies[].dhGroup

`string`

- rule: {"string":{"in":["DHGroup1","DHGroup14","DHGroup2","DHGroup2048","DHGroup24","ECP256","ECP384","None"]}}

### spec.vpnLinks[].ipsecPolicies[].pfsGroup

`string`

- rule: {"string":{"in":["ECP256","ECP384","None","PFS1","PFS14","PFS2","PFS2048","PFS24","PFSMM"]}}

### spec.vpnLinks[].customBgpAddresses

`[]AzureVpnGatewayConnectionCustomBgpAddress`

### spec.vpnLinks[].customBgpAddresses[].ipAddress

`string` · required

- rule: {"required":true,"string":{"ipv4":true}}

### spec.vpnLinks[].customBgpAddresses[].ipConfigurationId

`string`

- rule: {"string":{"in":["Instance0","Instance1"]}}

### spec.trafficSelectorPolicies

`[]AzureVpnGatewayConnectionTrafficSelectorPolicy`

### spec.trafficSelectorPolicies[].localAddressCidrs

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.trafficSelectorPolicies[].remoteAddressCidrs

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

## Validation Rules

- `vpn_link_names_unique`: vpn_links names must be unique on the connection

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureVpnGatewayConnection, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.connection_id` | `string` |  |
| `status.outputs.connection_name` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vpnGatewayId` | AzureVpnGateway | `status.outputs.vpn_gateway_id` |
| `spec.remoteVpnSiteId` | AzureVpnSite | `status.outputs.vpn_site_id` |
| `spec.routing.associatedRouteTableId` | AzureVirtualHub | `status.outputs.default_route_table_id` |
| `spec.routing.inboundRouteMapId` | AzureVirtualHub | `status.outputs.route_map_ids` |
| `spec.routing.outboundRouteMapId` | AzureVirtualHub | `status.outputs.route_map_ids` |
| `spec.routing.propagatedRouteTable.routeTableIds` | AzureVirtualHub | `status.outputs.default_route_table_id` |
| `spec.vpnLinks[].vpnSiteLinkId` | AzureVpnSite | `status.outputs.link_ids` |
| `spec.vpnLinks[].egressNatRuleIds` | AzureVpnGateway | `status.outputs.nat_rule_ids` |
| `spec.vpnLinks[].ingressNatRuleIds` | AzureVpnGateway | `status.outputs.nat_rule_ids` |

## See Also

- [Overview](../README.md)
