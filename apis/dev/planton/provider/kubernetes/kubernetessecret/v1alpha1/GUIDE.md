# KubernetesSecret Guide

The judgment this guide carries: there are THREE ways a secret reaches a
workload on this platform, and picking by habit — always minting a
Secret resource — misses what the other two paths were built to give
you. Choose by where the value's source of truth lives.

## The three paths

1. **This kind** — when the platform manifest IS the source of truth:
   you hold the value (a TLS pair you were issued, a partner's API key)
   and want it as a visible, referenceable node other resources can
   point at (workloads' `imagePullSecrets` reference it by name; the
   type vocabulary is on [reference.md](reference.md)).
2. **KubernetesExternalSecret** — when the truth lives in an external
   backend (AWS Secrets Manager, GCP Secret Manager, OpenBao/Vault...).
   The cluster copy becomes a refreshing projection and the external
   system stays authoritative. If the value already lives in a secrets
   manager, this is the right path — re-declaring it here forks the
   truth into two places that will drift.
3. **The container env `secrets` block** — for app-private values
   consumed by exactly one workload. The module materializes a managed
   Secret for you; no separate node earns its place.

What breaks when chosen wrong: path 1 for an externally-managed value
means rotation in the backend never reaches the cluster; path 3 for a
shared value means the second consumer copies it, and the copies drift.

## Values you never need to paste

Two of this kind's variants exist so that a value does NOT have to be
supplied by hand — reach for them before embedding material in a
manifest: a serviceAccountToken Secret is POPULATED BY THE CLUSTER's
token controller, and issued TLS certificates should come from
cert-manager (KubernetesCertificate materializes its own Secret) rather
than being pasted into a `tls` block that will expire unnoticed.

## Namespace ownership

A Secret exists to be consumed — its namespace is shared with its
consumers by definition. Wire `spec.namespace` to the application's
KubernetesNamespace via `valueFrom`; `createNamespace: true` on a Secret
is the multi-tenant trap —
[namespace-ownership pattern](../../../patterns/namespace-ownership.md).
Note the schema quirk: namespace is OPTIONAL here and defaults to the
cluster's `default` namespace when omitted — an omitted namespace is
almost always an accident, not a choice.

## On the diagram

A declared Secret is a node consumers reference — a reviewer sees which
workloads depend on which credential. Values buried in env blocks render
as nothing, which is fine precisely when nothing else consumes them.

## Pairs well with

- KubernetesServiceAccount — the identity a serviceAccountToken Secret
  binds to (wired by reference).
- KubernetesExternalSecret — the projection path when truth lives
  outside (path 2 above).
- KubernetesCertificate — the issued-TLS path that supersedes hand-pasted
  `tls` data.
