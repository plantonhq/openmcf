# KubernetesSecretStore: Research and Design

## Introduction

A SecretStore is the External Secrets Operator (ESO) custom resource that
represents one backend connection serving ExternalSecrets in a SINGLE
namespace. This component creates one SecretStore, named after the resource
(`metadata.name`), plus the credential Secret its configuration declares.

## The Upstream Model

ESO splits secret syncing into two custom resources:

- **SecretStore / ClusterSecretStore** — the connection: which backend, and
  how the controller authenticates to it. Upstream gives the two kinds an
  IDENTICAL spec; they differ only in scope, plus the cluster kind's
  `conditions` (which namespaces may use it). A namespaced SecretStore is
  usable only by ExternalSecrets in its own namespace.
- **ExternalSecret** — one sync declaration: which entries to read from a
  store and which Kubernetes Secret to materialize them into.

The Planton kinds mirror this split exactly: KubernetesSecretStore and
KubernetesClusterSecretStore share one `ExternalSecretsStoreConfig` proto
and one Pulumi spec builder, so the twins are structurally incapable of
drifting; scope and `conditions` are the only differences. The operator
installation itself is a third kind (KubernetesExternalSecretsOperator).

## The Grain: One Store Per Namespace/Team

The namespaced grain is the blast-radius choice. The store lives in
`spec.namespace`; its credential Secrets materialize there; only that
namespace's ExternalSecrets can sync through it. RBAC on the namespace is
the whole access story — no cluster-scoped fencing needed. The alternative
grain, KubernetesClusterSecretStore, serves every namespace and fences with
`conditions`; choose it for platform-wide backends every team shares.

## Backend Coverage (and the deliberate exclusions)

| Backend | Coverage |
|---|---|
| AWS | Secrets Manager / SSM Parameter Store / ACM export (`service`); region; assume-role ARN (the cross-account pattern); key `prefix` scoping; IRSA ServiceAccount auth, ambient identity, or static access keys |
| GCP Secret Manager | Project (FK to GcpProject); regional endpoint (`location`); Workload Identity ServiceAccount auth, ambient identity, or a static service-account key |
| Azure Key Vault | Vault URL; tenant; `auth_type` WorkloadIdentity / ManagedIdentity / ServicePrincipal; ServiceAccount federation or client-id/client-secret |
| Vault / OpenBao | Server URL, KV mount + version (v1/v2), enterprise namespace, private CA bundle; token / AppRole / Kubernetes auth methods. OpenBao speaks the Vault API — the same arm serves both |
| Kubernetes | Another cluster's Secrets as the backend: remote API server + CA, remote namespace, bearer-token or local-ServiceAccount auth — the cluster-to-cluster sync arm |
| Fake | Literal key/value/version entries served by ESO itself — test and sandbox only |

**Other upstream backends are deliberately unmodeled**: ESO supports dozens
of providers (IBM Secrets Manager, Oracle Vault, Akeyless, Doppler,
1Password, and more) whose demand does not justify first-class modeling.
They remain unmodeled until demand shows; a store for any of them can be
applied via KubernetesManifest without leaving the platform.

**Generators and PushSecret are unmodeled**: ESO's generator resources
(dynamic per-sync values like ECR tokens) and PushSecret (writing cluster
Secrets BACK to a backend) are separate machinery from the read-and-sync
path this family models.

## Authentication Model

Uniform across the cloud backends, three postures:

1. **Keyless via a referenced ServiceAccount** — `service_account_name` is a
   foreign key to KubernetesServiceAccount, whose workload-identity arms
   (IRSA / GKE Workload Identity / AKS Workload Identity) carry the cloud
   binding. ESO exchanges that ServiceAccount's token with the cloud;
   different stores can carry different identities — the multi-team
   posture. On a SecretStore the ServiceAccount defaults to the store's own
   namespace; an explicit `service_account_namespace` is unusual.
2. **Keyless via the operator's ambient identity** — leave auth empty; the
   controller's own ServiceAccount identity
   (KubernetesExternalSecretsOperator's `workload_identity`) or node
   identity is used.
3. **Declared static credentials** — sensitive spec values are materialized
   as a deterministically named Kubernetes Secret
   (`<resource-name>-credentials`, fixed data keys per backend) in the
   store's namespace; the CR references that Secret. Exactly-one contracts
   reject mixing postures, and paired-field contracts reject
   half-credentials.

## Credential Model Mechanics

A namespaced store reads its referenced Secrets from its OWN namespace —
the shared builder omits explicit namespaces from secretRefs when rendering
for a SecretStore (the CRD defaults them). The modules create credential
Secrets BEFORE the CR and make the CR depend on them, so ESO never observes
a store whose secretRefs dangle.

## Engine Mechanics

- **Pulumi**: the shared `externalsecretsstore` builder renders the
  CRD-JSON spec map; the CR applies as an untyped CustomResource. ESO's
  validating webhook checks the applied spec strictly, and the kind-cluster
  E2E lanes exercise the sync machinery live — shape errors fail loudly.
- **Terraform**: `kubectl_manifest` (alekc/kubectl) applies the CR — no
  cluster connection at plan time, so a store can be PLANNED before the
  operator's CRDs exist (single-run infra charts, offline plan proofs).
- **Neither engine waits for Ready**: readiness depends on external
  reachability (the cloud secrets API, Vault) that is not part of applying
  the resource — the never-block-on-a-controller posture.

## The Fake Backend's Role

The `fake` backend serves the literal entries declared in the spec — no
external account, no network, fully deterministic. It exists for pipelines,
tests, and evaluating the sync machinery: the kind-cluster E2E lane drives
a fake-backend store through a real ExternalSecret sync and asserts the
materialized Secret's contents. Never use it for real secrets: the values
sit in plain text in the store spec.

## Validation Highlights

Exactly-one contracts (backend arm, Vault auth method, one auth posture per
backend), paired-field contracts (AWS key pair, Azure service-principal
pair), vocabulary checks (AWS `service`, Azure `auth_type`, Vault KV
`version`), and Go-duration patterns on intervals — each with a teaching
message naming the fix.

## E2E

The fake backend reaches Ready with zero external dependencies — the
deterministic kind-cluster lane, verified end-to-end through an
ExternalSecret sync. Cloud-backend lanes need real accounts and identity
bindings and ride the batched real-cluster lanes.
