# AzureTrafficManagerProfile

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Offline-plan test manifest. Exercises the full surface: Performance
# routing, the complete monitor shape (HTTPS probe with path, custom
# status ranges, a Host header, zero-tolerance failures), an explicit
# low DNS TTL, Traffic View, and user tags merged over the derived
# ones. (max_return belongs to MultiValue routing -- its shape rides a
# dedicated offline plan per the profile's record.)
apiVersion: azure.planton.dev/v1alpha1
kind: AzureTrafficManagerProfile
metadata:
  name: test-traffic-manager-profile
  org: test-org
  env: dev
spec:
  resourceGroup:
    value: platform-rg
  name: app-director
  routingMethod: Performance
  dnsConfig:
    relativeName: contoso-app-director
    ttlSeconds: 30
  monitorConfig:
    protocol: HTTPS
    port: 443
    path: /healthz
    intervalInSeconds: 30
    timeoutInSeconds: 9
    toleratedNumberOfFailures: 0
    expectedStatusCodeRanges:
      - 200-299
      - 301-301
    customHeaders:
      - name: Host
        value: app.contoso.com
  trafficViewEnabled: true
  tags:
    cost-center: platform
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.routingMethod` | `string` | yes |  |  |
| `spec.dnsConfig` | `AzureTrafficManagerDnsConfig` | yes |  |  |
| `spec.dnsConfig.relativeName` | `string` | yes |  |  |
| `spec.dnsConfig.ttlSeconds` | `int32` |  | `60` |  |
| `spec.monitorConfig` | `AzureTrafficManagerMonitorConfig` | yes |  |  |
| `spec.monitorConfig.protocol` | `string` | yes |  |  |
| `spec.monitorConfig.port` | `int32` | yes |  |  |
| `spec.monitorConfig.path` | `string` |  |  |  |
| `spec.monitorConfig.intervalInSeconds` | `int32` |  | `30` |  |
| `spec.monitorConfig.timeoutInSeconds` | `int32` |  | `10` |  |
| `spec.monitorConfig.toleratedNumberOfFailures` | `int32` |  | `3` |  |
| `spec.monitorConfig.expectedStatusCodeRanges` | `[]string` |  |  |  |
| `spec.monitorConfig.customHeaders` | `[]AzureTrafficManagerCustomHeader` |  |  |  |
| `spec.monitorConfig.customHeaders[].name` | `string` | yes |  |  |
| `spec.monitorConfig.customHeaders[].value` | `string` | yes |  |  |
| `spec.enabled` | `bool` |  | `true` |  |
| `spec.maxReturn` | `int32` |  |  |  |
| `spec.trafficViewEnabled` | `bool` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.resourceGroup

`string | valueFrom` · required

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.routingMethod

`string` · required

- rule: {"required":true,"string":{"in":["Performance","Priority","Weighted","Geographic","Subnet","MultiValue"]}}

### spec.dnsConfig

`AzureTrafficManagerDnsConfig` · required

- rule: {"required":true}

### spec.dnsConfig.relativeName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.dnsConfig.ttlSeconds

`int32` · optional (explicit presence)

- default: `60`
- rule: {"int32":{"lte":2147483647,"gte":0}}

### spec.monitorConfig

`AzureTrafficManagerMonitorConfig` · required

- rule: {"required":true}
- rule: With interval_in_seconds 10 (fast interval), set timeout_in_seconds explicitly to 5-9 -- the default 10 does not fit inside the probe window

### spec.monitorConfig.protocol

`string` · required

- rule: {"required":true,"string":{"in":["HTTP","HTTPS","TCP"]}}

### spec.monitorConfig.port

`int32` · required · optional (explicit presence)

- rule: {"required":true,"int32":{"lte":65535,"gte":1}}

### spec.monitorConfig.path

`string`

### spec.monitorConfig.intervalInSeconds

`int32` · optional (explicit presence)

- default: `30`
- rule: {"int32":{"in":[10,30]}}

### spec.monitorConfig.timeoutInSeconds

`int32` · optional (explicit presence)

- default: `10`
- rule: {"int32":{"lte":10,"gte":5}}

### spec.monitorConfig.toleratedNumberOfFailures

`int32` · optional (explicit presence)

- default: `3`
- rule: {"int32":{"lte":9,"gte":0}}

### spec.monitorConfig.expectedStatusCodeRanges

`[]string`

- rule: {"repeated":{"items":{"string":{"pattern":"^[0-9]{3}-[0-9]{3}$"}}}}

### spec.monitorConfig.customHeaders

`[]AzureTrafficManagerCustomHeader`

### spec.monitorConfig.customHeaders[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.monitorConfig.customHeaders[].value

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.enabled

`bool` · optional (explicit presence)

- default: `true`

### spec.maxReturn

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":8,"gte":1}}

### spec.trafficViewEnabled

`bool`

### spec.tags

`map<string, string>`

## Validation Rules

- `azure_traffic_manager_multivalue_requires_max_return`: MultiValue routing requires max_return (1-8) -- the number of healthy addresses returned per answer

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureTrafficManagerProfile, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.traffic_manager_profile_id` | `string` |  |
| `status.outputs.traffic_manager_profile_name` | `string` |  |
| `status.outputs.fqdn` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureTrafficManagerEndpoint | `spec.profileId` | `status.outputs.traffic_manager_profile_id` |
| AzureTrafficManagerEndpoint | `spec.nested.targetProfileId` | `status.outputs.traffic_manager_profile_id` |

## See Also

- [Overview](../README.md)
