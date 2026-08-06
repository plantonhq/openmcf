# Schema-Validated Topic

The contract-enforced event stream: every published message is validated
against a shared schema before it enters the topic.

## What this preset creates

A topic attached to a referenced `GcpPubSubSchema` with JSON encoding.
Non-conforming publishes fail immediately at the producer with
`INVALID_ARGUMENT` — contract violations surface where they can be fixed,
never as processing failures deep inside consumers.

## Prerequisites

- A `GcpPubSubSchema` named `order-events-schema` (see the
  `GcpPubSubSchema` presets for an Avro event contract starting point).

## Remix ideas

- Switch `encoding` to `BINARY` for protobuf-typed schemas and compact
  wire payloads.
- Pair with a BigQuery-delivery subscription using `useTopicSchema` —
  the schema then drives the table's column mapping end to end.
- Evolve the contract by committing new revisions on the schema
  resource; the topic picks them up without any change here.
