# KubernetesKubeRayOperator Guide

The judgment this guide carries: one operator per cluster, watching
everything by default — the standard shape. The decision it anchors sits
on the CR side (GCS fault tolerance), so this guide stays short by design.

## Once per cluster, watches cluster-wide

The normal posture is a single install in the shared-cluster chart
watching every namespace; fence it only when isolation demands it. The
invisible-edge mechanism and the watch-posture table:
[operator-prerequisite pattern](../../_patterns/operator-prerequisite.md).
Its own dedicated namespace with `createNamespace: true` is the
sole-tenant case of the
[namespace-ownership pattern](../../_patterns/namespace-ownership.md).

## The judgment lives on the cluster

The high-stakes Ray decision — surviving head-pod loss via GCS fault
tolerance — is a KubernetesRayCluster field, not an operator setting; see
the [RayCluster guide](../kubernetesraycluster/GUIDE.md). Installing
this operator alone runs no Ray.

## Pairs well with

- KubernetesRayCluster — the Ray clusters this operator reconciles (see
  its [guide](../kubernetesraycluster/GUIDE.md)).
