---
title: "Protobuf Binary Contract"
description: "A protobuf-typed schema for high-volume telemetry where compact binary encoding matters."
type: "preset"
rank: "02"
presetSlug: "02-protobuf-binary-contract"
componentSlug: "pubsub-schema"
componentTitle: "Pub/Sub Schema"
provider: "gcp"
icon: "package"
order: 2
---

# Protobuf Binary Contract

A protobuf-typed schema for high-volume telemetry where compact binary
encoding matters.

## What this preset creates

A Pub/Sub schema named `device-telemetry` defined as a single proto3
message. Topics that attach it with `encoding: BINARY` validate the raw
protobuf bytes of every published message; `encoding: JSON` validates the
protojson form instead.

## When to choose PROTOCOL_BUFFER over AVRO

- Publishers already produce protobuf-serialized payloads (shared proto
  definitions across services).
- Wire size and serialization cost dominate — binary protobuf is the
  compact option.
- Note the trade-off: BigQuery delivery's `useTopicSchema` and Cloud
  Storage Avro export work most naturally with AVRO schemas; protobuf
  contracts usually pair with pull or push subscribers that own decoding.

## Remix ideas

- Evolve the message by committing new definitions — additive,
  tag-stable changes keep old publishers valid (each commit is a schema
  revision; a schema retains up to 20).
- Keep exactly one top-level message per schema — that is the contract
  Pub/Sub validates against.
