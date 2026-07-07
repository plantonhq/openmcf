---
title: "Pub/Sub Schema"
description: "Pub/Sub Schema deployment documentation"
icon: "package"
order: 100
componentName: "gcppubsubschema"
---

# GCP Pub/Sub Schema

Creates a Pub/Sub schema — the message contract publishers and subscribers agree on. One schema can validate messages on many topics: each topic attaches it by reference, so the event contract is evolved in one place and every attached topic enforces it at publish time.

## What Gets Created

A project-level `google_pubsub_schema` resource holding a typed message definition — an Avro record (JSON) or a protobuf message. Topics attach it through their `schemaSettings.schema` field; from then on Pub/Sub rejects any published message that does not conform to the schema.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **A GCP project** (the Pub/Sub API is enabled automatically)

## Quick Start

Create a file `schema.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpPubSubSchema
metadata:
  name: order-events-schema
spec:
  projectId:
    value: my-gcp-project
  schemaName: order-events
  type: AVRO
  definition: |
    {
      "type": "record",
      "name": "OrderEvent",
      "fields": [
        {"name": "order_id", "type": "string"},
        {"name": "amount_cents", "type": "long"}
      ]
    }
```

Deploy:

```shell
planton apply -f schema.yaml
```

Attach it from a topic:

```yaml
# In a GcpPubSubTopic spec:
schemaSettings:
  schema:
    valueFrom:
      kind: GcpPubSubSchema
      name: order-events-schema
      fieldPath: status.outputs.schema_id
  encoding: JSON
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `schemaName` | `string` | The GCP resource name of the schema. Immutable — renaming replaces the resource. | 3-255 chars; starts with a letter; letters, digits, `-_.~+%`; must not begin with the reserved `goog` prefix |
| `type` | `string` | The definition language. | `AVRO` or `PROTOCOL_BUFFER` |
| `definition` | `string` | The schema text (Avro JSON or a single protobuf message). Changing it commits a new revision in place. | Required |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `projectId` | `StringValueOrRef` | GCP project for the schema. Omitted rides the provider's default project; reference a `GcpProject` to compose |

## Revision Lifecycle

Definition changes never replace the schema — each change commits a new **revision** (a schema retains up to 20; beyond that, old revisions must be deleted manually). Attached topics accept messages conforming to any available revision, so evolve additively: add fields with defaults in Avro, add tagged fields in protobuf, and never repurpose existing ones.

## Stack Outputs

| Output | Description |
|--------|-------------|
| `schema_id` | Fully qualified schema path (`projects/{p}/schemas/{name}`) — what a topic's `schemaSettings.schema` consumes |
| `schema_name` | The short name of the schema |

## Operational Notes

- **Deletion ordering**: deleting a schema that topics still reference leaves them validating against the `_deleted-schema_` sentinel — publishes fail. Detach or recreate topics first.
- **AVRO pairs with delivery integrations**: BigQuery delivery (`useTopicSchema`) and Cloud Storage Avro export derive their layout from an Avro topic schema.

## Related Components

- [GcpPubSubTopic](../pubsub-topic/) — attaches this schema for publish-time validation
- [GcpPubSubSubscription](../pubsub-subscription/) — consumes schema-validated messages
