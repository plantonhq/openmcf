# KubernetesSecretStore Guide

The judgment this guide carries: this is the one-team grain of the
External Secrets story — its blast radius deliberately ends at the
namespace boundary, and that boundary is the reason to choose it over a
fenced cluster-wide store.

## When the extra node earns its place

The full cluster-vs-namespaced comparison lives in the
[KubernetesClusterSecretStore guide](../kubernetesclustersecretstore/GUIDE.md).
The short form: use THIS kind when one team owns both the backend and
its credentials — the store, its credential Secrets, and every consumer
live in the team's namespace, and nothing outside it can sync from the
connection. A platform-wide backend shared by many teams belongs in a
fenced ClusterSecretStore instead; duplicating per-namespace stores for
the same shared backend multiplies credential copies for no isolation
gain.

It is also the upstream default: an ExternalSecret's `storeRef.kind`
defaults to SecretStore, so same-namespace wiring needs no kind line.

## Namespace ownership

The store lives in the team's application namespace — a shared,
multi-tenant namespace by definition (consumers live there too). Wire
`spec.namespace` to the team's KubernetesNamespace via `valueFrom`;
`createNamespace: true` on a store is the multi-tenant trap —
[namespace-ownership pattern](../../_patterns/namespace-ownership.md).

## On the diagram

The store renders inside the team's namespace with "reads from" edges
from its ExternalSecrets — the team's secret topology stays visibly
self-contained, which is exactly the property this grain buys.

## Pairs well with

- KubernetesExternalSecretsOperator — required machinery, once per
  cluster.
- KubernetesExternalSecret — same-namespace consumers.
- KubernetesNamespace — the team namespace this store belongs to.
