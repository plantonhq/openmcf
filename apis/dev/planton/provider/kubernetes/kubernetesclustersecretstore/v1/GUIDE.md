# KubernetesClusterSecretStore Guide

The judgment this guide carries: a cluster-scoped store with no
`conditions` is an offer of its backend to EVERY namespace on the
cluster. Declare the fence with the store, on day one — the multi-tenancy
mistake here is invisible until the wrong namespace syncs a production
credential.

## Cluster grain vs namespaced grain (the comparison home)

- **This kind** — one backend connection shared platform-wide: the
  organization's AWS Secrets Manager, the platform OpenBao. Fence it
  with `conditions` (exact names, label selectors, or regexes — the
  vocabulary is on [reference.md](reference.md)) so a store holding
  production credentials is not readable from every dev namespace.
- **KubernetesSecretStore** — the same backend surface at namespace
  grain: only ExternalSecrets in its own namespace may use it, and its
  credential Secrets live there too. Reach for it when ONE team owns the
  backend credentials and the blast radius should end at their
  namespace boundary.

Note the reference default: an ExternalSecret's `storeRef.kind` defaults
to SecretStore — syncing from a cluster store requires saying
`kind: ClusterSecretStore` explicitly.

## Wiring conventions

- `secretsNamespace` is a foreign key defaulting to the
  KubernetesExternalSecretsOperator resource's namespace output — wire it
  with `valueFrom` so the store's dependency on the operator's home is
  explicit and ordered.
- The backend connection (including the Vault/OpenBao arm — point
  `server` at an in-cluster KubernetesOpenBao for the self-hosted chain)
  is the same `config` surface as the namespaced kind; the grain is the
  only difference.

## On the diagram

The store renders as a shared-layer node; every consuming ExternalSecret
draws a labeled "reads from" edge into it — the cluster's secret
topology becomes reviewable. A fenced store also tells a reviewer WHO may
consume it; an unfenced one cannot.

## Pairs well with

- KubernetesExternalSecretsOperator — required machinery; wired via
  `secretsNamespace`.
- KubernetesExternalSecret — the consumers (each says
  `kind: ClusterSecretStore` explicitly).
- KubernetesOpenBao — the self-hosted backend behind the store's
  Vault/OpenBao arm.
