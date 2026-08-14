# AzureTrafficManagerEndpoint

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Offline-plan test manifest. Exercises the external variant at depth:
# an explicit endpoint location (the Performance-routing requirement),
# always-serve, an explicit weight and priority, and a per-endpoint
# probe Host header. (The azure and nested variants ride dedicated
# offline plans per the profile's record; subnet and geo claims belong
# to Subnet/Geographic-routed profiles and ride the subnet-claims
# offline plan.)
apiVersion: azure.planton.dev/v1alpha1
kind: AzureTrafficManagerEndpoint
metadata:
  name: test-traffic-manager-endpoint
  org: test-org
  env: dev
spec:
  profileId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/platform-rg/providers/Microsoft.Network/trafficManagerProfiles/app-director
  name: eastus-app
  external:
    target:
      value: app-eastus.contoso.com
    endpointLocation: eastus
    alwaysServeEnabled: true
  weight: 100
  priority: 1
  customHeaders:
    - name: Host
      value: app-eastus.contoso.com
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.profileId` | `string \| valueFrom` | yes |  | AzureTrafficManagerProfile (`status.outputs.traffic_manager_profile_id`) |
| `spec.name` | `string` | yes |  |  |
| `spec.azure` | `AzureTrafficManagerAzureEndpoint` |  |  |  |
| `spec.azure.targetResourceId` | `string \| valueFrom` | yes |  |  |
| `spec.azure.alwaysServeEnabled` | `bool` |  | `false` |  |
| `spec.external` | `AzureTrafficManagerExternalEndpoint` |  |  |  |
| `spec.external.target` | `string \| valueFrom` | yes |  |  |
| `spec.external.endpointLocation` | `string` |  |  |  |
| `spec.external.alwaysServeEnabled` | `bool` |  | `false` |  |
| `spec.nested` | `AzureTrafficManagerNestedEndpoint` |  |  |  |
| `spec.nested.targetProfileId` | `string \| valueFrom` | yes |  | AzureTrafficManagerProfile (`status.outputs.traffic_manager_profile_id`) |
| `spec.nested.minimumChildEndpoints` | `int32` | yes |  |  |
| `spec.nested.minimumRequiredChildEndpointsIpv4` | `int32` |  |  |  |
| `spec.nested.minimumRequiredChildEndpointsIpv6` | `int32` |  |  |  |
| `spec.nested.endpointLocation` | `string` |  |  |  |
| `spec.weight` | `int32` |  | `1` |  |
| `spec.priority` | `int32` |  |  |  |
| `spec.enabled` | `bool` |  | `true` |  |
| `spec.geoMappings` | `[]string` |  |  |  |
| `spec.subnets` | `[]AzureTrafficManagerEndpointSubnet` |  |  |  |
| `spec.subnets[].first` | `string` | yes |  |  |
| `spec.subnets[].last` | `string` |  |  |  |
| `spec.subnets[].scope` | `int32` |  |  |  |
| `spec.customHeaders` | `[]AzureTrafficManagerEndpointCustomHeader` |  |  |  |
| `spec.customHeaders[].name` | `string` | yes |  |  |
| `spec.customHeaders[].value` | `string` | yes |  |  |

## Field Details

### spec.profileId

`string | valueFrom` · required

- references: AzureTrafficManagerProfile (`status.outputs.traffic_manager_profile_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureTrafficManagerProfile, name: <that resource's name>, fieldPath: status.outputs.traffic_manager_profile_id}} -- a bare string does not parse

### spec.name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.azure

`AzureTrafficManagerAzureEndpoint`

### spec.azure.targetResourceId

`string | valueFrom` · required

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.azure.alwaysServeEnabled

`bool` · optional (explicit presence)

- default: `false`

### spec.external

`AzureTrafficManagerExternalEndpoint`

### spec.external.target

`string | valueFrom` · required

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.external.endpointLocation

`string`

### spec.external.alwaysServeEnabled

`bool` · optional (explicit presence)

- default: `false`

### spec.nested

`AzureTrafficManagerNestedEndpoint`

### spec.nested.targetProfileId

`string | valueFrom` · required

- references: AzureTrafficManagerProfile (`status.outputs.traffic_manager_profile_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureTrafficManagerProfile, name: <that resource's name>, fieldPath: status.outputs.traffic_manager_profile_id}} -- a bare string does not parse

### spec.nested.minimumChildEndpoints

`int32` · required · optional (explicit presence)

- rule: {"required":true,"int32":{"gte":1}}

### spec.nested.minimumRequiredChildEndpointsIpv4

`int32` · optional (explicit presence)

- rule: {"int32":{"gte":0}}

### spec.nested.minimumRequiredChildEndpointsIpv6

`int32` · optional (explicit presence)

- rule: {"int32":{"gte":0}}

### spec.nested.endpointLocation

`string`

### spec.weight

`int32` · optional (explicit presence)

- default: `1`
- rule: {"int32":{"lte":1000,"gte":1}}

### spec.priority

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":1000,"gte":1}}

### spec.enabled

`bool` · optional (explicit presence)

- default: `true`

### spec.geoMappings

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.subnets

`[]AzureTrafficManagerEndpointSubnet`

### spec.subnets[].first

`string` · required

- rule: {"required":true,"string":{"ipv4":true}}

### spec.subnets[].last

`string`

- rule: last must be an IPv4 address

### spec.subnets[].scope

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":32,"gte":0}}

### spec.customHeaders

`[]AzureTrafficManagerEndpointCustomHeader`

### spec.customHeaders[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.customHeaders[].value

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

## Validation Rules

- `azure_traffic_manager_endpoint_exactly_one_variant`: Set exactly one endpoint variant -- azure, external, or nested -- the variant determines the endpoint type

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureTrafficManagerEndpoint, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.endpoint_id` | `string` |  |
| `status.outputs.endpoint_name` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.profileId` | AzureTrafficManagerProfile | `status.outputs.traffic_manager_profile_id` |
| `spec.nested.targetProfileId` | AzureTrafficManagerProfile | `status.outputs.traffic_manager_profile_id` |

## See Also

- [Overview](../README.md)
