# AzureVirtualNetworkGateway

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureVirtualNetworkGateway
metadata:
  name: test-virtual-network-gateway
spec:
  region: eastus
  resourceGroup:
    value: test-rg
  name: hub-vpn-gateway
  # VPN (site-to-site/point-to-site) is the default type; RouteBased the
  # default routing model. VPN_GW_1 is the production entry point.
  sku: VPN_GW_1
  ipConfigurations:
    # The subnet's ARM name must be EXACTLY "GatewaySubnet" (/27 or
    # larger recommended); VPN gateways require a Standard static public
    # IP on every configuration.
    - subnetId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/hub-vnet/subnets/GatewaySubnet
      publicIpAddressId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/publicIPAddresses/vpn-gw-pip
  bgpEnabled: true
  bgpSettings:
    asn: 65515
  # A NAT rule translating overlapping on-premises space; connections opt
  # in via their egress/ingress NAT rule id lists (the rule's ARM id
  # surfaces in the nat_rule_ids output under its name).
  natRules:
    - name: egress-overlap
      externalMappings:
        - addressSpace: "100.64.1.0/24"
      internalMappings:
        - addressSpace: "10.0.1.0/24"
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
| `spec.vpnType` | `enum` |  |  |  |
| `spec.sku` | `enum` |  |  |  |
| `spec.generation` | `enum` |  |  |  |
| `spec.ipConfigurations` | `[]AzureVirtualNetworkGatewayIpConfiguration` | yes |  |  |
| `spec.ipConfigurations[].name` | `string` |  |  |  |
| `spec.ipConfigurations[].subnetId` | `string \| valueFrom` | yes |  | AzureSubnet (`status.outputs.subnet_id`) |
| `spec.ipConfigurations[].publicIpAddressId` | `string \| valueFrom` |  |  | AzurePublicIp (`status.outputs.public_ip_id`) |
| `spec.ipConfigurations[].privateIpAddressAllocation` | `enum` |  |  |  |
| `spec.activeActive` | `bool` |  |  |  |
| `spec.privateIpAddressEnabled` | `bool` |  |  |  |
| `spec.edgeZone` | `string` |  |  |  |
| `spec.bgpEnabled` | `bool` |  |  |  |
| `spec.bgpSettings` | `AzureVirtualNetworkGatewayBgpSettings` |  |  |  |
| `spec.bgpSettings.asn` | `int64` |  |  |  |
| `spec.bgpSettings.peerWeight` | `int32` |  |  |  |
| `spec.bgpSettings.peeringAddresses` | `[]AzureVirtualNetworkGatewayBgpPeeringAddress` |  |  |  |
| `spec.bgpSettings.peeringAddresses[].ipConfigurationName` | `string` |  |  |  |
| `spec.bgpSettings.peeringAddresses[].apipaAddresses` | `[]string` | yes |  |  |
| `spec.customRouteAddressPrefixes` | `[]string` |  |  |  |
| `spec.defaultLocalNetworkGatewayId` | `string \| valueFrom` |  |  | AzureLocalNetworkGateway (`status.outputs.local_network_gateway_id`) |
| `spec.vpnClientConfiguration` | `AzureVirtualNetworkGatewayVpnClientConfiguration` |  |  |  |
| `spec.vpnClientConfiguration.addressSpaces` | `[]string` | yes |  |  |
| `spec.vpnClientConfiguration.aadTenant` | `string` |  |  |  |
| `spec.vpnClientConfiguration.aadAudience` | `string` |  |  |  |
| `spec.vpnClientConfiguration.aadIssuer` | `string` |  |  |  |
| `spec.vpnClientConfiguration.rootCertificates` | `[]AzureVirtualNetworkGatewayVpnClientRootCertificate` |  |  |  |
| `spec.vpnClientConfiguration.rootCertificates[].name` | `string` | yes |  |  |
| `spec.vpnClientConfiguration.rootCertificates[].publicCertData` | `string` | yes |  |  |
| `spec.vpnClientConfiguration.revokedCertificates` | `[]AzureVirtualNetworkGatewayVpnClientRevokedCertificate` |  |  |  |
| `spec.vpnClientConfiguration.revokedCertificates[].name` | `string` | yes |  |  |
| `spec.vpnClientConfiguration.revokedCertificates[].thumbprint` | `string` | yes |  |  |
| `spec.vpnClientConfiguration.radiusServerAddress` | `string` |  |  |  |
| `spec.vpnClientConfiguration.radiusServerSecret` | `string \| valueFrom` (sensitive) |  |  |  |
| `spec.vpnClientConfiguration.radiusServers` | `[]AzureVirtualNetworkGatewayVpnClientRadiusServer` |  |  |  |
| `spec.vpnClientConfiguration.radiusServers[].address` | `string` | yes |  |  |
| `spec.vpnClientConfiguration.radiusServers[].secret` | `string \| valueFrom` (sensitive) | yes |  |  |
| `spec.vpnClientConfiguration.radiusServers[].score` | `int32` |  |  |  |
| `spec.vpnClientConfiguration.ipsecPolicy` | `AzureVirtualNetworkGatewayVpnClientIpsecPolicy` |  |  |  |
| `spec.vpnClientConfiguration.ipsecPolicy.dhGroup` | `string` |  |  |  |
| `spec.vpnClientConfiguration.ipsecPolicy.ikeEncryption` | `string` |  |  |  |
| `spec.vpnClientConfiguration.ipsecPolicy.ikeIntegrity` | `string` |  |  |  |
| `spec.vpnClientConfiguration.ipsecPolicy.ipsecEncryption` | `string` |  |  |  |
| `spec.vpnClientConfiguration.ipsecPolicy.ipsecIntegrity` | `string` |  |  |  |
| `spec.vpnClientConfiguration.ipsecPolicy.pfsGroup` | `string` |  |  |  |
| `spec.vpnClientConfiguration.ipsecPolicy.saLifetimeSeconds` | `int32` |  |  |  |
| `spec.vpnClientConfiguration.ipsecPolicy.saDataSizeKilobytes` | `int32` |  |  |  |
| `spec.vpnClientConfiguration.vpnClientProtocols` | `[]string` |  |  |  |
| `spec.vpnClientConfiguration.vpnAuthTypes` | `[]string` |  |  |  |
| `spec.vpnClientConfiguration.clientConnections` | `[]AzureVirtualNetworkGatewayClientConnection` |  |  |  |
| `spec.vpnClientConfiguration.clientConnections[].name` | `string` | yes |  |  |
| `spec.vpnClientConfiguration.clientConnections[].policyGroupNames` | `[]string` | yes |  |  |
| `spec.vpnClientConfiguration.clientConnections[].addressPrefixes` | `[]string` | yes |  |  |
| `spec.policyGroups` | `[]AzureVirtualNetworkGatewayPolicyGroup` |  |  |  |
| `spec.policyGroups[].name` | `string` | yes |  |  |
| `spec.policyGroups[].policyMembers` | `[]AzureVirtualNetworkGatewayPolicyMember` | yes |  |  |
| `spec.policyGroups[].policyMembers[].name` | `string` | yes |  |  |
| `spec.policyGroups[].policyMembers[].type` | `string` |  |  |  |
| `spec.policyGroups[].policyMembers[].value` | `string` | yes |  |  |
| `spec.policyGroups[].isDefault` | `bool` |  |  |  |
| `spec.policyGroups[].priority` | `int32` |  |  |  |
| `spec.bgpRouteTranslationForNatEnabled` | `bool` |  |  |  |
| `spec.dnsForwardingEnabled` | `bool` |  |  |  |
| `spec.ipSecReplayProtectionEnabled` | `bool` |  | `true` |  |
| `spec.minimumScaleUnit` | `int32` |  |  |  |
| `spec.maximumScaleUnit` | `int32` |  |  |  |
| `spec.remoteVnetTrafficEnabled` | `bool` |  |  |  |
| `spec.virtualWanTrafficEnabled` | `bool` |  |  |  |
| `spec.natRules` | `[]AzureVirtualNetworkGatewayNatRule` |  |  |  |
| `spec.natRules[].name` | `string` | yes |  |  |
| `spec.natRules[].mode` | `enum` |  |  |  |
| `spec.natRules[].type` | `enum` |  |  |  |
| `spec.natRules[].externalMappings` | `[]AzureVirtualNetworkGatewayNatRuleMapping` | yes |  |  |
| `spec.natRules[].externalMappings[].addressSpace` | `string` | yes |  |  |
| `spec.natRules[].externalMappings[].portRange` | `string` |  |  |  |
| `spec.natRules[].internalMappings` | `[]AzureVirtualNetworkGatewayNatRuleMapping` | yes |  |  |
| `spec.natRules[].internalMappings[].addressSpace` | `string` | yes |  |  |
| `spec.natRules[].internalMappings[].portRange` | `string` |  |  |  |
| `spec.natRules[].ipConfigurationId` | `string` |  |  |  |
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

- rule: Gateway names start with a letter or number, end with a letter, number, or underscore, and may contain alphanumerics, underscores, periods, and hyphens
- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.type

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_virtual_network_gateway_type_unspecified`
- `VPN`
- `EXPRESS_ROUTE`

### spec.vpnType

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_virtual_network_gateway_vpn_type_unspecified`
- `ROUTE_BASED`
- `POLICY_BASED`

### spec.sku

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_virtual_network_gateway_sku_unspecified`
- `BASIC`
- `STANDARD`
- `HIGH_PERFORMANCE`
- `ULTRA_PERFORMANCE`
- `VPN_GW_1`
- `VPN_GW_2`
- `VPN_GW_3`
- `VPN_GW_4`
- `VPN_GW_5`
- `VPN_GW_1_AZ`
- `VPN_GW_2_AZ`
- `VPN_GW_3_AZ`
- `VPN_GW_4_AZ`
- `VPN_GW_5_AZ`
- `ER_GW_1_AZ`
- `ER_GW_2_AZ`
- `ER_GW_3_AZ`
- `ER_GW_SCALE`

### spec.generation

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_virtual_network_gateway_generation_unspecified`
- `GENERATION1`
- `GENERATION2`
- `NONE`

### spec.ipConfigurations

`[]AzureVirtualNetworkGatewayIpConfiguration` · required

- rule: {"repeated":{"minItems":"1","maxItems":"3"}}

### spec.ipConfigurations[].name

`string`

### spec.ipConfigurations[].subnetId

`string | valueFrom` · required

- references: AzureSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.ipConfigurations[].publicIpAddressId

`string | valueFrom`

- references: AzurePublicIp (`status.outputs.public_ip_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzurePublicIp, name: <that resource's name>, fieldPath: status.outputs.public_ip_id}} -- a bare string does not parse

### spec.ipConfigurations[].privateIpAddressAllocation

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_virtual_network_gateway_ip_allocation_unspecified`
- `DYNAMIC`
- `STATIC`

### spec.activeActive

`bool`

### spec.privateIpAddressEnabled

`bool`

### spec.edgeZone

`string`

### spec.bgpEnabled

`bool`

### spec.bgpSettings

`AzureVirtualNetworkGatewayBgpSettings`

### spec.bgpSettings.asn

`int64`

- rule: {"int64":{"gte":"0"}}

### spec.bgpSettings.peerWeight

`int32`

- rule: {"int32":{"lte":100,"gte":0}}

### spec.bgpSettings.peeringAddresses

`[]AzureVirtualNetworkGatewayBgpPeeringAddress`

- rule: {"repeated":{"maxItems":"2"}}

### spec.bgpSettings.peeringAddresses[].ipConfigurationName

`string`

### spec.bgpSettings.peeringAddresses[].apipaAddresses

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"ipv4":true}}}}

### spec.customRouteAddressPrefixes

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.defaultLocalNetworkGatewayId

`string | valueFrom`

- references: AzureLocalNetworkGateway (`status.outputs.local_network_gateway_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureLocalNetworkGateway, name: <that resource's name>, fieldPath: status.outputs.local_network_gateway_id}} -- a bare string does not parse

### spec.vpnClientConfiguration

`AzureVirtualNetworkGatewayVpnClientConfiguration`

- rule: Entra ID authentication needs all three of aad_tenant, aad_audience, and aad_issuer
- rule: radius_server_address and radius_server_secret are set together

### spec.vpnClientConfiguration.addressSpaces

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.vpnClientConfiguration.aadTenant

`string`

### spec.vpnClientConfiguration.aadAudience

`string`

### spec.vpnClientConfiguration.aadIssuer

`string`

### spec.vpnClientConfiguration.rootCertificates

`[]AzureVirtualNetworkGatewayVpnClientRootCertificate`

### spec.vpnClientConfiguration.rootCertificates[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vpnClientConfiguration.rootCertificates[].publicCertData

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vpnClientConfiguration.revokedCertificates

`[]AzureVirtualNetworkGatewayVpnClientRevokedCertificate`

### spec.vpnClientConfiguration.revokedCertificates[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vpnClientConfiguration.revokedCertificates[].thumbprint

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vpnClientConfiguration.radiusServerAddress

`string`

### spec.vpnClientConfiguration.radiusServerSecret

`string | valueFrom` · sensitive

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.vpnClientConfiguration.radiusServers

`[]AzureVirtualNetworkGatewayVpnClientRadiusServer`

### spec.vpnClientConfiguration.radiusServers[].address

`string` · required

- rule: {"required":true,"string":{"ipv4":true}}

### spec.vpnClientConfiguration.radiusServers[].secret

`string | valueFrom` · required · sensitive

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.vpnClientConfiguration.radiusServers[].score

`int32`

- rule: {"int32":{"lte":30,"gte":1}}

### spec.vpnClientConfiguration.ipsecPolicy

`AzureVirtualNetworkGatewayVpnClientIpsecPolicy`

### spec.vpnClientConfiguration.ipsecPolicy.dhGroup

`string`

- rule: {"string":{"in":["DHGroup1","DHGroup14","DHGroup2","DHGroup2048","DHGroup24","ECP256","ECP384","None"]}}

### spec.vpnClientConfiguration.ipsecPolicy.ikeEncryption

`string`

- rule: {"string":{"in":["AES128","AES192","AES256","DES","DES3","GCMAES128","GCMAES256"]}}

### spec.vpnClientConfiguration.ipsecPolicy.ikeIntegrity

`string`

- rule: {"string":{"in":["GCMAES128","GCMAES256","MD5","SHA1","SHA256","SHA384"]}}

### spec.vpnClientConfiguration.ipsecPolicy.ipsecEncryption

`string`

- rule: {"string":{"in":["AES128","AES192","AES256","DES","DES3","GCMAES128","GCMAES192","GCMAES256","None"]}}

### spec.vpnClientConfiguration.ipsecPolicy.ipsecIntegrity

`string`

- rule: {"string":{"in":["GCMAES128","GCMAES192","GCMAES256","MD5","SHA1","SHA256"]}}

### spec.vpnClientConfiguration.ipsecPolicy.pfsGroup

`string`

- rule: {"string":{"in":["ECP256","ECP384","None","PFS1","PFS14","PFS2","PFS2048","PFS24","PFSMM"]}}

### spec.vpnClientConfiguration.ipsecPolicy.saLifetimeSeconds

`int32`

- rule: {"int32":{"lte":172799,"gte":300}}

### spec.vpnClientConfiguration.ipsecPolicy.saDataSizeKilobytes

`int32`

- rule: {"int32":{"gte":1024}}

### spec.vpnClientConfiguration.vpnClientProtocols

`[]string`

- rule: {"repeated":{"items":{"string":{"in":["IkeV2","OpenVPN","SSTP"]}}}}

### spec.vpnClientConfiguration.vpnAuthTypes

`[]string`

- rule: {"repeated":{"maxItems":"3","items":{"string":{"in":["Certificate","AAD","Radius"]}}}}

### spec.vpnClientConfiguration.clientConnections

`[]AzureVirtualNetworkGatewayClientConnection`

### spec.vpnClientConfiguration.clientConnections[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vpnClientConfiguration.clientConnections[].policyGroupNames

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.vpnClientConfiguration.clientConnections[].addressPrefixes

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.policyGroups

`[]AzureVirtualNetworkGatewayPolicyGroup`

### spec.policyGroups[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.policyGroups[].policyMembers

`[]AzureVirtualNetworkGatewayPolicyMember` · required

- rule: {"repeated":{"minItems":"1"}}

### spec.policyGroups[].policyMembers[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.policyGroups[].policyMembers[].type

`string`

- rule: {"string":{"in":["AADGroupId","CertificateGroupId","RadiusAzureGroupId"]}}

### spec.policyGroups[].policyMembers[].value

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.policyGroups[].isDefault

`bool`

### spec.policyGroups[].priority

`int32`

- rule: {"int32":{"gte":0}}

### spec.bgpRouteTranslationForNatEnabled

`bool`

### spec.dnsForwardingEnabled

`bool`

### spec.ipSecReplayProtectionEnabled

`bool` · optional (explicit presence)

- default: `true`

### spec.minimumScaleUnit

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":40,"gte":1}}

### spec.maximumScaleUnit

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":40,"gte":1}}

### spec.remoteVnetTrafficEnabled

`bool`

### spec.virtualWanTrafficEnabled

`bool`

### spec.natRules

`[]AzureVirtualNetworkGatewayNatRule`

### spec.natRules[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.natRules[].mode

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_virtual_network_gateway_nat_rule_mode_unspecified`
- `EGRESS_SNAT`
- `INGRESS_SNAT`

### spec.natRules[].type

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_virtual_network_gateway_nat_rule_type_unspecified`
- `STATIC_NAT`
- `DYNAMIC_NAT`

### spec.natRules[].externalMappings

`[]AzureVirtualNetworkGatewayNatRuleMapping` · required

- rule: {"repeated":{"minItems":"1"}}

### spec.natRules[].externalMappings[].addressSpace

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.natRules[].externalMappings[].portRange

`string`

### spec.natRules[].internalMappings

`[]AzureVirtualNetworkGatewayNatRuleMapping` · required

- rule: {"repeated":{"minItems":"1"}}

### spec.natRules[].internalMappings[].addressSpace

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.natRules[].internalMappings[].portRange

`string`

### spec.natRules[].ipConfigurationId

`string`

### spec.tags

`map<string, string>`

## Validation Rules

- `sku_required`: Choose the gateway SKU explicitly -- it sizes throughput, cost, and zone redundancy (VPN_GW_1 is the common production entry point for VPN gateways)
- `vpn_sku_vocabulary`: VPN gateways use BASIC, STANDARD, HIGH_PERFORMANCE, or the VPN_GW_1..5[_AZ] SKUs -- the ER_GW/ULTRA_PERFORMANCE SKUs are ExpressRoute-only
- `express_route_sku_vocabulary`: ExpressRoute gateways use STANDARD, HIGH_PERFORMANCE, ULTRA_PERFORMANCE, ER_GW_1_AZ..3_AZ, or ER_GW_SCALE
- `policy_based_requires_basic_sku`: Policy-based VPN gateways support only the BASIC SKU (legacy IKEv1 -- prefer route-based for anything new)
- `generation1_sku_vocabulary`: Generation1 VPN gateways use BASIC, STANDARD, HIGH_PERFORMANCE, VPN_GW_1..3, or VPN_GW_1_AZ..3_AZ
- `generation2_sku_vocabulary`: Generation2 VPN gateways start at VPN_GW_2: use VPN_GW_2..5 or VPN_GW_2_AZ..5_AZ
- `generation_is_vpn_only`: The generation knob applies to VPN gateways only -- leave it unset (or NONE) on ExpressRoute gateways
- `express_route_gateway_has_no_public_ips`: ExpressRoute gateways get Azure-managed addressing -- remove public_ip_address_id from every ip_configuration
- `vpn_gateway_requires_public_ips`: VPN gateways require a public IP on every ip_configuration (tunnels terminate on it)
- `active_active_requires_two_ip_configurations`: An active-active gateway is a two-instance pair -- give it (at least) two ip_configurations, each with its own public IP
- `active_active_is_vpn_only`: Active-active mode applies to VPN gateways only
- `vpn_client_configuration_is_vpn_only`: Point-to-site (vpn_client_configuration) runs on VPN gateways only
- `scale_units_require_each_other`: minimum_scale_unit and maximum_scale_unit are set together (both bound the ER_GW_SCALE autoscaler)
- `scale_units_only_on_ergwscale`: Autoscale bounds apply only to the ER_GW_SCALE SKU
- `ergwscale_requires_scale_units`: The ER_GW_SCALE SKU requires minimum_scale_unit and maximum_scale_unit (Azure's autoscale contract)
- `scale_unit_floor_not_above_ceiling`: minimum_scale_unit cannot exceed maximum_scale_unit
- `nat_rule_names_unique`: NAT rule names must be unique on the gateway

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureVirtualNetworkGateway, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.virtual_network_gateway_id` | `string` |  |
| `status.outputs.virtual_network_gateway_name` | `string` |  |
| `status.outputs.nat_rule_ids` | `map<string, string>` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.ipConfigurations[].subnetId` | AzureSubnet | `status.outputs.subnet_id` |
| `spec.ipConfigurations[].publicIpAddressId` | AzurePublicIp | `status.outputs.public_ip_id` |
| `spec.defaultLocalNetworkGatewayId` | AzureLocalNetworkGateway | `status.outputs.local_network_gateway_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureVirtualNetworkGatewayConnection | `spec.virtualNetworkGatewayId` | `status.outputs.virtual_network_gateway_id` |
| AzureVirtualNetworkGatewayConnection | `spec.peerVirtualNetworkGatewayId` | `status.outputs.virtual_network_gateway_id` |

## See Also

- [Overview](../README.md)
