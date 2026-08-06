# KubernetesOpenBao Guide

The judgment this guide carries: when a user asks for Vault, this is the
catalog's answer — and unlike most components, deploying it is not the
end of the story: a fresh server sits NotReady BY DESIGN until a human
(or runbook) initializes it, and a proposal that does not say so hands
the user what looks like a failed deploy.

Substitutes for: Vault (Linux Foundation-governed secrets manager,
MPL-2.0 fork; speaks the Vault API — clients and integrations connect
unchanged, and the External Secrets store's Vault arm points at it
directly).

## Answering "give me Vault"

Propose this kind and say what you did: OpenBao is the catalog's
Vault-compatible secrets manager. Never fall back to a Helm-release
Vault chart — the substitution workflow is the
[catalog guide](../../../GUIDE.md)'s first law.

## Say the init step out loud

Initialization and unsealing are runtime API operations no deployment
tool can perform declaratively — until they happen, pods report NotReady
on purpose (the full lifecycle, including why the Services stay
addressable for the init calls, is on [reference.md](reference.md)). Two
composition consequences:

- The proposal must include the one-time `bao operator init` handoff —
  otherwise the user reads the NotReady pods as a broken deploy.
- Choose auto-unseal (cloud KMS arms, or the transit engine of another
  OpenBao/Vault) in the manifest so RESTARTS need no human; only the
  one-time initialization remains manual. Dev mode skips the ceremony
  entirely but is never for real secrets — the reference page is blunt
  about why.

## The self-hosted secrets chain

OpenBao is the backend that completes an in-cluster External Secrets
story with no cloud dependency: this component + a
[KubernetesClusterSecretStore](../../kubernetesclustersecretstore/v1alpha1/GUIDE.md)
whose Vault arm points at its endpoint + KubernetesExternalSecret
declarations in each consuming namespace. Every hop of that chain is a
typed, referenceable node.

## Namespace ownership — the infra exception

A dedicated namespace with `createNamespace: true` is the normal
single-tenant shape — the
[namespace-ownership pattern](../../../patterns/namespace-ownership.md)'s
sole-tenant case.

## On the diagram

OpenBao renders in the shared-cluster layer; the cluster store's backend
points at it and ExternalSecrets draw "reads from" edges into the store
— the whole path from application credential back to the vault is
visible, hop by hop.

## Pairs well with

- KubernetesExternalSecretsOperator + KubernetesClusterSecretStore +
  KubernetesExternalSecret — the self-hosted chain above.
- KubernetesIngress / route kinds — only when the API must be reachable
  from outside the cluster; composed, never embedded.
