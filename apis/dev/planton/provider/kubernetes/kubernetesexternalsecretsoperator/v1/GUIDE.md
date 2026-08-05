# KubernetesExternalSecretsOperator Guide

The judgment this guide carries: like cert-manager, this operator is
machinery that does nothing visible by itself — an architecture that
declares external secrets without it fails silently, with ExternalSecret
resources nothing reconciles. It is also strictly a singleton: the
release name is fixed, so a second installation is not a supported shape.

## Install once, in the shared-cluster chart

One installation serves the whole cluster: every store and every
ExternalSecret depends on this single controller machinery (the
three-component split and the fixed-release-name constraint are on
[reference.md](reference.md)). It belongs in the shared-cluster chart.
Application environments declare stores and ExternalSecrets — never
their own operator.

## The chain this component roots

1. **This operator** — once, shared chart.
2. **A store** — the backend connection:
   [KubernetesClusterSecretStore](../../kubernetesclustersecretstore/v1/GUIDE.md)
   for platform-wide backends, KubernetesSecretStore for one team's.
3. **KubernetesExternalSecret** — one per secret to sync, in the
   consuming application's namespace.

Store credential Secrets for cluster-scoped stores conventionally live in
this operator's install namespace — the ClusterSecretStore spec's
`secretsNamespace` defaults its reference to this component's namespace
output; wire it with `valueFrom` so the dependency is explicit.

## Namespace ownership — the infra exception

A dedicated "external-secrets" namespace with `createNamespace: true` is
the normal single-tenant shape — the
[namespace-ownership pattern](../../../patterns/namespace-ownership.md)'s
sole-tenant case.

## On the diagram

The operator renders in the shared-cluster layer; each ClusterSecretStore
draws its `secretsNamespace` reference edge into it. Like every
operator-backed capability, ExternalSecrets do not draw an edge to the
operator itself — reviewers verify the operator node exists.

## Pairs well with

- KubernetesClusterSecretStore / KubernetesSecretStore — the backend
  connections this operator serves.
- KubernetesExternalSecret — the per-secret sync declarations.
- KubernetesOpenBao — the self-hosted backend option: a complete
  in-cluster secrets chain with no cloud dependency.
