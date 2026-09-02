# DigitalOcean Kafka Schema

Registers a schema subject (Avro, JSON Schema, or Protobuf) in a DigitalOcean managed Kafka cluster's schema registry, so producers and consumers agree on message structure. Every field is create-only: the provider has no update path, so any change -- including a whitespace-only reformat of the definition, which is compared verbatim -- destroys the subject and re-registers it, dropping all previously registered versions. The owning cluster is wired by reference or supplied as a literal UUID.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Schema Registry Subject** -- the named subject on the referenced cluster's registry, carrying your schema definition

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Kafka Database Cluster** -- a DigitalOceanDatabaseCluster running the `kafka` engine (the registry exists only on Kafka clusters).

### DigitalOcean Account

- Nothing beyond the cluster: schema subjects are free API objects on it.

## Deploy

### Console

Open the deployment store, find **DigitalOcean Kafka Schema**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Avro Event Schema** preset in the [Presets](#presets) tab to register a topic's value schema under the `<topic>-value` naming convention.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: digital-ocean.planton.dev/v1alpha1
kind: DigitalOceanDatabaseKafkaSchema
metadata:
  name: orders-value-schema
  org: acme-corp
  env: prod
spec:
  cluster:
    value: "2f6a8f0e-3b1c-4c8e-9f2d-7a5b4c3d2e1f"
  subjectName: orders-value
  schemaType: avro
  schema: '{"type":"record","name":"Order","namespace":"com.acme.orders","fields":[{"name":"id","type":"string"},{"name":"amountCents","type":"long"}]}'
```

```shell
planton apply -f kafka-schema.yaml
```

This registers the `orders-value` subject on the referenced Kafka cluster's registry with a two-field Avro record as its founding schema. A Stack Job tracks the provisioning in real time.

### InfraChart

When the Kafka cluster deploys in the same InfraPipeline, wire the subject to it with ValueFromRef instead of a literal UUID:

```yaml
spec:
  cluster:
    valueFrom:
      kind: DigitalOceanDatabaseCluster
      name: events-kafka
      fieldPath: status.outputs.cluster_id
```

The InfraPipeline resolves the dependency graph, deploys the cluster first, then registers the subject in its registry.

## Key Configuration

These are the most important decisions when configuring a Kafka schema subject. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Every change is a replacement** -- The provider has no update path: changing `schema`, `schemaType`, `subjectName`, or `cluster` destroys the subject and re-registers the new document as version 1, dropping every previously registered version. Consumers that pin older schema versions lose them the moment the replacement lands.

**The definition is compared verbatim** -- No JSON normalization, no key reordering: a whitespace-only reformat counts as a change and triggers the destroy-and-drop above. Keep the manifest's `schema` string byte-stable -- single-line, machine-formatted, never hand-prettified after the fact.

**Founding schema, not evolution channel** -- If your consumers rely on registry-mediated compatibility across versions, do not evolve schemas through this resource. Evolve them through your producers' registry client (which appends versions) and use this resource only to declare the founding schema of a subject.

**Subject naming** -- `subjectName` is the subject's API identity, unique within the cluster's registry. Follow the registry's `<topic>-value` (and `<topic>-key`) convention so serializer libraries resolve the schema automatically from the topic name.

**Schema language** -- `schemaType` accepts exactly `avro`, `json`, or `protobuf`, lowercase and case-sensitive. Changing it later replaces the subject.

**Compatibility level is not here** -- The registry's subject compatibility level (BACKWARD, FULL, and so on) has no surface in this resource; it stays whatever the registry defaults to or whatever was set out-of-band.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **DigitalOceanDatabaseCluster** | `cluster` | `status.outputs.cluster_id` |

The referenced cluster must run the `kafka` engine; a literal cluster UUID is accepted in place of the reference.

### What This Component Provides

After provisioning, `status.outputs` carries only the subject's identity pair -- `cluster_id` and `subject_name` -- both echoes of resolved inputs. The registry's internal numeric schema id is discarded by the provider and deliberately not exported. Producers and consumers fetch the schema by subject name through the cluster's registry endpoint, authenticating with the cluster's connection outputs and user credentials -- there is no output here for downstream Cloud Resources to wire.

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Avro value schema for an event topic** -- an Avro record registered under `<topic>-value` so producer and consumer serializers resolve it automatically. The default shape for event pipelines. Start from the **Avro Event Schema** preset.

**Strict JSON contract** -- a JSON Schema closed to unknown properties with enumerated fields, turning a loose JSON pipeline into a validated contract without moving to Avro. Start from the **JSON Contract Schema** preset.

## Works With

- [**DigitalOcean Database Cluster**](/cloud-catalog/digital-ocean-database-cluster) -- the Kafka-engine cluster whose registry holds the subject
- [**DigitalOcean Kafka Topic**](/cloud-catalog/digital-ocean-database-kafka-topic) -- the topic whose messages the subject describes, paired through the `<topic>-value` naming convention
- [**DigitalOcean Database User**](/cloud-catalog/digital-ocean-database-user) -- credentials producers and consumers use to reach the cluster and its registry
