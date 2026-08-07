# KubernetesVelero Guide

The judgment this guide carries: a backup that lives inside the cluster
it protects is not disaster recovery — WHERE backups land is the central
decision, and volume data needs its own explicit path or only the
Kubernetes objects are protected.

## Where backups land — and the in-cluster caveat

The `backupStorage.backend` decides the store and loads the matching
provider plugin (S3/S3-compatible, GCS, Azure Blob — keyless on
EKS/GKE/AKS via the identity annotations; the reference page has the
arms). An in-cluster
[KubernetesSeaweedFs](../kubernetesseaweedfs/v1alpha1/reference.md) works as
a target and is honest for dev or for protecting against accidental
deletion — but for genuine disaster recovery the store must survive the
cluster's death: an external bucket, not an in-cluster one. Say which
posture the proposal takes.

## Objects are not volumes — pick the data path

Backing up Kubernetes resources does NOT capture volume contents by
itself. Two explicit paths (reference page): CSI snapshots (needs a
snapshot-capable CSI driver + the feature flag) or file-system backup
via the node agent (kopia — works on any volume type). A Velero proposal
that skips this choice protects manifests and silently loses data.

## Once per cluster

CRDs and the node agent are cluster-scoped; one server owns the backup
records (release name fixed). Shared chart, own namespace
([namespace-ownership pattern](../../_patterns/namespace-ownership.md)
sole-tenant case).

## On the diagram

Velero renders as a shared-layer node; its store is configuration, not a
drawn edge — reviewers ask "where do backups land, and does volume data
travel?" because the diagram will not answer either.

## Pairs well with

- KubernetesSeaweedFs — an in-cluster target (with the DR caveat above).
- The workload and database kinds whose namespaces its backups protect —
  scheduled backups and store credentials live in this spec's own
  surface.
