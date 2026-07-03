# GCP Workload Identity Pool Provider

Attaches one external issuer (GitHub Actions, an AWS account, a SAML IdP, or an X.509 CA) to a Workload Identity Pool — the piece that makes keyless authentication into Google Cloud actually work. No service-account keys are created, stored, or rotated.

## What Gets Created

When you deploy a GcpWorkloadIdentityPoolProvider resource, Planton provisions:

- **Workload Identity Pool Provider** — a `google_iam_workload_identity_pool_provider` inside the referenced pool, addressable as `projects/<number>/locations/global/workloadIdentityPools/<poolId>/providers/<providerId>`

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing Workload Identity Pool** — referenced via `workloadIdentityPoolId`
- **IAM permissions** — `roles/iam.workloadIdentityPoolAdmin` on the target project

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpWorkloadIdentityPoolProvider
metadata:
  name: github-oidc
spec:
  workloadIdentityPoolId:
    valueFrom:
      kind: GcpWorkloadIdentityPool
      name: github-actions-pool
      fieldPath: status.outputs.workload_identity_pool_id
  workloadIdentityPoolProviderId: github-oidc
  attributeMapping:
    google.subject: assertion.sub
    attribute.repository: assertion.repository
  attributeCondition: assertion.repository_owner == "my-org"
  oidc:
    issuerUri: https://token.actions.githubusercontent.com
```

```shell
planton apply -f provider.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `workloadIdentityPoolId` | `StringValueOrRef` | — | Required. The pool this provider belongs to (bare pool ID; reference a GcpWorkloadIdentityPool). Immutable. |
| `workloadIdentityPoolProviderId` | `string` | — | Required. Provider ID (4-32 chars; lowercase letters, digits, hyphens; `gcp-` prefix reserved). Immutable. |
| `projectId` | `StringValueOrRef` | provider default | Project that owns the pool. Immutable. |
| `displayName` | `string` | `""` | Console-visible name (max 32 chars). Mutable. |
| `description` | `string` | `""` | Which issuer this trusts (max 256 chars). Mutable. |
| `disabled` | `bool` | `false` | Kill switch — rejects new token exchanges. Mutable. |
| `attributeMapping` | `map(string)` | issuer defaults | Claim → Google-attribute CEL mappings. Required for OIDC (must include `google.subject`). |
| `attributeCondition` | `string` | accept all | CEL gate restricting which otherwise valid credentials are accepted — always set one for multi-tenant issuers. |
| `aws` \| `oidc` \| `saml` \| `x509` | `object` | — | Required, exactly one. The external issuer; the type cannot change on a live provider. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `name` | Full provider resource name — the audience string tokens are minted for and web-identity provider configurations consume |
| `workload_identity_pool_provider_id` | The bare provider ID |
| `state` | `ACTIVE`, or `DELETED` while soft-deleted (~30 days; the ID cannot be reused until permanent deletion) |

## Related Components

- [GcpWorkloadIdentityPool](/docs/catalog/gcp/gcpworkloadidentitypool) — the pool this provider attaches to
- [GcpServiceAccount](/docs/catalog/gcp/gcpserviceaccount) — grant `roles/iam.workloadIdentityUser` to this provider's principals for impersonation
- [GcpProjectIamMember](/docs/catalog/gcp/gcpprojectiammember) — grants roles directly to principals mapped by this provider
