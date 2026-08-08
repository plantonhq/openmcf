# AzureLocalNetworkGateway

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureLocalNetworkGateway
metadata:
  name: test-local-network-gateway
spec:
  region: eastus
  resourceGroup:
    value: test-rg
  name: hq-datacenter
  # The on-premises VPN device's public endpoint (exactly one of
  # gatewayAddress or gatewayFqdn).
  gatewayAddress: "203.0.113.10"
  # The prefixes reachable behind the device -- Azure routes these into
  # the tunnel. Leave empty only when bgpSettings carries the routing.
  addressSpaces:
    - "192.168.100.0/24"
    - "192.168.101.0/24"
  bgpSettings:
    asn: 65010
    # An IP INSIDE the tunnel (the device's tunnel interface), not the
    # device's public address.
    bgpPeeringAddress: "10.255.255.1"
    peerWeight: 0
  tags:
    purpose: hack-test
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.gatewayAddress` | `string` |  |  |  |
| `spec.gatewayFqdn` | `string` |  |  |  |
| `spec.addressSpaces` | `[]string` |  |  |  |
| `spec.bgpSettings` | `AzureLocalNetworkGatewayBgpSettings` |  |  |  |
| `spec.bgpSettings.asn` | `int64` |  |  |  |
| `spec.bgpSettings.bgpPeeringAddress` | `string` | yes |  |  |
| `spec.bgpSettings.peerWeight` | `int32` |  |  |  |
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

### spec.gatewayAddress

`string`

- rule: gateway_address must be an IPv4 address

### spec.gatewayFqdn

`string`

### spec.addressSpaces

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.bgpSettings

`AzureLocalNetworkGatewayBgpSettings`

### spec.bgpSettings.asn

`int64`

- rule: {"int64":{"gte":"1"}}

### spec.bgpSettings.bgpPeeringAddress

`string` · required

- rule: {"required":true,"string":{"ipv4":true}}

### spec.bgpSettings.peerWeight

`int32`

- rule: {"int32":{"lte":100,"gte":0}}

### spec.tags

`map<string, string>`

## Validation Rules

- `exactly_one_gateway_endpoint`: Describe the on-premises endpoint with exactly one of gateway_address (static public IP) or gateway_fqdn (re-resolved name)
- `routing_source_required`: Azure needs a routing source for the site: static address_spaces, bgp_settings, or both -- an empty object routes nothing

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureLocalNetworkGateway, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.local_network_gateway_id` | `string` |  |
| `status.outputs.local_network_gateway_name` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureVirtualNetworkGateway | `spec.defaultLocalNetworkGatewayId` | `status.outputs.local_network_gateway_id` |
| AzureVirtualNetworkGatewayConnection | `spec.localNetworkGatewayId` | `status.outputs.local_network_gateway_id` |

## See Also

- [Overview](../README.md)
