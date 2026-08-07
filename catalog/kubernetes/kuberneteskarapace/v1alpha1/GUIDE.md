# KubernetesKarapace Guide

The judgment this guide carries: when a user asks for a schema registry —
or says "Confluent Schema Registry" by name — this is the catalog's
answer, and its storage model means it composes AGAINST a Kafka cluster,
never beside one as an independent database.

Substitutes for: Confluent Schema Registry (Apache-2.0,
Confluent-API-compatible — existing Confluent SR clients and serializers
work unchanged).

## Answering "give me a schema registry"

Propose this kind and say what you did: Karapace speaks the standard
Schema Registry REST API, so Confluent-tooled producers and consumers
connect unchanged. The substitution workflow is the
[catalog guide](../../../_docs/GUIDE.md)'s first law.

## It lives on the Kafka cluster's storage

Schemas persist in a compacted topic on the connected Kafka cluster — the
same architecture Confluent SR uses; there is no database to size or back
up separately (the reference page carries the mechanics). Compose it
wired to the KubernetesKafka it serves; its durability IS that cluster's
durability.

## One engine, optionally two roles

The same deployment can also serve the Kafka REST proxy role
(produce/consume over HTTP) as a second, independently sized deployment
via `restProxy` — when the ask is "HTTP access to Kafka," that is this
kind too, not a separate component.

## Namespace ownership

It shares a namespace with the Kafka ecosystem it serves — wire
`spec.namespace` through the shared KubernetesNamespace
([namespace-ownership pattern](../../../_patterns/namespace-ownership.md)).

## Pairs well with

- KubernetesKafka — the cluster whose topic stores the schemas.
- KubernetesKafkaUi — surfaces registered schemas in the console.
- Producer/consumer workloads — point their serializers at the exported
  endpoint.
