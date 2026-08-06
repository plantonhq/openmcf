# GCP Workload Identity Pool

Creates a Workload Identity Pool — the trust boundary for keyless authentication into Google Cloud. External identities (GitHub Actions, GitLab CI, AWS workloads, enterprise IdPs) federate through the pool instead of holding long-lived service-account keys.

## What Gets Created

When you deploy a GcpWorkloadIdentityPool resource, Planton provisions:

- **Workload Identity Pool** — a `google_iam_workload_identity_pool` in the target project, addressable as `projects/<number>/locations/global/workloadIdentityPools/<poolId>`

The pool holds no issuer configuration; attach one GcpWorkloadIdentityPoolProvider per external issuer.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing GCP project** — referenced via `projectId`
- **IAM permissions** — `roles/iam.workloadIdentityPoolAdmin` on the target project

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpWorkloadIdentityPool
metadata:
  name: github-actions-pool
spec:
  projectId:
    value: my-gcp-project-123
  workloadIdentityPoolId: github-actions
  displayName: GitHub Actions
  description: Keyless federation for the engineering org's CI pipelines
```

```shell
planton apply -f pool.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `workloadIdentityPoolId` | `string` | — | Required. Pool ID (4-32 chars; lowercase letters, digits, hyphens; `gcp-` prefix reserved). Immutable. |
| `projectId` | `StringValueOrRef` | provider default | Project that owns the pool. Can reference a GcpProject. Immutable. |
| `displayName` | `string` | `""` | Console-visible name (max 32 chars). Mutable. |
| `description` | `string` | `""` | What this pool federates and who owns it (max 256 chars). Mutable. |
| `disabled` | `bool` | `false` | Kill switch — a disabled pool rejects all token exchanges. Mutable. |
| `mode` | `string` | `FEDERATION_ONLY` | `FEDERATION_ONLY` or `TRUST_DOMAIN` (managed workload identities). Immutable. |
| `inlineCertificateIssuanceConfig` | `object` | — | mTLS certificate issuance for TRUST_DOMAIN pools (region→CA-pool map, key algorithm, lifetime, rotation). |
| `inlineTrustConfig` | `object` | — | PEM trust anchors for additional trusted trust domains. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `name` | Full resource name — the handle IAM principals (`principal://iam.googleapis.com/<name>/subject/...`) and providers are built from |
| `workload_identity_pool_id` | The bare pool ID; feed directly into a GcpWorkloadIdentityPoolProvider's `workloadIdentityPoolId` field |
| `state` | `ACTIVE`, or `DELETED` while soft-deleted (~30 days; the ID cannot be reused until permanent deletion) |

## Related Components

- [GcpWorkloadIdentityPoolProvider](/docs/catalog/gcp/gcpworkloadidentitypoolprovider) — attaches an external issuer to this pool
- [GcpServiceAccount](/docs/catalog/gcp/gcpserviceaccount) — the identity federated principals typically impersonate
- [GcpProjectIamMember](/docs/catalog/gcp/gcpprojectiammember) — grants roles to principals built from this pool
- [GcpProject](/docs/catalog/gcp/gcpproject) — provides the project that owns the pool
