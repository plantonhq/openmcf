# KubernetesMysql Guide

The judgment this guide carries: applications connect to this cluster
through its PROXY service, never a database pod — and proposals sized
below three nodes are trading away the exact property (quorum-certified
writes) that makes this kind production-grade.

## Wire applications to the proxy handle

The service chart on [reference.md](reference.md) is the contract:
writes go to `<name>-haproxy` (or `<name>-proxysql`), which routes to a
healthy node and re-routes on failure without client changes. An
application pointed at a pod or the internal `<name>-pxc` service
bypasses failure handling — it works until the first node event, which
is the worst possible time to discover it. Credentials come from the
operator-generated `<name>-secrets` Secret — wire by reference, never
copy.

## Three nodes is the floor, on purpose

Galera is synchronous multi-primary: a committed transaction exists on
every node, so losing one loses nothing — but that guarantee needs a
quorum, and the operator rejects sizes below 3 unless the `unsafe`
block explicitly opts into the single-node development posture. A
proposal for anything user-facing starts at 3.

## The operator must watch THIS namespace

KubernetesPerconaMysqlOperator is the registry prerequisite; its watch
scope and the two placement shapes (co-located vs shared-chart
cluster-wide) are the
[operator guide](../../kubernetesperconamysqloperator/v1alpha1/GUIDE.md)'s
judgment. The coupling is invisible in the manifests
([operator-prerequisite pattern](../../../patterns/operator-prerequisite.md)) —
check the watch scope.

## Namespace ownership

The cluster shares its namespace with consumers or a co-located
operator — the multi-tenant case where `createNamespace: true` is
wrong; wire `spec.namespace` to a dedicated KubernetesNamespace —
[namespace-ownership pattern](../../../patterns/namespace-ownership.md).

## On the diagram

Cluster and operator render as separate nodes; applications relate to
the cluster through the shared namespace and the credentials Secret.
The proxy layer lives inside this spec — routing choices change no
node count.

## Pairs well with

- KubernetesPerconaMysqlOperator — required; placement judgment in its
  guide.
- KubernetesNamespace — the shared namespace's owner (pattern above).
- Application workloads — proxy service + generated credentials, both
  by reference.
