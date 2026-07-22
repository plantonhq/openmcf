# KubernetesExternalSecret: Research and Design

## Introduction

An ExternalSecret is the External Secrets Operator (ESO) custom resource
that declares one sync: read entries from a store's backend, materialize
them as a Kubernetes Secret, keep the Secret fresh. This component creates
one ExternalSecret, named after the resource (`metadata.name`). The
materialized Secret is created by the OPERATOR, not by the IaC modules —
the modules apply the declaration; the controller does the syncing.

## The Upstream Model

ESO splits secret syncing into two custom resources: the store
(SecretStore / ClusterSecretStore — WHERE secrets live and HOW to
authenticate; identical specs upstream, scope and namespace conditions
being the only differences) and the ExternalSecret (WHAT to sync and into
WHICH Secret). The Planton kinds mirror the split: store connections are
KubernetesSecretStore / KubernetesClusterSecretStore, the operator
installation is KubernetesExternalSecretsOperator, and this kind carries
one sync declaration.

Because the declaration carries no connection details, it is
backend-agnostic: the same ExternalSecret syncs from AWS Secrets Manager,
GCP Secret Manager, Azure Key Vault, or Vault/OpenBao depending only on
which store it references — and the cluster it runs on (EKS, GKE, AKS,
self-managed) is independent of the backend the store points at.

## data vs dataFrom vs template

Upstream gives three ways to shape the materialized Secret, and the spec
models all three:

- **`data`** — explicit entries: one backend key (optionally one `property`
  within a structured JSON entry, one pinned `version`, one decoding
  strategy) maps to one Secret key. Precise and reviewable; the form to
  prefer for application credentials.
- **`dataFrom`** — bulk pulls, two sources per entry (exactly one set):
  `extract` pulls ALL properties of one structured backend entry (each
  property becomes a Secret key); `find` pulls EVERY backend entry matching
  a name regex and/or tags (each matched entry becomes a Secret key).
  Ordered `rewrite` steps (regex source → target with capture groups)
  reshape the resulting keys — the idiom is stripping path prefixes.
- **`target.template`** — applied when materializing: sets the Secret
  `type` (e.g. `kubernetes.io/dockerconfigjson`, `kubernetes.io/tls`),
  stamps labels/annotations, and renders `data` values as Go templates over
  the synced keys (e.g. reshape a username/password pair into a connection
  string). `merge_policy` decides whether templated keys REPLACE the synced
  keys (upstream default) or merge over them.

`data` and `dataFrom` combine freely; at least one of the two is required
(a sync that syncs nothing is an always-misconfiguration, rejected at
validation).

## Creation and Deletion Policies

The `target` carries the materialized Secret's lifecycle:

- **`creation_policy`** — Owner (default: ESO creates and owns via
  ownerReference, failing on an unowned name collision), Orphan (create but
  never garbage-collect), Merge (write synced keys into an EXISTING Secret
  ESO does not own), None (never write — niche, paired with external
  tooling).
- **`deletion_policy`** — Retain (upstream default, the safe posture: the
  Secret outlives the ExternalSecret), Delete (the Secret dies with the
  ExternalSecret and keys pruned from the backend are pruned from the
  Secret), Merge (remove only the keys this ExternalSecret owns).
- **`immutable`** — marks the materialized Secret immutable: consumers get
  kubelet-cache stability, but any refresh that would change data FAILS
  (immutable Secrets cannot be updated). Pair with `refresh_policy:
  CreatedOnce`.

## The Refresh Model

`refresh_interval` (Go duration, upstream default `1h`) is how often ESO
re-reads the backend; `0s` means fetch exactly once. `refresh_policy`
gates WHEN syncing happens at all: CreatedOnce (create, never touch again —
drift stays), Periodic (the standard posture), OnChange (only when the
ExternalSecret's spec changes). Rotation in the backend propagates within
one refresh interval — workloads consuming the Secret via volume mounts see
updated values without restart; env-var consumers need a pod restart, which
is reload tooling's job, deliberately outside this kind.

## Deliberately Unmodeled

Stated timelessly, for the record:

- **Generators** — ESO's generator resources (per-sync dynamic values:
  ECR/GCR tokens, UUIDs, passwords) are separate machinery from
  backend-read syncing; unmodeled until demand shows.
- **PushSecret** — writing cluster Secrets BACK to a backend is the reverse
  data flow; unmodeled.
- **`template.templateFrom`** — composing templates from ConfigMaps/Secrets
  is unmodeled; the inline `template.data` map covers the real cases
  (dockerconfigjson, connection strings, TLS reshaping).
- Any of these is expressible via KubernetesManifest without leaving the
  platform.

## Engine Mechanics

- **Pulumi**: the typed spec renders to the CRD-JSON shape
  (`spec_builder.go`) and applies as an untyped CustomResource — the same
  posture as the store kinds and the cert-manager family. ESO's validating
  webhook checks the applied spec strictly, and the kind-cluster E2E lanes
  verify the full sync loop live (store Ready, ExternalSecret synced,
  Secret contents), so shape errors fail loudly.
- **Terraform**: `kubectl_manifest` (alekc/kubectl) applies the CR — no
  cluster connection at plan time, so the declaration can be PLANNED before
  the operator's CRDs exist (single-run infra charts, offline plan proofs).
  The locals render the same shape as the Pulumi builder; twins kept in
  lockstep.
- **Neither engine waits for the sync**: the materialized Secret appears
  when the operator reaches the backend, which is not part of applying the
  resource. The E2E verifier (not the module) asserts synced state.

## The secret_name Output

The modules pin the CR's `target.name` to the resolved Secret name
(`target.name` when set, else `metadata.name`) so the exported handle can
never drift from what the operator creates. `status.outputs.secret_name`
is the composition hook: workloads wire env `valueFrom` and volume
`secretName` references to it, exactly as they would to a KubernetesSecret
or a KubernetesCertificate's TLS Secret.

## Validation Highlights

At-least-one contract (`data` or `data_from`), exactly-one contract per
`data_from` entry (extract vs find), find-criteria contract (regexp and/or
tags — path alone only scopes), vocabulary checks (store kind, refresh
policy, creation/deletion policies, merge policy, decoding strategy), and
Go-duration patterns on intervals — each with a teaching message naming the
fix.

## E2E

The kind-cluster lane syncs from a fake-backend store — fully
deterministic, no external account — and asserts the materialized Secret's
contents, covering the whole family (operator, store, sync) in one loop.
Cloud-backend lanes ride the batched real-cluster lanes.
