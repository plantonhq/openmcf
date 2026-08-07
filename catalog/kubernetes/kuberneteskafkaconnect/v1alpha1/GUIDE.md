# KubernetesKafkaConnect Guide

The judgment this guide carries: on this platform Connect is
declarative-ONLY — anything created through its REST API gets reverted —
and the stock image runs no real integrations, so the plugin decision is
part of the architecture, not an afterthought.

## Declarative only — the REST API is not yours

The module pins Strimzi's `use-connector-resources` annotation, so
connectors are managed exclusively through KubernetesKafkaConnector
resources; connectors created or modified via the Connect REST API are
REVERTED by the operator (the reference page states the contract).
Propose pipes as resources — visible, reviewable nodes — never as
runbook curl commands.

## The stock image runs nothing real — pick a plugin arm

The stock image carries only the MirrorMaker 2 connectors: every real
integration (Debezium CDC, object-store sinks, search indexes) needs one
of the plugin arms — a prebuilt `image`, OCI `plugins`, or an in-cluster
`build` (the decision surface is on [reference.md](reference.md)). A
connector declaring a class the workers do not carry fails with
"class not found" — so the Connect cluster's plugin arm and its
connectors' classes are ONE decision, made together.

## Placement, operator, namespace

The Strimzi operator reconciles this kind
([operator-prerequisite pattern](../../../_patterns/operator-prerequisite.md);
watch judgment in the
[Strimzi operator guide](../../kubernetesstrimzikafkaoperator/v1alpha1/GUIDE.md)),
and KubernetesKafkaConnector resources must live in THIS cluster's
namespace — same place-bound contract as topics and users (their
reference pages state it). Wire `spec.namespace` through the Kafka
ecosystem's shared KubernetesNamespace
([namespace-ownership pattern](../../../_patterns/namespace-ownership.md)).

## On the diagram

Connect cluster and each connector render as nodes — the data pipeline
topology (what flows from where to where) is reviewable. REST-created
connectors would render as nothing, which is the other reason the
declarative contract is right.

## Pairs well with

- KubernetesKafka — the cluster whose topics the pipes read and write.
- KubernetesKafkaConnector — one per pipe, in this namespace.
- KubernetesStrimziKafkaOperator — required machinery.
