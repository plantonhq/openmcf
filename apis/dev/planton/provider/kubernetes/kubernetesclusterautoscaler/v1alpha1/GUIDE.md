# KubernetesClusterAutoscaler Guide

The judgment this guide carries: this component earns its place only on
clusters whose capacity is organized as pre-defined node groups — and
deploying it on GKE or AKS, where a managed autoscaler already exists, is
a real and common mistake.

## When this, when Karpenter, when neither

- **This kind** scales EXISTING node groups up and down (EC2 ASGs, Azure
  VMSS, Cluster API MachineDeployments) — right for EKS-with-ASGs,
  Cluster API, and self-managed clusters.
- **[KubernetesKarpenter](../../kuberneteskarpenter/v1alpha1/GUIDE.md)** is the
  modern AWS alternative — right-sized nodes on demand, no pre-made
  groups. The node-scaling comparison home is its guide.
- **Neither on GKE/AKS**: those ship a MANAGED cluster autoscaler as a
  toggle on the node pool itself (on Planton, through the cluster kinds).
  Deploying this component there fights the managed one — the exception,
  never the default (the reference page states it). Confirm the cluster's
  shape before proposing this kind.

## Once per cluster

The autoscaler leader-elects and owns the cluster-wide scale decision;
a second install fights the first (release name fixed). One install in
the shared-cluster chart, its own namespace the
[namespace-ownership pattern](../../../patterns/namespace-ownership.md)'s
sole-tenant case.

## On the diagram

Renders as a shared-layer node acting on node groups the cluster kind
defines — it has no per-workload edges. Its presence (or wrongful
presence on a managed-autoscaler cluster) is a reviewer check, not a
drawn relationship.

## Pairs well with

- KubernetesKarpenter — the on-demand alternative (comparison in its
  guide).
- The cluster kinds — which define the node groups this scales (and, on
  GKE/AKS, already carry the managed autoscaler).
