# KubernetesOpenSearchOperator Guide

The judgment this guide carries: this operator watches ALL namespaces by
default — the friendly posture — and its fence, unusually, is SINGULAR:
`watchNamespace` restricts it to exactly one namespace, not a list.

## Watch posture: all namespaces, or exactly one

The default (cluster-wide watch) is the platform shape: one operator in
the shared chart, OpenSearch clusters declared anywhere. The fence is
different from every other operator in the catalog: `watchNamespace`
takes ONE namespace (pair it with `useRoleBindings` to drop cluster-wide
RBAC — the reference page carries both). Fencing to one namespace means
one operator per team namespace — a heavier shape; choose it only when
RBAC isolation genuinely demands it. The invisible-edge mechanism:
[operator-prerequisite pattern](../../../patterns/operator-prerequisite.md).

## The judgment lives on the cluster

The high-stakes OpenSearch decisions — the demo-credentials trap, the
security config — are KubernetesOpenSearch fields; see its
[guide](../../kubernetesopensearch/v1/GUIDE.md). Installing this operator
alone runs no OpenSearch.

## Namespace ownership — the infra exception

A dedicated single-tenant namespace with `createNamespace: true` is the
normal shape — the
[namespace-ownership pattern](../../../patterns/namespace-ownership.md)'s
sole-tenant case.

## Pairs well with

- KubernetesOpenSearch — the clusters this operator reconciles (see its
  [guide](../../kubernetesopensearch/v1/GUIDE.md)).
