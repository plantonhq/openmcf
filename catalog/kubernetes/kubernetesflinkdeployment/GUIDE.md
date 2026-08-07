# KubernetesFlinkDeployment Guide

The judgment this guide carries: one field decides whether this is a
job-scoped cluster or a shared runtime, and stateful upgrade modes
silently require durable storage the operator will reject the deployment
without. Both are architecture decisions, not tuning.

## Application vs session — the `job` field

With `spec.job` set, this is an APPLICATION cluster: it exists to run
that one job and follows its lifecycle — the recommended production shape,
one cluster per pipeline for isolation. Without `job`, it is a SESSION
cluster: an empty runtime accepting jobs submitted separately, for many
short-lived jobs sharing warm capacity (the field doc on
[reference.md](v1alpha1/reference.md)). Choose by workload, and say which in the
proposal — they are operationally different systems.

## Stateful upgrades require durable storage — compose it

Any `job.upgradeMode` beyond "stateless" requires `state.checkpointsDir`,
and "savepoint" also requires `state.savepointsDir` — both on storage
every pod can reach, or the operator rejects the deployment. The catalog
answer is S3-compatible object storage: compose a
[KubernetesSeaweedFs](../kubernetesseaweedfs/v1alpha1/reference.md) and point
`state.s3` at it. "last-state" upgrades pair with
`state.highAvailability`, also required for the operator's rollback
feature. A stateful-upgrade Flink job proposed with no state storage is a
deployment that will not apply.

## Operator watch scope — the install-failure trap

KubernetesFlinkOperator is the registry prerequisite. Its watch-namespace
handling is unusually strict — a listed-but-absent namespace fails the
install outright — so the judgment lives in the
[operator guide](../kubernetesflinkoperator/GUIDE.md); the
invisible-edge mechanism is the
[operator-prerequisite pattern](../../_patterns/operator-prerequisite.md).

## Namespace ownership

Wire `spec.namespace` to a dedicated KubernetesNamespace rather than
`createNamespace: true`
([namespace-ownership pattern](../../_patterns/namespace-ownership.md)).

## Pairs well with

- KubernetesFlinkOperator — required; watch-scope install trap in its
  guide.
- KubernetesSeaweedFs — S3-compatible state storage for stateful
  upgrades.
- KubernetesKafka — the common source/sink for Flink streaming pipelines.
