# AzureAvailabilitySet

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Deep-shape example for docs and offline validation: an availability
# set with explicit domain counts, managed alignment, and a proximity
# placement group. References are literal values so the manifest
# validates standalone.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureAvailabilitySet
metadata:
  name: test-availability-set
  id: test-availability-set
  org: test-org
  env: test
spec:
  resourceGroup:
    value: test-rg
  name: web-avset
  region: eastus
  platformUpdateDomainCount: 5
  platformFaultDomainCount: 3
  managed: true
  proximityPlacementGroupId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Compute/proximityPlacementGroups/web-ppg
  tags:
    tier: web
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.region` | `string` | yes |  |  |
| `spec.platformUpdateDomainCount` | `int32` |  |  |  |
| `spec.platformFaultDomainCount` | `int32` |  |  |  |
| `spec.managed` | `bool` |  |  |  |
| `spec.proximityPlacementGroupId` | `string \| valueFrom` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.resourceGroup

`string | valueFrom` · required

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

- rule: Availability set names are up to 80 letters, numbers, dots, dashes, and underscores, starting with a letter or number and ending with a letter, number, or underscore
- rule: {"required":true}

### spec.region

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.platformUpdateDomainCount

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":20,"gte":1}}

### spec.platformFaultDomainCount

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":3,"gte":1}}

### spec.managed

`bool` · optional (explicit presence)

### spec.proximityPlacementGroupId

`string | valueFrom`

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.tags

`map<string, string>`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureAvailabilitySet, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.availability_set_id` | `string` |  |
| `status.outputs.availability_set_name` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureVirtualMachine | `spec.availability.availabilitySetId` | `status.outputs.availability_set_id` |

## See Also

- [Overview](../README.md)
