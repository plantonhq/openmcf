# KubernetesTektonOperator Guide

The judgment this guide carries: this operator is a fixed-name singleton
that installs Tekton's LIFECYCLE MANAGER and nothing else — the cluster's
actual Tekton is a separate declaration, and the two have a strict
teardown order.

## Manager only; one per cluster

The operator installs into a fixed `tekton-operator` namespace with
upstream-fixed resource names — exactly one per cluster, and it is
installed with automatic component installation DISABLED so the
[KubernetesTekton](../../kubernetestekton/v1alpha1/GUIDE.md) declaration is the
single owner of the cluster's TektonConfig. Installing this operator alone
deploys no Tekton components. The prerequisite relationship (no diagram
edge) is the
[operator-prerequisite pattern](../../../patterns/operator-prerequisite.md)'s
singleton posture.

## No version field, and destroy order matters

No version selector by design: the operator and the TektonConfig schema
the KubernetesTekton kind renders against are pinned together, so a
user-selectable version would drift the surface (the reference page
explains). And never destroy this operator while a KubernetesTekton
exists — its teardown finalizers need a running operator, or deletion
hangs. The destroy-order rule is stated from the config side too; honor it
in any teardown proposal.

## Namespace ownership — the infra exception

The fixed `tekton-operator` namespace is a dedicated sole-tenant one — the
[namespace-ownership pattern](../../../patterns/namespace-ownership.md)'s
normal infra case.

## Pairs well with

- KubernetesTekton — the one TektonConfig this operator reconciles (see
  its [guide](../../kubernetestekton/v1alpha1/GUIDE.md)).
