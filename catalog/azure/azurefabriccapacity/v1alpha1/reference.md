# AzureFabricCapacity

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Deep-shape example for docs and offline validation: an F2 capacity
# with two administrators and tags. References are literal values so
# the manifest validates standalone.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureFabricCapacity
metadata:
  name: test-fabric-capacity
  id: test-fabric-capacity
  org: test-org
  env: test
spec:
  resourceGroup:
    value: test-rg
  name: testorgfabric
  region: eastus
  skuName: F2
  administrationMembers:
    - admin@testorg.example
    - 11111111-2222-3333-4444-555555555555
  tags:
    costCenter: analytics
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.region` | `string` | yes |  |  |
| `spec.skuName` | `string` | yes |  |  |
| `spec.administrationMembers` | `[]string` | yes |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.resourceGroup

`string | valueFrom` · required

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

- rule: Capacity names must be 3-63 characters of lowercase letters and numbers, starting with a letter
- rule: {"required":true}

### spec.region

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.skuName

`string` · required

- rule: {"required":true,"string":{"in":["F2","F4","F8","F16","F32","F64","F128","F256","F512","F1024","F2048"]}}

### spec.administrationMembers

`[]string` · required

- rule: {"repeated":{"minItems":"1","unique":true,"items":{"string":{"minLen":"1"}}}}

### spec.tags

`map<string, string>`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureFabricCapacity, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.fabric_capacity_id` | `string` |  |
| `status.outputs.fabric_capacity_name` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |

## See Also

- [Overview](../README.md)
