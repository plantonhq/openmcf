# AzureEventgridDomain

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Deep-shape example for docs and offline validation: a pinned-topics
# CloudEvents domain with an IP allowlist. References are literal
# values so the manifest validates standalone.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureEventgridDomain
metadata:
  name: test-eventgrid-domain
  id: test-eventgrid-domain
  org: test-org
  env: test
spec:
  resourceGroup:
    value: test-rg
  # The name becomes a PUBLIC DNS hostname, unique across the region.
  name: test-org-tenant-events
  region: eastus
  inputSchema: CloudEventSchemaV1_0
  # The pinned-topics governance posture: topics exist only as declared
  # AzureEventgridDomainTopic resources.
  autoCreateTopicWithFirstSubscription: false
  autoDeleteTopicWithLastSubscription: false
  publicNetworkAccessEnabled: true
  localAuthEnabled: true
  inboundIpRules:
    - 203.0.113.0/24
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
| `spec.inputSchema` | `string` |  | `EventGridSchema` |  |
| `spec.inputMappingFields` | `AzureEventgridDomainInputMappingFields` |  |  |  |
| `spec.inputMappingFields.id` | `string` |  |  |  |
| `spec.inputMappingFields.topic` | `string` |  |  |  |
| `spec.inputMappingFields.eventTime` | `string` |  |  |  |
| `spec.inputMappingFields.eventType` | `string` |  |  |  |
| `spec.inputMappingFields.subject` | `string` |  |  |  |
| `spec.inputMappingFields.dataVersion` | `string` |  |  |  |
| `spec.inputMappingDefaultValues` | `AzureEventgridDomainInputMappingDefaultValues` |  |  |  |
| `spec.inputMappingDefaultValues.eventType` | `string` |  |  |  |
| `spec.inputMappingDefaultValues.subject` | `string` |  |  |  |
| `spec.inputMappingDefaultValues.dataVersion` | `string` |  |  |  |
| `spec.autoCreateTopicWithFirstSubscription` | `bool` |  | `true` |  |
| `spec.autoDeleteTopicWithLastSubscription` | `bool` |  | `true` |  |
| `spec.publicNetworkAccessEnabled` | `bool` |  | `true` |  |
| `spec.localAuthEnabled` | `bool` |  | `true` |  |
| `spec.inboundIpRules` | `[]string` |  |  |  |
| `spec.identity` | `AzureEventgridDomainIdentity` |  |  |  |
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

- rule: Domain names must be 3-50 characters of letters, numbers, and hyphens
- rule: {"required":true}

### spec.region

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.inputSchema

`string` · optional (explicit presence)

- default: `EventGridSchema`
- rule: {"string":{"in":["EventGridSchema","CloudEventSchemaV1_0","CustomEventSchema"]}}

### spec.inputMappingFields

`AzureEventgridDomainInputMappingFields`

### spec.inputMappingFields.id

`string`

### spec.inputMappingFields.topic

`string`

### spec.inputMappingFields.eventTime

`string`

### spec.inputMappingFields.eventType

`string`

### spec.inputMappingFields.subject

`string`

### spec.inputMappingFields.dataVersion

`string`

### spec.inputMappingDefaultValues

`AzureEventgridDomainInputMappingDefaultValues`

### spec.inputMappingDefaultValues.eventType

`string`

### spec.inputMappingDefaultValues.subject

`string`

### spec.inputMappingDefaultValues.dataVersion

`string`

### spec.autoCreateTopicWithFirstSubscription

`bool` · optional (explicit presence)

- default: `true`

### spec.autoDeleteTopicWithLastSubscription

`bool` · optional (explicit presence)

- default: `true`

### spec.publicNetworkAccessEnabled

`bool` · optional (explicit presence)

- default: `true`

### spec.localAuthEnabled

`bool` · optional (explicit presence)

- default: `true`

### spec.inboundIpRules

`[]string`

- rule: {"repeated":{"maxItems":"128","items":{"string":{"minLen":"1"}}}}

### spec.identity

`AzureEventgridDomainIdentity`

- rule: identity_ids is required for USER_ASSIGNED and must be empty for SYSTEM_ASSIGNED

### spec.identity.type

`enum` · required

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_eventgrid_domain_identity_type_unspecified`
- `SYSTEM_ASSIGNED`
- `USER_ASSIGNED`

### spec.identity.identityIds

`[]string | valueFrom`

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.tags

`map<string, string>`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureEventgridDomain, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.domain_id` | `string` |  |
| `status.outputs.domain_name` | `string` |  |
| `status.outputs.endpoint` | `string` |  |
| `status.outputs.primary_access_key` | `string` |  |
| `status.outputs.secondary_access_key` | `string` |  |
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
| AzureEventgridDomainTopic | `spec.domainId` | `status.outputs.domain_id` |

## See Also

- [Overview](../README.md)
