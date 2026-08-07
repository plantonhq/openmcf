# KubernetesClickHouse Guide

The judgment this guide carries: this is the catalog's ClickHouse, and
its two silent failure modes are the operator that must be watching this
namespace and the coordination service that replication quietly requires
but nothing enforces at deploy time.

## The operator, and its own-namespace default

KubernetesAltinityOperator is the registry prerequisite, and it watches
ONLY its own namespace unless widened — the
[operator-prerequisite pattern](../../_patterns/operator-prerequisite.md)
is the full story, and the [operator guide](../kubernetesaltinityoperator/GUIDE.md)
carries the watch-scope judgment. A ClickHouse cluster in an application
namespace the operator does not watch sits unreconciled with no error in
its own manifest.

## Replication needs coordination — wire it deliberately

Setting `replicas` above 1 (and any `ON CLUSTER` DDL) requires a
coordination service — ClickHouse Keeper or ZooKeeper — declared through
`spec.coordination` (the field doc on [reference.md](v1alpha1/reference.md) is the
authority). A replicated cluster proposed without it is a topology that
cannot actually replicate; name the coordination component in the
proposal alongside the shard/replica counts.

## Namespace ownership

The cluster shares its namespace with consumers or a co-located operator
— the multi-tenant case where `createNamespace: true` is wrong. Wire
`spec.namespace` to a dedicated KubernetesNamespace —
[namespace-ownership pattern](../../_patterns/namespace-ownership.md).

## On the diagram

Cluster and operator render as separate nodes with no edge between them
(the operator-prerequisite pattern); the shard×replica topology and the
coordination dependency live inside this spec and the coordination node.
An architecture showing a replicated ClickHouse with no coordination node
is showing a cluster that will not replicate.

## Pairs well with

- KubernetesAltinityOperator — required; watch-scope judgment in its
  guide.
- KubernetesNamespace — the shared namespace's owner (pattern above).
- KubernetesKafka — the common ingest path into ClickHouse for
  event-analytics architectures.
