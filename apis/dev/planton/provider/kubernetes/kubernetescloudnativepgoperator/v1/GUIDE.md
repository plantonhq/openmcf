# KubernetesCloudNativePgOperator Guide

The judgment this guide carries: one CloudNativePG operator serves every
PostgreSQL cluster on the cluster — and the coupling that actually bites
is not installation but CONFIGURATION: a database's backup declaration is
inert unless THIS operator was installed with the matching plugin.

## Once per cluster, in the shared-cluster chart

Install it once; application environments declare KubernetesPostgres
clusters, never their own operator (the
[Postgres guide](../../kubernetespostgres/v1/GUIDE.md) approaches the
same coupling from the database side). The dependency draws no diagram
edge — the
[operator-prerequisite pattern](../../../patterns/operator-prerequisite.md)
is the mechanism. Its "cnpg-system" namespace is the sole-tenant case of
the [namespace-ownership pattern](../../../patterns/namespace-ownership.md).

## The backup-plugin coupling

A KubernetesPostgres declaring `spec.backup` requires this operator
installed with `barmanCloudPlugin.enabled` — the backup objects the
database renders are reconciled by that plugin, and with a plugin-less
operator they are silently inert infrastructure (the Postgres guide
names this from its side; the plugin toggle lives HERE). When any
database in the architecture declares backups, verify this operator's
plugin arm in the same review.

## On the diagram

The operator renders in the shared-cluster layer; databases render in
their environments with no edge to it. Reviewers verify the operator
node exists — and, when backups are declared anywhere, that its plugin
configuration matches.

## Pairs well with

- KubernetesPostgres — the clusters this operator reconciles (see its
  [guide](../../kubernetespostgres/v1/GUIDE.md)).
