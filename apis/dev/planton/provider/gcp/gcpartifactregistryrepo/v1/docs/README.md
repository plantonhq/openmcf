# GcpArtifactRegistryRepo — Deep Dive

## The problem this resource solves

Every workload on GCP consumes artifacts — container images into GKE/Cloud Run/Cloud Functions, language packages into builds, OS packages into VMs — and the repository that stores them is load-bearing infrastructure: if it is unavailable, deploys stop; if it grows unbounded, cost grows with it; if its access is wrong, either builds break or artifacts leak. This kind models the repository as a first-class node so the artifact supply chain is explicit and reviewable: which repositories exist, what they cache, who can pull and push, what gets cleaned up, and which resources consume them by reference.

## Where it sits in the composition

- **GcpProject** — the owning project (`projectId` reference, or ambient).
- **GcpArtifactRegistryRepo** — this resource, in any of its three serving modes. Repositories compose with *each other*: a virtual repository's upstream policies and a remote repository's custom upstream both reference another repository's `repository_path` output.
- **GcpCloudFunction** — its `dockerRepository` takes this repository's `repository_path` for user-managed build-artifact storage.
- **GcpKmsKey** — `kmsKeyName` for CMEK-protected artifact storage.
- **GcpServiceAccount** — `iamMembers[].member` resolves a service account's `member` output for pull/push grants.
- **GKE / Cloud Run image references** — workloads pull `{registry_uri}/{image}:{tag}`; the registry URI output is exactly that prefix.

## Lifecycle contract

| Property | Behavior |
|---|---|
| `format`, `mode`, `location`, `projectId`, `kmsKeyName`, `repositoryId` | Immutable (ForceNew) — replacement deletes every stored artifact |
| `remoteRepositoryConfig` (upstream arms) | Immutable — a cache's upstream is its identity |
| `remoteRepositoryConfig.upstreamCredentials`, `disableUpstreamValidation` | Mutable — credential rotation never recreates the cache |
| `description`, `labels`, `dockerConfig`, `cleanupPolicies`, `cleanupPolicyDryRun`, `virtualRepositoryConfig`, `vulnerabilityScanningEnablement`, `iamMembers` | Mutable in place |
| Deletion | Deletes the repository and every artifact version in it — there is no force-destroy gate; protect precious repositories with IAM and KEEP policies |

## Cleanup policies (the cost model)

A busy CI pipeline pushes thousands of images a month, and every superseded digest keeps billing. Cleanup policies are the declarative answer:

- **DELETE** policies match versions by age (`olderThan`/`newerThan`), tag state (`UNTAGGED` digests are the classic target), and package/tag/version prefixes.
- **KEEP** policies protect versions — either by condition or by `mostRecentVersions.keepCount` ("always keep the last N per package"). On overlap, KEEP always wins.
- `cleanupPolicyDryRun` runs the pipeline in log-only mode — validate new policies in Cloud Audit Logs before letting them delete.

The standard pairing: one DELETE for old untagged versions, one KEEP for the N most recent — busy repositories stay bounded without ever deleting something a rollback needs.

## Remote repositories (the hermetic-build arm)

A remote repository caches exactly one upstream: the well-known public registries (`DOCKER_HUB`, `MAVEN_CENTRAL`, `NPMJS`, `PYPI`), the OS mirror trees (Apt/Yum bases + path), or — via `commonRepository.uri` — any custom registry or another Artifact Registry repository. Authenticated upstreams take a username plus a **Secret Manager secret version path**: the password itself never enters the spec, the state, or the module — only the reference to where it lives. The Artifact Registry service agent reads the secret at pull time (`roles/secretmanager.secretAccessor`).

## Virtual repositories (the one-endpoint arm)

Virtual repositories serve from multiple upstreams by priority — first-party repos above shared caches — and are pull-only: CI keeps pushing to standard repositories, consumers pull one URL forever while the serving graph changes underneath by editing upstream policies (mutable in place).

## Access model

Additive IAM only: each `iamMembers` entry grants one role to one member on this repository and composes safely with grants made anywhere else — removal subtracts exactly that grant. Public access is the same mechanism (`allUsers` + `roles/artifactregistry.reader`), never a separate flag, so the ACL surface reads uniformly. Authoritative bindings/policies (which clobber grants they do not list) are deliberately not modeled.

## Recorded scope decisions

- **Repository-level IAM trio kinds** (`_iam_binding`/`_iam_policy` as standalone kinds) — not modeled: authoritative semantics are hostile to composition; the additive `iamMembers` list covers the per-repository grant surface.
- **`google_artifact_registry_rule`** and **`google_artifact_registry_project_config`** — deferred: project-scoped/rule surfaces outside the 90/10 line for repository delivery.
- **`google_artifact_registry_vpcsc_config`** — beta-only; not modeled on the released GA line.
- **`deletion_policy`** — absent from the released 6.x Terraform resource; the bridged Pulumi provider's client-side flag is pinned to `DELETE` for byte-identical destroy behavior (PARITY note in both modules).
- **`registry_uri` as a provider attribute** — absent from the released 6.x resource; both modules construct it identically from resolved attributes.
