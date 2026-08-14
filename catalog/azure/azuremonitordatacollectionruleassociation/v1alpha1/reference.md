# AzureMonitorDataCollectionRuleAssociation

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Deep-shape example for docs and offline validation: one VM attached
# to a data collection rule. References are literal ARM ids so the
# manifest validates standalone.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureMonitorDataCollectionRuleAssociation
metadata:
  name: test-dcra
  id: test-dcra
  org: test-org
  env: test
spec:
  targetResourceId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/app-rg/providers/Microsoft.Compute/virtualMachines/app-vm
  # Name the association after the rule it binds -- association names
  # are what fleet listings show.
  name: linux-baseline-assoc
  dataCollectionRuleId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Insights/dataCollectionRules/linux-baseline
  description: attaches the app VM to the Linux baseline rule
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.targetResourceId` | `string \| valueFrom` | yes |  |  |
| `spec.name` | `string` |  |  |  |
| `spec.dataCollectionRuleId` | `string \| valueFrom` |  |  | AzureMonitorDataCollectionRule (`status.outputs.data_collection_rule_id`) |
| `spec.dataCollectionEndpointId` | `string \| valueFrom` |  |  |  |
| `spec.description` | `string` |  |  |  |

## Field Details

### spec.targetResourceId

`string | valueFrom` · required

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.name

`string`

### spec.dataCollectionRuleId

`string | valueFrom`

- references: AzureMonitorDataCollectionRule (`status.outputs.data_collection_rule_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureMonitorDataCollectionRule, name: <that resource's name>, fieldPath: status.outputs.data_collection_rule_id}} -- a bare string does not parse

### spec.dataCollectionEndpointId

`string | valueFrom`

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.description

`string`

## Validation Rules

- `dcra_exactly_one_binding`: set exactly one of data_collection_rule_id and data_collection_endpoint_id
- `dcra_rule_binding_requires_name`: name is required when data_collection_rule_id is set

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureMonitorDataCollectionRuleAssociation, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.data_collection_rule_association_id` | `string` |  |
| `status.outputs.data_collection_rule_association_name` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.dataCollectionRuleId` | AzureMonitorDataCollectionRule | `status.outputs.data_collection_rule_id` |

## See Also

- [Overview](../README.md)
