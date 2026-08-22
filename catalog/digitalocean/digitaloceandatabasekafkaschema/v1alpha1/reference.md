# DigitalOceanDatabaseKafkaSchema

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `digital-ocean.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

DigitalOceanDatabaseKafkaSchemaSpec models the full
digitalocean_database_kafka_schema_registry resource surface: one schema
subject registered in a DigitalOcean managed Kafka cluster's schema
registry.

EVERY field is create-only. There is no update path in the provider:
any change -- including evolving the schema definition, even a
whitespace-only reformat (the definition is compared verbatim, never
normalized) -- destroys the subject and re-registers it, which DROPS all
previously registered versions of the subject. Treat schema evolution as
a deliberate replacement, never a casual edit.

## Example

```yaml
# Reference manifests for DigitalOceanDatabaseKafkaSchema --
# protovalidate-valid, embedded as the reference page's Example block, and
# the documents the offline tofu plans render. Two documents: an Avro
# record subject and a JSON-schema subject.
apiVersion: digital-ocean.planton.dev/v1alpha1
kind: DigitalOceanDatabaseKafkaSchema
metadata:
  name: orders-value-schema
spec:
  # Literal cluster UUID; use valueFrom to reference a
  # DigitalOceanDatabaseCluster resource instead.
  cluster:
    value: aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee
  subjectName: orders-value
  schemaType: avro
  schema: '{"type":"record","name":"Order","fields":[{"name":"id","type":"string"},{"name":"amountCents","type":"long"}]}'
---
apiVersion: digital-ocean.planton.dev/v1alpha1
kind: DigitalOceanDatabaseKafkaSchema
metadata:
  name: shipment-events-schema
spec:
  cluster:
    value: aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee
  subjectName: shipment-events-value
  schemaType: json
  schema: '{"$schema":"http://json-schema.org/draft-07/schema#","type":"object","properties":{"shipmentId":{"type":"string"},"status":{"type":"string"}},"required":["shipmentId","status"]}'
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.cluster` | `string \| valueFrom` | yes |  | DigitalOceanDatabaseCluster (`status.outputs.cluster_id`) |
| `spec.subjectName` | `string` | yes |  |  |
| `spec.schemaType` | `string` | yes |  |  |
| `spec.schema` | `string` | yes |  |  |

## Field Details

### spec.cluster

`string | valueFrom` · required

The Kafka database cluster whose schema registry the subject is
registered in. Use a literal cluster UUID or a reference to a
DigitalOceanDatabaseCluster resource (the cluster must run the kafka
engine). Changing it replaces the subject.

- references: DigitalOceanDatabaseCluster (`status.outputs.cluster_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: DigitalOceanDatabaseCluster, name: <that resource's name>, fieldPath: status.outputs.cluster_id}} -- a bare string does not parse

### spec.subjectName

`string` · required

Name of the schema subject. Unique within the cluster's registry; the
subject name IS the API identity. Changing it replaces the subject.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.schemaType

`string` · required

Schema definition language: avro, json, or protobuf (case-sensitive).
Changing it replaces the subject.

- rule: {"required":true,"string":{"in":["avro","json","protobuf"]}}

### spec.schema

`string` · required

The schema definition itself, in the language schema_type names (e.g.
an Avro record document as JSON). Stored and compared verbatim --
formatting is significant. Changing it replaces the subject and drops
all prior versions.

- rule: {"required":true,"string":{"minLen":"1"}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: DigitalOceanDatabaseKafkaSchema, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.cluster_id` | `string` | UUID of the Kafka database cluster whose registry holds the subject. |
| `status.outputs.subject_name` | `string` | Name of the registered schema subject (its API identity within the cluster's registry). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.cluster` | DigitalOceanDatabaseCluster | `status.outputs.cluster_id` |

## See Also

- [Overview](../README.md)
