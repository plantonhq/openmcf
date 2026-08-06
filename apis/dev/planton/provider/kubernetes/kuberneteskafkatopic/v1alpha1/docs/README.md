# KubernetesKafkaTopic: Research and Design

## Introduction

A KafkaTopic is the Strimzi custom resource that manages a single topic
within a Kafka cluster. This component creates one KafkaTopic, named
after the resource (`metadata.name`). The topic itself is created by the
cluster's TOPIC OPERATOR, not by the IaC modules — the modules apply the
declaration; the operator does the reconciling.

## The Upstream Model

Strimzi's topic operator runs per cluster (part of the entity operators
a KubernetesKafka deploys by default) and performs the corresponding
Kafka operation whenever a KafkaTopic is created, deleted, or changed.
Configuration flows ONE WAY: from the resource to the topic. Changes
made to a managed topic outside its resource are reverted at the next
reconciliation. Topics without a KafkaTopic resource are left alone —
the operator manages only what is declared.

When more than one KafkaTopic names the same Kafka topic (via
`spec.topicName`), the oldest resource manages the topic; newer ones
fail with a resource conflict and report NotReady.

## The Placement Contract

The topic operator watches a SINGLE namespace — the Kafka cluster's own
— and selects resources by the `strimzi.io/cluster` label. Both halves
are load-bearing:

- A KafkaTopic in any other namespace is accepted by the API server and
  then silently never reconciled — no error, no topic.
- A KafkaTopic whose label names a cluster that does not exist in that
  namespace is equally invisible to the operator.

The spec therefore requires `namespace` (which MUST be the cluster's
namespace) and `kafka_cluster` (rendered as the label); the modules
stamp the label unconditionally.

## Naming

`metadata.name` is the topic name by default. Kafka permits names
Kubernetes metadata cannot carry — `.`, `_`, and uppercase (e.g.
`orders.v1_DLQ`) — so `spec.topic_name` overrides the Kafka-side name
while the resource keeps a valid Kubernetes name. Once the topic is
created, its name cannot be changed (a `topicName` rename is detected
and rejected).

The spec validates the Kafka rules: up to 249 characters of
alphanumerics, `.`, `_` and `-`. One hazard the format rule cannot
catch: names differing only by `.` vs `_` COLLIDE in Kafka's internal
metrics — avoid coexisting pairs like `orders.v1` and `orders_v1`.

## Partitions and Replicas

- **Partitions grow, never shrink.** Kafka has no partition shrink; the
  topic operator rejects a decrease at reconcile time
  (`PartitionDecreaseException`, resource NotReady). Increasing
  partitions changes key-to-partition mapping for keyed topics — plan
  counts up front where partitioning is semantic. Empty = the cluster's
  `num.partitions` default.
- **Replicas are capped by the broker count.** A replication factor
  above the cluster's broker count is rejected at reconcile time — the
  resource reports NotReady and nothing is created. Changing `replicas`
  on an existing topic requires Cruise Control on the cluster (Strimzi
  submits it as a replicas-change operation). Empty = the cluster's
  `default.replication.factor`. The production norm is 3, paired with
  `min.insync.replicas: "2"` — survives one broker loss without losing
  acknowledged writes.

## Deletion Deletes Data

Deleting the resource deletes the TOPIC AND ITS DATA — the topic
operator propagates deletion to Kafka. Two cluster-side switches govern
this: `delete.topic.enable` must be `true` in the cluster's config (the
Kafka default), and the operator's finalizer
(`strimzi.io/topic-operator`, on by default) holds the resource until
Kafka-side deletion completes. A namespace stuck terminating because the
operator is gone is the known failure mode of that finalizer — Strimzi
documents removing finalizers manually as the escape hatch.

## Auto-Created Topics vs Declared Topics

Brokers auto-create topics by default (`auto.create.topics.enable`)
when a client touches a nonexistent name, with default configuration.
An auto-created topic that races a KafkaTopic declaration lands first,
and the operator's follow-up reconfiguration can fail — for example
when the declared partition count is below the broker default. Strimzi
recommends disabling auto-creation on clusters where topics are
declared; the boundary language on the kind's README carries this.

## Deliberately Unmodeled

Stated timelessly, for the record: the KafkaTopic CRD's spec has
exactly four fields — `topicName`, `partitions`, `replicas`, `config` —
and the Planton spec models all four. There is nothing on the CRD left
unmodeled. What stays outside the kind is operational machinery, not
spec surface:

- **Partition reassignment / rebalancing** — `KafkaRebalance` is an
  operational verb against the cluster (and requires Cruise Control),
  not a property of a topic.
- **Unmanaging and finalizer surgery** — Strimzi's
  `strimzi.io/managed: "false"` annotation flow and finalizer removal
  are operational procedures, applied with kubectl when needed.

## Engine Mechanics

- **Pulumi**: one untyped CustomResource (`kafka.strimzi.io/v1`
  KafkaTopic) — the Strimzi CRDs type `config` with
  `x-kubernetes-preserve-unknown-fields`, which crd2pulumi cannot
  carry, so no generated package is shipped for the Kafka family (the
  same ruling as the KubernetesKafka module).
- **Terraform**: one `kubectl_manifest` (alekc/kubectl) — no cluster
  connection at plan time, so topics plan before the Strimzi CRDs exist
  (single-run infra charts, offline plan proofs). The locals render the
  same shape as the Pulumi module; twins kept in lockstep.
- **No namespace resource, deliberately**: the namespace belongs to the
  KubernetesKafka resource's lifecycle — this kind only places into it.
- **No await machinery**: reconciliation belongs to the topic operator,
  not to applying the resource. `kubectl get kafkatopic` shows Ready /
  NotReady with a teaching message.

## Outputs

`namespace` and `topic_name` — the actual Kafka topic name
(`spec.topic_name` when set, else `metadata.name`), resolved
identically in both engines so the handle can never drift from what the
operator creates. The bootstrap endpoint deliberately comes from the
KubernetesKafka resource's own outputs.

## E2E

The kind-cluster lane deploys a single-node Kafka fixture and asserts
reconciliation: a minimal topic and a full-surface topic (name override
with dots and underscores, compaction config, explicit sizing) — with
`replicas: 1`, because the fixture has one broker and the operator
rejects a higher replication factor.
