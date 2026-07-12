---
title: "Backward-Compatible Avro Group"
description: "This preset creates an Avro schema group with BACKWARD compatibility -- the standard evolution policy for analytics pipelines, where consumers upgrade first and must keep reading data written with..."
type: "preset"
rank: "01"
presetSlug: "01-backward-compatible-avro"
componentSlug: "event-hub-schema-group"
componentTitle: "Event Hub Schema Group"
provider: "azure"
icon: "package"
order: 1
---

# Backward-Compatible Avro Group

This preset creates an Avro schema group with BACKWARD compatibility --
the standard evolution policy for analytics pipelines, where consumers
upgrade first and must keep reading data written with older schemas.

## When to Use

- Producer/consumer fleets exchanging compact schema-referencing Avro
  payloads through the Azure SDK's schema-registry serializers
- Kafka-client estates using the registry's Kafka-compatible surface

## Key Configuration Choices

- **`schemaCompatibility: BACKWARD`** -- the registry rejects schema
  versions that would break existing readers (delete fields and add
  optional fields freely; upgrade consumers before producers)
- **Everything is ForceNew** -- the group has no update surface; a
  policy change replaces the group and drops its registered schemas,
  so choose the policy deliberately
- **STANDARD+ namespace** -- Azure rejects schema groups on BASIC at
  apply time

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-telemetry-hubs` | The AzureEventHubNamespace's Planton resource name | Your streaming composition |
| `telemetry-schemas` | The group name serializers address | Your schema taxonomy |
