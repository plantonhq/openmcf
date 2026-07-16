# AzureEventHubSchemaGroup -- Design Research

## The Resource

An Event Hubs schema group
(`Microsoft.EventHub/namespaces/schemagroups`) is a named collection in
the namespace's schema registry: every schema in the group shares a
serialization format and an evolution policy. The component maps onto
`azurerm_eventhub_namespace_schema_group` (azurerm v4.x,
`internal/services/eventhub/eventhub_namespace_schema_registry_resource.go`),
parity-verified against pulumi-azure v6
(`eventhub.NamespaceSchemaGroup`).

The registry is what lets producers and consumers exchange compact,
schema-referencing payloads: serializers register or resolve schemas
against the group at runtime -- the Azure SDK's schema-registry
serializers directly, and Kafka clients through the registry's
Kafka-compatible surface.

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `namespace_id` | `namespace_id` | The typed parent id; FK → AzureEventHubNamespace; ForceNew |
| `name` | `schema_group_name` | Required, ForceNew; 1-256 chars, alphanumeric ends, `[-._]` interior |
| `schema_compatibility` | enum | NONE/BACKWARD/FORWARD, required (unspecified rejected); ForceNew |
| `schema_type` | enum | AVRO/JSON, required; ForceNew |

Outputs: `schema_group_id`, `schema_group_name`.

## The All-ForceNew Verdict

This resource has NO update surface: Azure exposes no mutable
properties on a schema group, so every spec field is ForceNew. Any
change replaces the group, and the schemas registered inside it are
dropped with the old group. The spec documents this prominently --
operators should treat the group as append-only infrastructure and
choose name, policy, and format deliberately at creation.

## Recorded Skips (with reasons)

- **The "Unknown" schema type** -- Azure's API carries an "Unknown"
  placeholder value alongside Avro and Json; it is not a real
  serialization format and is deliberately not modeled. The enum
  offers the two formats a group can actually be created with.
- **The tier contract stays apply-time** -- the registry requires a
  STANDARD or higher namespace, which depends on the referenced
  namespace's live sku that validation cannot see. Documented on the
  spec; Azure rejects schema groups on BASIC verbatim.
- **No tags** -- ARM does not support tags on Event Hubs entities;
  the platform's identity tags live on the parent namespace.

## Policy Semantics Worth Knowing

- **BACKWARD**: new schemas can READ data written with older ones
  (delete fields, add optional fields) -- upgrade consumers first. The
  standard choice for analytics pipelines.
- **FORWARD**: old schemas can read data written with NEW ones (add
  fields, delete optional fields) -- upgrade producers first.
- **NONE**: no checking -- fine for single-team streams where
  producers and consumers deploy together.

## Operational Behavior Worth Knowing

- **Schema groups create in seconds** -- the namespace dominates any
  composed deploy.
- **Replacement is destructive**: a policy or format change looks like
  a one-line edit but drops every registered schema version with the
  old group.

## Composition

- `namespace_id` → `AzureEventHubNamespace.status.outputs.namespace_id`
- `schema_group_name` output ← serializer configuration (paired with
  the namespace's fully-qualified hostname)
