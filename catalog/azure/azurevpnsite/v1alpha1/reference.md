# AzureVpnSite

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureVpnSite
metadata:
  name: test-vpn-site
spec:
  region: eastus
  resourceGroup:
    value: test-rg
  name: branch-london
  virtualWanId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/virtualWans/global-wan
  addressCidrs:
    - "192.168.10.0/24"
  deviceVendor: Cisco
  deviceModel: ISR4331
  links:
    - name: primary-isp
      providerName: ExampleCarrier
      speedInMbps: 200
      ipAddress: "203.0.113.10"
    - name: backup-isp
      speedInMbps: 100
      ipAddress: "198.51.100.7"
      bgp:
        asn: 65010
        peeringAddress: "169.254.21.2"
  o365Policy:
    trafficCategory:
      optimizeEndpointEnabled: true
  tags:
    purpose: hack-test
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.virtualWanId` | `string \| valueFrom` | yes |  | AzureVirtualWan (`status.outputs.virtual_wan_id`) |
| `spec.addressCidrs` | `[]string` |  |  |  |
| `spec.deviceVendor` | `string` |  |  |  |
| `spec.deviceModel` | `string` |  |  |  |
| `spec.links` | `[]AzureVpnSiteLink` |  |  |  |
| `spec.links[].name` | `string` | yes |  |  |
| `spec.links[].providerName` | `string` |  |  |  |
| `spec.links[].speedInMbps` | `int32` |  |  |  |
| `spec.links[].ipAddress` | `string` |  |  |  |
| `spec.links[].fqdn` | `string` |  |  |  |
| `spec.links[].bgp` | `AzureVpnSiteLinkBgp` |  |  |  |
| `spec.links[].bgp.asn` | `int64` |  |  |  |
| `spec.links[].bgp.peeringAddress` | `string` | yes |  |  |
| `spec.o365Policy` | `AzureVpnSiteO365Policy` |  |  |  |
| `spec.o365Policy.trafficCategory` | `AzureVpnSiteO365TrafficCategory` |  |  |  |
| `spec.o365Policy.trafficCategory.allowEndpointEnabled` | `bool` |  |  |  |
| `spec.o365Policy.trafficCategory.defaultEndpointEnabled` | `bool` |  |  |  |
| `spec.o365Policy.trafficCategory.optimizeEndpointEnabled` | `bool` |  |  |  |
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

- rule: {"required":true,"string":{"pattern":"^[^'<>%&:?/+]+$"}}

### spec.virtualWanId

`string | valueFrom` · required

- references: AzureVirtualWan (`status.outputs.virtual_wan_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVirtualWan, name: <that resource's name>, fieldPath: status.outputs.virtual_wan_id}} -- a bare string does not parse

### spec.addressCidrs

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.deviceVendor

`string`

### spec.deviceModel

`string`

### spec.links

`[]AzureVpnSiteLink`

- rule: Give the link a public endpoint: ip_address (static public IP) or fqdn (re-resolved name)

### spec.links[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.links[].providerName

`string`

### spec.links[].speedInMbps

`int32`

- rule: {"int32":{"gte":0}}

### spec.links[].ipAddress

`string`

- rule: ip_address must be an IP address

### spec.links[].fqdn

`string`

### spec.links[].bgp

`AzureVpnSiteLinkBgp`

### spec.links[].bgp.asn

`int64`

- rule: {"int64":{"gte":"1"}}

### spec.links[].bgp.peeringAddress

`string` · required

- rule: peering_address must be an IP address
- rule: {"required":true}

### spec.o365Policy

`AzureVpnSiteO365Policy`

### spec.o365Policy.trafficCategory

`AzureVpnSiteO365TrafficCategory`

### spec.o365Policy.trafficCategory.allowEndpointEnabled

`bool`

### spec.o365Policy.trafficCategory.defaultEndpointEnabled

`bool`

### spec.o365Policy.trafficCategory.optimizeEndpointEnabled

`bool`

### spec.tags

`map<string, string>`

## Validation Rules

- `link_names_unique`: Link names must be unique on the site -- each is the key connections and the link_ids output use

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureVpnSite, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.vpn_site_id` | `string` |  |
| `status.outputs.vpn_site_name` | `string` |  |
| `status.outputs.link_ids` | `map<string, string>` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.virtualWanId` | AzureVirtualWan | `status.outputs.virtual_wan_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureVpnGatewayConnection | `spec.remoteVpnSiteId` | `status.outputs.vpn_site_id` |
| AzureVpnGatewayConnection | `spec.vpnLinks[].vpnSiteLinkId` | `status.outputs.link_ids` |

## See Also

- [Overview](../README.md)
