# KubernetesKafkaMirrorMaker2 Guide

The judgment this guide carries: this is the catalog's MIGRATION ON-RAMP
for Kafka — when a user asks to move off Confluent Cloud, MSK, or any
external Kafka, the answer is not an export script; it is this kind
mirroring into a KubernetesKafka until consumers cut over with their
offsets intact.

## The migration workflow it anchors

1. Stand up the destination
   [KubernetesKafka](../../kuberneteskafka/v1/GUIDE.md) (plus its
   operator).
2. Declare this mirror: source = the running external cluster, target =
   the new one. It replicates topics, records, AND consumer positions
   continuously — checkpointing translates source offsets to target
   offsets (the mechanics are on [reference.md](reference.md)).
3. Move producers, then cut consumers over to the target with offsets
   intact; retire the mirror when the source drains.

Propose it whenever the ask smells like "migrate/replicate/DR my Kafka" —
the alternative people reach for (snapshot-and-replay) loses consumer
positions, which is exactly what this kind preserves.

## Shape and placement

One TARGET cluster per mirror, one or more sources; it runs on the
Strimzi operator like the rest of the family
([operator-prerequisite pattern](../../../patterns/operator-prerequisite.md))
and lives in the target cluster's namespace, wired through the shared
KubernetesNamespace
([namespace-ownership pattern](../../../patterns/namespace-ownership.md)).

## On the diagram

The mirror renders as its own node between source and target — a
migration in progress is VISIBLE in the architecture, which is exactly
what a reviewer wants during a cutover window.

## Pairs well with

- KubernetesKafka — the target (and, when in-cluster, the source).
- KubernetesStrimziKafkaOperator — required machinery.
- KubernetesKafkaUi — watching lag and mirrored topics during cutover.
