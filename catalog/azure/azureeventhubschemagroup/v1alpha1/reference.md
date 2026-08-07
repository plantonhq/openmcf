# AzureEventHubSchemaGroup

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureEventHubSchemaGroupSpec** defines a schema group in an Event Hubs
namespace's schema registry: a named collection of event schemas with a
shared serialization format and compatibility (evolution) policy.

The schema registry lets producers and consumers exchange compact,
schema-referencing payloads instead of embedding schemas in every event:
serializers register or resolve schemas against the group at runtime
(the Azure SDK's schema-registry serializers, and Kafka clients through
the registry's Kafka-compatible surface). The compatibility policy is
what makes schema EVOLUTION safe -- it controls which changes the
registry accepts as new schema versions within the group.

**Tier contract enforced by Azure at apply time**: the schema registry
requires a STANDARD or higher namespace (BASIC namespaces reject schema
groups).

**This resource has no update surface**: every field is fixed at
creation -- Azure exposes no mutable properties on a schema group, so
any change replaces it (the registered schemas inside it are dropped
with it; treat the group as append-only infrastructure).

**ForceNew fields**: all of them (`namespace_id`, `schema_group_name`,
`schema_compatibility`, `schema_type`).

## Example

```yaml
# Offline-plan manifest: a schema group exercising the enum wire maps
# (BACKWARD -> Backward, AVRO -> Avro).
apiVersion: azure.planton.dev/v1alpha1
kind: AzureEventHubSchemaGroup
metadata:
  name: test-eh-schema-group
spec:
  namespaceId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.EventHub/namespaces/hack-eventhub-ns
  schemaGroupName: order-events
  schemaCompatibility: BACKWARD
  schemaType: AVRO
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespaceId` | `string \| valueFrom` | yes |  | AzureEventHubNamespace (`status.outputs.namespace_id`) |
| `spec.schemaGroupName` | `string` | yes |  |  |
| `spec.schemaCompatibility` | `enum` |  |  |  |
| `spec.schemaType` | `enum` |  |  |  |

## Field Details

### spec.namespaceId

`string | valueFrom` · required

The Event Hubs namespace whose registry holds the group, by ARM ID.
References an AzureEventHubNamespace's namespace_id output. Fixed at
creation.

- references: AzureEventHubNamespace (`status.outputs.namespace_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureEventHubNamespace, name: <that resource's name>, fieldPath: status.outputs.namespace_id}} -- a bare string does not parse

### spec.schemaGroupName

`string` · required

The group's name -- unique within the namespace, 1-256 characters.
Starts and ends with a letter or number; letters, numbers, periods,
hyphens, and underscores in between. Serializers address the group
by this name.

**ForceNew**: renaming replaces the group and drops its registered
schemas.

- rule: schema_group_name must start and end with a letter or number and may contain letters, numbers, periods, hyphens, and underscores (max 256 characters)
- rule: {"required":true,"string":{"minLen":"1","maxLen":"256"}}

### spec.schemaCompatibility

`enum`

The evolution policy: which changes the registry accepts as new
versions of an existing schema in this group.

- NONE: any change is accepted -- no compatibility checking. Fine
  for single-team streams where producers and consumers deploy
  together.
- BACKWARD: new schemas can READ data written with older ones
  (delete fields, add optional fields) -- upgrade consumers first.
  The standard choice for analytics pipelines.
- FORWARD: old schemas can read data written with NEW ones (add
  fields, delete optional fields) -- upgrade producers first.

**ForceNew**: fixed at creation.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_event_hub_schema_compatibility_unspecified` -- Not specified -- invalid; choose an explicit policy.
- `NONE` -- No compatibility checking -- any change is accepted.
- `BACKWARD` -- New schemas read old data; upgrade consumers first.
- `FORWARD` -- Old schemas read new data; upgrade producers first.

### spec.schemaType

`enum`

The serialization format for every schema in the group. AVRO is the
Event Hubs registry's first-class format; JSON covers JSON-Schema
payloads.

**ForceNew**: fixed at creation.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_event_hub_schema_type_unspecified` -- Not specified -- invalid; choose an explicit format. (Azure's API also carries an "Unknown" placeholder type; it is not a real format and is deliberately not modeled.)
- `AVRO` -- Apache Avro -- the registry's first-class schema format.
- `JSON` -- JSON Schema documents.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureEventHubSchemaGroup, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.schema_group_id` | `string` | The Azure Resource Manager ID of the schema group. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.EventHub/namespaces/{ns}/schemagroups/{name} |
| `status.outputs.schema_group_name` | `string` | The group's name -- what schema-registry serializers address. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespaceId` | AzureEventHubNamespace | `status.outputs.namespace_id` |

## See Also

- [Overview](../README.md)
