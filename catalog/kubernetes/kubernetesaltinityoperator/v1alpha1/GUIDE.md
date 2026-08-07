# KubernetesAltinityOperator Guide

The judgment this guide carries: this operator watches ONLY its own
namespace until told otherwise, and its default operator credentials are
the publicly documented ones — two settings that decide whether an
architecture actually works and whether it is safe.

## Watch scope — the own-namespace trap

The default (empty `watchNamespaces`) confines the operator to its
install namespace: a KubernetesClickHouse anywhere else sits
unreconciled with no error. Widen it deliberately — list every namespace
that will hold clusters, or `[".*"]` for the whole cluster. The
mechanism and why this never shows on the diagram:
[operator-prerequisite pattern](../../../_patterns/operator-prerequisite.md).
The shared-chart, watch-everything shape is the usual platform posture;
`namespaceScopedRbac` is only sound in the own-namespace-alone shape (the
field doc on [reference.md](reference.md) explains why).

## Set operator credentials outside dev

`operatorCredentials` unset means the operator connects to every managed
ClickHouse with the publicly documented default password (reference page
is blunt about this). For anything beyond a throwaway cluster, set a real
password — the operator's account exists on every cluster it manages, so
a default here is a cluster-wide exposure, not a per-instance one.

## Once per cluster, shared-cluster chart; the infra-namespace exception

One operator reconciles every ClickHouse cluster; declare it once in the
shared-cluster chart. A dedicated single-tenant namespace with
`createNamespace: true` is the normal form here — the
[namespace-ownership pattern](../../../_patterns/namespace-ownership.md)'s
sole-tenant case.

## Pairs well with

- KubernetesClickHouse — the clusters this operator reconciles (see its
  [guide](../../kubernetesclickhouse/v1alpha1/GUIDE.md)).
