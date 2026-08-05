# KubernetesPerconaMysqlOperator Guide

The judgment this guide carries: like its MongoDB sibling, this
operator has TWO legitimate placement shapes — co-located beside its
databases (the upstream default; watch block empty) or shared-chart
with `watch.clusterWide: true` / a fenced `watch.namespaces` list.
Choose deliberately: the default scope means a database declared in
another namespace sits unreconciled with no error anywhere.

## Choosing the shape

Co-located suits one team owning both operator and databases, keeping
operator-upgrade blast radius inside their namespace. The shared-chart
cluster-wide shape suits the platform posture — application
environments declare KubernetesMysql clusters without touching the
shared layer. The coupling is never a manifest edge; when a cluster
sits NotReady with nothing acting on it, check the watch scope first.

## Uninstall does not take the databases — but read before upgrading

The CRD lifecycle on [reference.md](reference.md): uninstalling the
release never cascade-deletes database clusters, and chart upgrades run
new operator code against EXISTING CRDs — apply the new release's CRDs
by hand when its notes call for it. Fold that check into any upgrade
proposal.

## Namespace ownership

Shared-chart shape: dedicated single-tenant namespace,
`createNamespace: true` is the normal form (the
[namespace-ownership pattern](../../../patterns/namespace-ownership.md)'s
sole-tenant case). Co-located shape: the namespace is shared with the
databases — dedicated KubernetesNamespace, everyone wired through it.

## Pairs well with

- KubernetesMysql — the clusters this operator reconciles (see its
  [guide](../../kubernetesmysql/v1/GUIDE.md)).
