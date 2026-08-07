# KubernetesRabbitMq Guide

The judgment this guide carries: pick this over Kafka by workload shape,
size it on an odd replica count from the start, and never copy the
operator-generated credentials — the cluster cannot be shrunk later, so
the initial shape matters more than usual.

## RabbitMQ or Kafka?

This is the task-queue / work-distribution / RPC broker; for append-only
event STREAMING at scale, KubernetesKafka is the kind (the spec doc on
[reference.md](v1alpha1/reference.md) states the split). Choosing by habit rather
than workload shape is the common mistake — a streaming firehose on
RabbitMQ and a request/reply queue on Kafka both fight their broker.

## Odd replicas, and you cannot scale down

Production runs an ODD replica count (3, 5, 7) so quorum queues and the
Raft metadata store survive node loss; 2 replicas lose availability on
any single failure. The operator does NOT support scaling down (removed
brokers strand their queue replicas) — sizing down means migrating to a
new cluster. Size for the ceiling you expect, because shrinking is a
migration, not a config change.

## Credentials are operator-generated — wire, never copy

Admin credentials live in the operator-generated `<name>-default-user`
Secret (the naming contract is on [reference.md](v1alpha1/reference.md)) — consume
them by reference; they never pass through this spec.

## Operator, namespace, and the destroy warning

KubernetesRabbitMqOperator is the registry prerequisite
([operator-prerequisite pattern](../../_patterns/operator-prerequisite.md));
its guide carries the all-namespaces default AND the critical warning
that destroying the operator cascade-deletes every RabbitMQ cluster. The
cluster shares its namespace with consumers — wire `spec.namespace` to a
dedicated KubernetesNamespace, not `createNamespace: true`
([namespace-ownership pattern](../../_patterns/namespace-ownership.md)).

## Pairs well with

- KubernetesRabbitMqOperator — required; the destroy-safety warning is in
  its guide.
- KubernetesNamespace — the shared namespace's owner.
- Application workloads — connect via the exported endpoints and the
  generated credentials Secret.
