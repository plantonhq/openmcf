# Minimal Installation

This preset installs the External Secrets Operator with chart defaults: one
controller replica, CRDs installed with the release and kept on uninstall,
webhook and cert-controller enabled. No ambient cloud identity is
configured — every store authenticates through its own auth block
(per-store ServiceAccounts or credential Secrets). This is the right
starting point for most clusters.

## When to Use

- First ESO installation on a cluster, before any stores exist
- Clusters where stores will carry their own identities or static
  credentials (including token-based backends like Vault, and any
  self-managed/kind cluster)
- When you prefer per-store identity isolation over one ambient identity

## Key Configuration Choices

- **Chart defaults everywhere** — one replica of each component, default
  chart version pin, `info` logging
- **CRDs install and survive uninstall** (spec defaults) — uninstalling the
  release never cascade-deletes ExternalSecret/SecretStore objects
- **No `workloadIdentity`** — stores must bring their own auth; add ambient
  identity later without reinstalling (see preset 02)
- **Owned namespace** (`createNamespace: true`) — the module creates and
  labels the `external-secrets` namespace

## Placeholders to Replace

No placeholders — this preset is directly deployable.

## Related Presets

- **02-eks-ambient-identity** — bind one IAM role as the fallback identity
  for every store on an EKS cluster
- **03-tuned-multi-team** — sizing and concurrency for clusters with many
  teams and hundreds of ExternalSecrets
