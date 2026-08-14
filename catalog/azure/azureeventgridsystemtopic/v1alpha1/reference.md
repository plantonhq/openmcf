# AzureEventgridSystemTopic

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Deep-shape example for docs and offline validation: a storage
# account's blob-event stream with a system-assigned identity for
# secured delivery. References are literal values so the manifest
# validates standalone.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureEventgridSystemTopic
metadata:
  name: test-eventgrid-system-topic
  id: test-eventgrid-system-topic
  org: test-org
  env: test
spec:
  resourceGroup:
    value: test-rg
  name: appdata-storage-events
  # Must match the SOURCE's region ("Global" for subscription/
  # resource-group sources).
  region: eastus
  sourceResourceId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Storage/storageAccounts/appdata
  topicType: Microsoft.Storage.StorageAccounts
  identity:
    type: SYSTEM_ASSIGNED
  tags:
    costCenter: platform
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.region` | `string` | yes |  |  |
| `spec.sourceResourceId` | `string \| valueFrom` | yes |  |  |
| `spec.topicType` | `string` | yes |  |  |
| `spec.identity` | `AzureEventgridSystemTopicIdentity` |  |  |  |
| `spec.identity.type` | `enum` | yes |  |  |
| `spec.identity.identityIds` | `[]string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.resourceGroup

`string | valueFrom` · required

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

- rule: System topic names must be 3-128 characters of letters, numbers, and hyphens
- rule: {"required":true}

### spec.region

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.sourceResourceId

`string | valueFrom` · required

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.topicType

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.identity

`AzureEventgridSystemTopicIdentity`

- rule: identity_ids is required for USER_ASSIGNED and SYSTEM_AND_USER_ASSIGNED and must be empty for SYSTEM_ASSIGNED

### spec.identity.type

`enum` · required

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_eventgrid_system_topic_identity_type_unspecified`
- `SYSTEM_ASSIGNED`
- `USER_ASSIGNED`
- `SYSTEM_AND_USER_ASSIGNED`

### spec.identity.identityIds

`[]string | valueFrom`

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.tags

`map<string, string>`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureEventgridSystemTopic, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.system_topic_id` | `string` |  |
| `status.outputs.system_topic_name` | `string` |  |
| `status.outputs.metric_resource_id` | `string` |  |
| `status.outputs.identity_principal_id` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.identity.identityIds` | AzureUserAssignedIdentity | `status.outputs.identity_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureEventgridEventSubscription | `spec.systemTopicId` | `status.outputs.system_topic_id` |

## See Also

- [Overview](../README.md)
