# KubernetesPerconaMongoOperator Guide

The judgment this guide carries: this operator has TWO legitimate
placement shapes — and the upstream default (databases living beside
their operator, own-namespace watch) is the opposite of what the other
operator kinds on this platform default to in practice. Choose the
shape deliberately; do not assume the shared-chart shape silently
works.

## Two placement shapes, one choice

- **Co-located (upstream default)**: the operator sits in the SAME
  namespace as the KubernetesMongodb clusters it reconciles, watch
  block empty. Right when one team owns both operator and databases,
  and the blast radius of operator upgrades should end at their
  namespace.
- **Shared-chart, cluster-wide**: one operator in the shared-cluster
  chart with `watch.clusterWide: true` (or a fenced `watch.namespaces`
  list), serving every application environment. Right for the
  platform shape where teams declare databases without touching the
  shared-cluster layer.

Either way the coupling is invisible in the manifests (the
[operator-prerequisite pattern](../../_patterns/operator-prerequisite.md)
explains why). When a KubernetesMongodb sits NotReady with nothing
reconciling it, the watch scope is the first thing to check.

## Uninstall does not take the databases — but read before upgrading

The CRD lifecycle on [reference.md](v1alpha1/reference.md) is the operational
contract: uninstalling the release never cascade-deletes database
clusters (the upstream safety posture), and a chart upgrade runs new
operator code against EXISTING CRDs — the release notes decide when
CRDs must be applied by hand. Fold that check into any upgrade
proposal.

## Namespace ownership

In the shared-chart shape, a dedicated single-tenant namespace with
`createNamespace: true` is the normal form (the
[namespace-ownership pattern](../../_patterns/namespace-ownership.md)'s
sole-tenant case). In the co-located shape the namespace is shared with
the databases — use a dedicated KubernetesNamespace and wire everyone
through it.

## Pairs well with

- KubernetesMongodb — the clusters this operator reconciles (see its
  [guide](../kubernetesmongodb/GUIDE.md)).
