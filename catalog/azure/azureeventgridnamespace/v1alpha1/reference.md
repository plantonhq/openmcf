# AzureEventgridNamespace

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Deep-shape example for docs and offline validation: an MQTT-enabled
# namespace with capacity, an inbound IP rule, a combined identity,
# and the routing enrichments -- the route topic is included here so
# the offline plan renders the full MQTT block. References are literal
# values so the manifest validates standalone.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureEventgridNamespace
metadata:
  name: test-eventgrid-namespace
  id: test-eventgrid-namespace
  org: test-org
  env: test
spec:
  resourceGroup:
    value: test-rg
  name: acme-events-hub
  region: eastus
  capacity: 2
  publicNetworkAccessEnabled: true
  inboundIpRules:
    - 203.0.113.0/24
  identity:
    type: SYSTEM_AND_USER_ASSIGNED
    identityIds:
      - value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/test-uai
  topicSpacesConfiguration:
    alternativeAuthenticationNameSources:
      - ClientCertificateSubject
    maximumClientSessionsPerAuthenticationName: 3
    maximumSessionExpiryInHours: 4
    routeTopicId:
      value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.EventGrid/topics/test-route-topic
    dynamicRoutingEnrichments:
      - key: clientname
        value: ${client.authenticationName}
    staticRoutingEnrichments:
      - key: source
        value: mqtt-broker
  tags:
    costCenter: platform
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.region` | `string` | yes |  |  |
| `spec.capacity` | `int32` |  | `1` |  |
| `spec.publicNetworkAccessEnabled` | `bool` |  | `true` |  |
| `spec.inboundIpRules` | `[]string` |  |  |  |
| `spec.identity` | `AzureEventgridNamespaceIdentity` |  |  |  |
| `spec.identity.type` | `enum` | yes |  |  |
| `spec.identity.identityIds` | `[]string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.topicSpacesConfiguration` | `AzureEventgridNamespaceTopicSpacesConfiguration` |  |  |  |
| `spec.topicSpacesConfiguration.alternativeAuthenticationNameSources` | `[]string` |  |  |  |
| `spec.topicSpacesConfiguration.maximumClientSessionsPerAuthenticationName` | `int32` |  | `1` |  |
| `spec.topicSpacesConfiguration.maximumSessionExpiryInHours` | `int32` |  | `1` |  |
| `spec.topicSpacesConfiguration.routeTopicId` | `string \| valueFrom` |  |  | AzureEventgridTopic (`status.outputs.topic_id`) |
| `spec.topicSpacesConfiguration.dynamicRoutingEnrichments` | `[]AzureEventgridNamespaceRoutingEnrichment` |  |  |  |
| `spec.topicSpacesConfiguration.dynamicRoutingEnrichments[].key` | `string` | yes |  |  |
| `spec.topicSpacesConfiguration.dynamicRoutingEnrichments[].value` | `string` | yes |  |  |
| `spec.topicSpacesConfiguration.staticRoutingEnrichments` | `[]AzureEventgridNamespaceRoutingEnrichment` |  |  |  |
| `spec.topicSpacesConfiguration.staticRoutingEnrichments[].key` | `string` | yes |  |  |
| `spec.topicSpacesConfiguration.staticRoutingEnrichments[].value` | `string` | yes |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.resourceGroup

`string | valueFrom` · required

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

- rule: Namespace names must be 3-50 characters
- rule: {"required":true}

### spec.region

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.capacity

`int32` · optional (explicit presence)

- default: `1`
- rule: {"int32":{"lte":40,"gte":1}}

### spec.publicNetworkAccessEnabled

`bool` · optional (explicit presence)

- default: `true`

### spec.inboundIpRules

`[]string`

- rule: {"repeated":{"maxItems":"128","items":{"string":{"minLen":"1"}}}}

### spec.identity

`AzureEventgridNamespaceIdentity`

- rule: identity_ids is required for USER_ASSIGNED and SYSTEM_AND_USER_ASSIGNED and must be empty for SYSTEM_ASSIGNED

### spec.identity.type

`enum` · required

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_eventgrid_namespace_identity_type_unspecified`
- `SYSTEM_ASSIGNED`
- `USER_ASSIGNED`
- `SYSTEM_AND_USER_ASSIGNED`

### spec.identity.identityIds

`[]string | valueFrom`

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.topicSpacesConfiguration

`AzureEventgridNamespaceTopicSpacesConfiguration`

### spec.topicSpacesConfiguration.alternativeAuthenticationNameSources

`[]string`

- rule: {"repeated":{"items":{"string":{"in":["ClientCertificateSubject","ClientCertificateDns","ClientCertificateUri","ClientCertificateIp","ClientCertificateEmail"]}}}}

### spec.topicSpacesConfiguration.maximumClientSessionsPerAuthenticationName

`int32` · optional (explicit presence)

- default: `1`
- rule: {"int32":{"lte":100,"gte":1}}

### spec.topicSpacesConfiguration.maximumSessionExpiryInHours

`int32` · optional (explicit presence)

- default: `1`
- rule: {"int32":{"lte":8,"gte":1}}

### spec.topicSpacesConfiguration.routeTopicId

`string | valueFrom`

- references: AzureEventgridTopic (`status.outputs.topic_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureEventgridTopic, name: <that resource's name>, fieldPath: status.outputs.topic_id}} -- a bare string does not parse

### spec.topicSpacesConfiguration.dynamicRoutingEnrichments

`[]AzureEventgridNamespaceRoutingEnrichment`

### spec.topicSpacesConfiguration.dynamicRoutingEnrichments[].key

`string` · required

- rule: {"required":true,"string":{"minLen":"1","maxLen":"20"}}

### spec.topicSpacesConfiguration.dynamicRoutingEnrichments[].value

`string` · required

- rule: {"required":true,"string":{"minLen":"1","maxLen":"128"}}

### spec.topicSpacesConfiguration.staticRoutingEnrichments

`[]AzureEventgridNamespaceRoutingEnrichment`

### spec.topicSpacesConfiguration.staticRoutingEnrichments[].key

`string` · required

- rule: {"required":true,"string":{"minLen":"1","maxLen":"20"}}

### spec.topicSpacesConfiguration.staticRoutingEnrichments[].value

`string` · required

- rule: {"required":true,"string":{"minLen":"1","maxLen":"128"}}

### spec.tags

`map<string, string>`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureEventgridNamespace, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.namespace_id` | `string` |  |
| `status.outputs.namespace_name` | `string` |  |
| `status.outputs.identity_principal_id` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.identity.identityIds` | AzureUserAssignedIdentity | `status.outputs.identity_id` |
| `spec.topicSpacesConfiguration.routeTopicId` | AzureEventgridTopic | `status.outputs.topic_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureEventgridNamespaceTopic | `spec.namespaceId` | `status.outputs.namespace_id` |

## See Also

- [Overview](../README.md)
