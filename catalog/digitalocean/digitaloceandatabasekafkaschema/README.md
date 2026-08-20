# DigitalOcean Database Kafka Schema

Built for 100% parity with the Terraform DigitalOcean provider's `digitalocean_database_kafka_schema_registry` resource at the pinned provider version.

## What this component models

One schema subject registered in a DigitalOcean managed Kafka cluster's schema registry: the subject's name, its definition language (Avro, JSON Schema, or Protobuf), and the definition itself.

The component covers the provider's full argument surface:

- `cluster` -- the owning Kafka cluster, wired by reference (or a literal cluster UUID)
- `subject_name` -- the registry subject (create-only)
- `schema_type` -- `avro`, `json`, or `protobuf` (create-only)
- `schema` -- the definition document (create-only; compared verbatim)

## Quick start

```yaml
apiVersion: digital-ocean.planton.dev/v1alpha1
kind: DigitalOceanDatabaseKafkaSchema
metadata:
  name: orders-value-schema
spec:
  cluster:
    valueFrom:
      kind: DigitalOceanDatabaseCluster
      name: my-kafka-cluster
      fieldPath: status.outputs.cluster_id
  subjectName: orders-value
  schemaType: avro
  schema: '{"type":"record","name":"Order","fields":[{"name":"id","type":"string"}]}'
```

Deploy with either provisioner; both produce identical resources and outputs.

## Outputs

| Output | Description |
|---|---|
| `cluster_id` | UUID of the Kafka cluster whose registry holds the subject |
| `subject_name` | The registered subject's name (its API identity) |

## Behavior worth knowing

- **EVERY field is create-only.** The provider has no update path: any change -- including schema evolution and even a whitespace-only reformat of the definition -- destroys the subject and re-registers it, which **drops all previously registered versions**. Treat evolution as a deliberate replacement.
- **The definition is compared verbatim.** There is no JSON normalization; keep the document byte-stable in your manifest.
- **Import is excluded.** The provider's importer is broken at the pinned version (it never restores the subject name); the import map records the exclusion and its re-evaluate trigger.

## Module layout

- `iac/tf/` -- OpenTofu/Terraform module (provider pinned `~> 2.99`)
- `iac/pulumi/` -- Pulumi module (Go, pulumi-digitalocean SDK)
- Both engines wire the same spec fields and export the same outputs; behavioral parity is the contract.
