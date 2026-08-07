# KubernetesFlinkOperator Guide

The judgment this guide carries: this operator's watch-namespace list is
not just a filter — it is a set of namespaces the install actively writes
into, so a name that does not exist fails the whole install. That makes
watch scope the decision to get right before anything reconciles.

## Watch namespaces must already exist

Empty `watchNamespaces` = watch every namespace, the normal
one-operator-per-cluster posture. With a list set, the chart scopes RBAC
and the admission webhook to exactly those namespaces AND plants job
ServiceAccount/Role/RoleBinding into each — but does NOT create them, so
a listed namespace that is absent fails the install with
`namespaces "…" not found` (the field doc on [reference.md](reference.md)).
When fencing the operator, the namespaces must exist first (declare them
as KubernetesNamespace resources). The general invisible-edge mechanism:
[operator-prerequisite pattern](../../../_patterns/operator-prerequisite.md).

## Once per cluster; the infra-namespace exception

One operator reconciles every Flink cluster; declare it once in the
shared-cluster chart. Its own dedicated namespace with
`createNamespace: true` is the normal sole-tenant shape
([namespace-ownership pattern](../../../_patterns/namespace-ownership.md)).

## Pairs well with

- KubernetesFlinkDeployment — the Flink clusters this operator reconciles
  (see its [guide](../../kubernetesflinkdeployment/v1alpha1/GUIDE.md)).
