# KubernetesStrimziKafkaOperator Guide

The judgment this guide carries: installing this operator is necessary
but not sufficient — its WATCH SCOPE decides which namespaces' Kafka
clusters actually get reconciled, and the default (its own namespace
only) is almost never what a platform installation wants.

## Set the watch scope deliberately

The platform shape is one operator in the shared-cluster chart with
`watch.anyNamespace: true`, so every application environment can declare
Kafka clusters without touching the shared layer again. Use
`watch.namespaces` (a fenced list) when team isolation genuinely requires
it — accepting that every new Kafka-bearing namespace then needs a
shared-chart change. Leaving the watch block empty confines the operator
to its own namespace: correct only when the Kafka clusters live there
too. The scope vocabulary and multi-install rules are on
[reference.md](reference.md).

## Once per cluster, in the shared chart

One operator reconciles any number of Kafka clusters. Declare it once in
the shared-cluster chart; application environments declare
KubernetesKafka clusters, never their own operator. (A second install is
a special, fenced shape — its constraints are on the reference page.)

## Namespace ownership — the infra exception

A dedicated single-tenant namespace with `createNamespace: true` is the
normal shape here — the
[namespace-ownership pattern](../../../patterns/namespace-ownership.md)'s
sole-tenant case.

## On the diagram

The operator renders in the shared-cluster layer, but nothing draws an
edge from a Kafka cluster to it — the dependency is enforced by the
platform's prerequisite metadata, not by a manifest reference. When
reviewing an architecture with Kafka in it, verify the operator node
exists and its watch scope covers the cluster's namespace; the diagram
alone will not surface a mismatch.

## Pairs well with

- KubernetesKafka — the clusters this operator reconciles (see its
  [guide](../../kuberneteskafka/v1/GUIDE.md)).
- KubernetesKafkaTopic / KubernetesKafkaUser — reconciled by each
  cluster's entity operators, which this operator manages.
