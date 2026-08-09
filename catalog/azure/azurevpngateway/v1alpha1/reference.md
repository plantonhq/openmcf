# AzureVpnGateway

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureVpnGateway
metadata:
  name: test-vpn-gateway
spec:
  region: eastus
  resourceGroup:
    value: test-rg
  name: hub-vpn-gateway
  virtualHubId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/virtualHubs/hub-eastus
  # The unset optional fields apply ARM's defaults: routing preference
  # "Microsoft Network", scale unit 1.
  bgpRouteTranslationForNatEnabled: true
  bgpSettings:
    asn: 65515
    peerWeight: 10
    instance0BgpPeeringAddress:
      customIps:
        - "169.254.21.5"
    instance1BgpPeeringAddress:
      customIps:
        - "169.254.22.5"
  natRules:
    - name: branch-overlap
      mode: INGRESS_SNAT
      type: STATIC_NAT
      externalMappings:
        - addressSpace: "100.64.10.0/24"
      internalMappings:
        - addressSpace: "192.168.10.0/24"
  tags:
    purpose: hack-test
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.virtualHubId` | `string \| valueFrom` | yes |  | AzureVirtualHub (`status.outputs.virtual_hub_id`) |
| `spec.routingPreference` | `enum` |  | `MICROSOFT_NETWORK` |  |
| `spec.scaleUnit` | `int32` |  | `1` |  |
| `spec.bgpRouteTranslationForNatEnabled` | `bool` |  |  |  |
| `spec.bgpSettings` | `AzureVpnGatewayBgpSettings` |  |  |  |
| `spec.bgpSettings.asn` | `int64` |  |  |  |
| `spec.bgpSettings.peerWeight` | `int32` |  |  |  |
| `spec.bgpSettings.instance0BgpPeeringAddress` | `AzureVpnGatewayInstanceBgpPeeringAddress` |  |  |  |
| `spec.bgpSettings.instance0BgpPeeringAddress.customIps` | `[]string` | yes |  |  |
| `spec.bgpSettings.instance1BgpPeeringAddress` | `AzureVpnGatewayInstanceBgpPeeringAddress` |  |  |  |
| `spec.bgpSettings.instance1BgpPeeringAddress.customIps` | `[]string` | yes |  |  |
| `spec.natRules` | `[]AzureVpnGatewayNatRule` |  |  |  |
| `spec.natRules[].name` | `string` | yes |  |  |
| `spec.natRules[].mode` | `enum` |  |  |  |
| `spec.natRules[].type` | `enum` |  |  |  |
| `spec.natRules[].externalMappings` | `[]AzureVpnGatewayNatRuleMapping` | yes |  |  |
| `spec.natRules[].externalMappings[].addressSpace` | `string` | yes |  |  |
| `spec.natRules[].externalMappings[].portRange` | `string` |  |  |  |
| `spec.natRules[].internalMappings` | `[]AzureVpnGatewayNatRuleMapping` | yes |  |  |
| `spec.natRules[].internalMappings[].addressSpace` | `string` | yes |  |  |
| `spec.natRules[].internalMappings[].portRange` | `string` |  |  |  |
| `spec.natRules[].ipConfiguration` | `enum` |  |  |  |
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

### spec.virtualHubId

`string | valueFrom` · required

- references: AzureVirtualHub (`status.outputs.virtual_hub_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVirtualHub, name: <that resource's name>, fieldPath: status.outputs.virtual_hub_id}} -- a bare string does not parse

### spec.routingPreference

`enum` · optional (explicit presence)

- default: `MICROSOFT_NETWORK`
- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_vpn_gateway_routing_preference_unspecified`
- `MICROSOFT_NETWORK`
- `INTERNET`

### spec.scaleUnit

`int32` · optional (explicit presence)

- default: `1`
- rule: {"int32":{"gte":0}}

### spec.bgpRouteTranslationForNatEnabled

`bool`

### spec.bgpSettings

`AzureVpnGatewayBgpSettings`

### spec.bgpSettings.asn

`int64`

- rule: {"int64":{"gte":"1"}}

### spec.bgpSettings.peerWeight

`int32`

- rule: {"int32":{"lte":100,"gte":0}}

### spec.bgpSettings.instance0BgpPeeringAddress

`AzureVpnGatewayInstanceBgpPeeringAddress`

### spec.bgpSettings.instance0BgpPeeringAddress.customIps

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"ipv4":true}}}}

### spec.bgpSettings.instance1BgpPeeringAddress

`AzureVpnGatewayInstanceBgpPeeringAddress`

### spec.bgpSettings.instance1BgpPeeringAddress.customIps

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"ipv4":true}}}}

### spec.natRules

`[]AzureVpnGatewayNatRule`

### spec.natRules[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.natRules[].mode

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_vpn_gateway_nat_rule_mode_unspecified`
- `EGRESS_SNAT`
- `INGRESS_SNAT`

### spec.natRules[].type

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_vpn_gateway_nat_rule_type_unspecified`
- `STATIC_NAT`
- `DYNAMIC_NAT`

### spec.natRules[].externalMappings

`[]AzureVpnGatewayNatRuleMapping` · required

- rule: {"repeated":{"minItems":"1"}}

### spec.natRules[].externalMappings[].addressSpace

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.natRules[].externalMappings[].portRange

`string`

### spec.natRules[].internalMappings

`[]AzureVpnGatewayNatRuleMapping` · required

- rule: {"repeated":{"minItems":"1"}}

### spec.natRules[].internalMappings[].addressSpace

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.natRules[].internalMappings[].portRange

`string`

### spec.natRules[].ipConfiguration

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_vpn_gateway_nat_rule_ip_configuration_unspecified`
- `INSTANCE_0`
- `INSTANCE_1`

### spec.tags

`map<string, string>`

## Validation Rules

- `nat_rule_names_unique`: NAT rule names must be unique on the gateway -- each is the key the nat_rule_ids output uses

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureVpnGateway, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.vpn_gateway_id` | `string` |  |
| `status.outputs.vpn_gateway_name` | `string` |  |
| `status.outputs.bgp_asn` | `int64` |  |
| `status.outputs.public_ip_addresses` | `[]string` |  |
| `status.outputs.private_ip_addresses` | `[]string` |  |
| `status.outputs.nat_rule_ids` | `map<string, string>` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.virtualHubId` | AzureVirtualHub | `status.outputs.virtual_hub_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureVpnGatewayConnection | `spec.vpnGatewayId` | `status.outputs.vpn_gateway_id` |
| AzureVpnGatewayConnection | `spec.vpnLinks[].egressNatRuleIds` | `status.outputs.nat_rule_ids` |
| AzureVpnGatewayConnection | `spec.vpnLinks[].ingressNatRuleIds` | `status.outputs.nat_rule_ids` |

## See Also

- [Overview](../README.md)
