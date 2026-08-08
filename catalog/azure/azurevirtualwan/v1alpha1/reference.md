# AzureVirtualWan

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureVirtualWan
metadata:
  name: test-virtual-wan
spec:
  region: eastus
  resourceGroup:
    value: test-rg
  name: global-wan
  # The unset optional fields apply ARM's defaults: type "Standard"
  # (the full-mesh tier), branch-to-branch transit on, VPN encryption
  # on, no Office 365 local breakout.
  tags:
    purpose: hack-test
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.disableVpnEncryption` | `bool` |  |  |  |
| `spec.allowBranchToBranchTraffic` | `bool` |  | `true` |  |
| `spec.office365LocalBreakoutCategory` | `enum` |  | `NONE` |  |
| `spec.type` | `string` |  | `Standard` |  |
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

### spec.disableVpnEncryption

`bool`

### spec.allowBranchToBranchTraffic

`bool` · optional (explicit presence)

- default: `true`

### spec.office365LocalBreakoutCategory

`enum` · optional (explicit presence)

- default: `NONE`
- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_virtual_wan_office365_breakout_category_unspecified`
- `NONE`
- `ALL`
- `OPTIMIZE`
- `OPTIMIZE_AND_ALLOW`

### spec.type

`string` · optional (explicit presence)

- default: `Standard`

### spec.tags

`map<string, string>`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureVirtualWan, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.virtual_wan_id` | `string` |  |
| `status.outputs.virtual_wan_name` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |

## See Also

- [Overview](../README.md)
