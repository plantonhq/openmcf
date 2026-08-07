# KubernetesRayCluster Guide

The judgment this guide carries: a Ray cluster keeps its entire control
state in the head pod's memory, so losing the head loses everything unless
you compose fault tolerance — and the Ray version and image must move in
lockstep. Neither fails at apply time.

## Head loss wipes the cluster — enable GCS fault tolerance

The head pod's Global Control Store holds all job, actor, and worker state
IN MEMORY. Without fault tolerance, deleting or losing the head rebuilds an
empty cluster — every job and registration gone. Enable
`gcsFaultTolerance` backed by an external Redis-protocol store for
head-loss survival: compose a
[KubernetesValkey](../kubernetesvalkey/v1alpha1/reference.md) and point it
there (the field doc on [reference.md](v1alpha1/reference.md)). Worker pods are
always expendable; the head is the single point of failure this setting
removes. Any long-lived Ray cluster proposal should include it.

## Version and image move together

`rayVersion` must match the Ray inside `image` — the operator reads the
version to shape commands but runs the image as given, so a mismatch fails
at runtime, not apply. The spec derives image from `rayVersion` by
default; only an explicit `image` override can break the lockstep, so
override both together or neither.

## Operator and namespace

KubernetesKubeRayOperator is the registry prerequisite and watches
cluster-wide by default
([operator-prerequisite pattern](../../_patterns/operator-prerequisite.md);
watch judgment in its [guide](../kuberneteskuberayoperator/GUIDE.md)).
Wire `spec.namespace` to a dedicated KubernetesNamespace rather than
`createNamespace: true`
([namespace-ownership pattern](../../_patterns/namespace-ownership.md)).

## Pairs well with

- KubernetesKubeRayOperator — required (see its
  [guide](../kuberneteskuberayoperator/GUIDE.md)).
- KubernetesValkey — the Redis-protocol store backing GCS fault
  tolerance.
