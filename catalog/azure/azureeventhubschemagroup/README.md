# AzureEventHubSchemaGroup

A schema group in an Event Hubs namespace's schema registry: a named
collection of event schemas with a shared serialization format
(AVRO/JSON) and a compatibility policy governing evolution. Producers
and consumers exchange compact, schema-referencing payloads instead of
embedding schemas in every event -- serializers register or resolve
schemas against the group at runtime, including Kafka clients through
the registry's Kafka-compatible surface.

## When to Use

Use AzureEventHubSchemaGroup when you need:

- **Compact event payloads** -- events carry a schema reference, not
  the schema itself; the Azure SDK's schema-registry serializers
  resolve it at runtime
- **Safe schema evolution** -- the compatibility policy controls which
  changes the registry accepts as new schema versions within the group
- **A Confluent-style registry for Kafka clients** -- the registry
  exposes a Kafka-compatible surface, so Kafka serializer stacks work
  against it

## Key Configuration

- `namespace_id` -- the Event Hubs namespace whose registry holds the
  group, referenced from an AzureEventHubNamespace output (STANDARD or
  higher -- Azure rejects schema groups on BASIC at apply time)
- `schema_group_name` -- unique within the namespace, 1-256 characters;
  what serializers address
- `schema_compatibility` -- the evolution policy: NONE (no checking),
  BACKWARD (new schemas read old data; upgrade consumers first -- the
  standard for analytics pipelines), or FORWARD (old schemas read new
  data; upgrade producers first)
- `schema_type` -- AVRO (the registry's first-class format) or JSON

**Every field is ForceNew.** Azure exposes no mutable properties on a
schema group, so any change replaces it -- and the schemas registered
inside it are dropped with it. Choose the name, policy, and format
deliberately, and treat the group as append-only infrastructure.

## Composition

```yaml
namespaceId:
  valueFrom:
    kind: AzureEventHubNamespace
    name: telemetry-hubs
    fieldPath: status.outputs.namespace_id
```

Serializers address `status.outputs.schema_group_name` alongside the
namespace's fully-qualified hostname.

## Documentation

- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
