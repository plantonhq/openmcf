# GcpGkeWorkloadIdentityBinding Guide

The judgment this guide protects: this binding is only HALF of the
Workload Identity handshake, and the half it is not — the Kubernetes
side — fails silently. The IAM grant can be perfect while pods still run
as the node's identity, and nothing errors until an API call is denied
somewhere downstream.

## The handshake has two halves

This binding grants `roles/iam.workloadIdentityUser` on the Google
service account to one KSA principal
(`serviceAccount:{project}.svc.id.goog[{namespace}/{name}]`). The other
half lives in Kubernetes: the KSA must carry the
`iam.gke.io/gcp-service-account: {gsa-email}` annotation, and the pods
must run on nodes with `workloadMetadataMode: GKE_METADATA` (the Planton
pool default posture). Miss the annotation and workloads quietly use the
node SA; miss GKE_METADATA and they see raw VM metadata — both are
misconfigurations that look like "it works" until an IAM denial says
otherwise.

## Exactly one KSA per binding, by design

The spec names one namespace/name pair — not a list — because the IAM
member string is per-principal and reviewable one grant at a time. Ten
workloads needing the same GSA is ten bindings, and that is the point:
revoking one workload's access is deleting one resource, not editing a
shared list. The role is fixed to workloadIdentityUser; broader grants
on the service account belong on `GcpServiceAccount.iam_members`, where
arbitrary-role additive IAM lives.

## On the diagram

The binding is an edge-piece: it references the `GcpServiceAccount` (by
email) and names a KSA that exists only inside a cluster — so it renders
as the visible bridge between the GCP identity graph and the Kubernetes
workloads that `GcpGkeCluster`'s workload pool hosts.

## Pairs well with

- `GcpServiceAccount` — the identity being granted; keep it minimal and
  purpose-built per workload.
- `GcpGkeCluster` — provides the workload pool
  (`{project}.svc.id.goog`) the principal grammar rides on.
- `GcpGkeNodePool` — run the pods with `workloadMetadataMode:
  GKE_METADATA` (the metadata-server half of the handshake).
