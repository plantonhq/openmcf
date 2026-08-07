# KubernetesPostgres Guide

The judgment this guide carries: what this component assumes is already on
the cluster, and who should own its namespace — the two places
agent-composed architectures with a database go wrong before a single
PostgreSQL setting matters.

## The architecture must include the operator

This component does not run PostgreSQL by itself: it renders a CloudNativePG
`Cluster` custom resource, and
**[KubernetesCloudNativePgOperator](../../kubernetescloudnativepgoperator/v1alpha1/GUIDE.md)
must be on the cluster** to reconcile it — proposing a database without it
deploys a custom resource nothing acts on, with no error anywhere (the
[operator-prerequisite pattern](../../../_patterns/operator-prerequisite.md)
is the general mechanism).

Two couplings to get right:

- **Backups couple to the operator's configuration.** Declaring
  `spec.backup` here requires the operator component installed with
  `barmanCloudPlugin.enabled` — the backup objects this spec renders are
  reconciled by that plugin. A backup declared on the database with a
  plugin-less operator is silently inert infrastructure (the operator's
  guide carries the plugin side).
- **One operator per cluster, many databases.** The operator is
  cluster-scoped; compose it once (typically in the shared-cluster chart),
  then any number of KubernetesPostgres resources in application
  environments.

## Namespace ownership

`spec.namespace` is a required foreign key targeting KubernetesNamespace.
`createNamespace: true` makes THIS database the namespace's owner in IaC
state: the namespace is created before the cluster and — per the field's
own contract — **deleted with the resource**. Safe for a namespace whose
only tenant is this database; wrong the moment a cache, an app, or a second
database shares it (the second `createNamespace` deploy fails on
already-exists, and destroying the database would delete the neighbors'
namespace). The judgment, the failure story, and the `valueFrom` wiring:
[namespace-ownership pattern](../../../_patterns/namespace-ownership.md).

## On the diagram

Wiring `spec.namespace` via `valueFrom` draws the database's namespace edge
on the architecture diagram; the operator component appears as its own node
in the shared-cluster layer. An architecture whose diagram shows database,
namespace, and operator is one a user can actually reason about — a
`createNamespace` flag and an assumed operator show nothing.

## Pairs well with

- KubernetesCloudNativePgOperator — required, once per cluster.
- KubernetesNamespace — the namespace owner (pattern above).
- Application workloads connect through the `<name>-rw` / `<name>-ro`
  Services and `<name>-app` credential Secret the reference page's naming
  contract documents.
