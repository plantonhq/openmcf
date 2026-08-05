# KubernetesKafka Guide

The judgment this guide carries: an agent-composed Kafka goes silently
inert in two ways — an operator that is not watching the cluster's
namespace, and topics or users declared where the entity operators never
look. Neither mistake errors at deploy time.

## The operator must watch THIS namespace

KubernetesStrimziKafkaOperator is the registry prerequisite, but
installing it is not enough: by default the operator watches its own
namespace only, so a Kafka cluster in an application namespace needs the
operator's watch scope to cover it. Placement judgment and watch-scope
wiring: the
[operator's guide](../../kubernetesstrimzikafkaoperator/v1/GUIDE.md).
Note the coupling never appears as a reference edge in the manifests —
reviewers must check the watch scope deliberately.

## Topics and users are architecture, and they are place-bound

Declare topics and users as KubernetesKafkaTopic / KubernetesKafkaUser
resources — each is a visible, reviewable node with its own retention and
ACL surface — and declare them IN THE KAFKA CLUSTER'S NAMESPACE: its
entity operators watch only there (the coupling is on
[reference.md](reference.md)). A topic manifest in the app's namespace is
a resource nothing reconciles.

## Namespace ownership

A Kafka namespace is inherently shared — the cluster, its topics, its
users, and usually its client workloads live together. That is exactly
the multi-tenant case where `createNamespace: true` on any one of them is
wrong: use a dedicated KubernetesNamespace and wire everyone's
`spec.namespace` through it —
[namespace-ownership pattern](../../../patterns/namespace-ownership.md).

## On the diagram

Cluster, topics, and users each render as nodes in the shared namespace —
the data topology becomes something a reviewer can read. Topics
auto-created by client applications render as nothing.

## Pairs well with

- KubernetesStrimziKafkaOperator — required; watch scope judgment above.
- KubernetesKafkaTopic / KubernetesKafkaUser — the data topology, in this
  cluster's namespace.
- KubernetesKarapace — Confluent-API-compatible schema registry when
  producers and consumers need enforced schemas.
- KubernetesKafkaUi — the web console for the cluster.
