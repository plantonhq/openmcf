# GCP Service Account IAM Member

Grants one role to one identity ON a GCP service account — the least-privilege unit of impersonation and account-usage control. The grant merges into the account's IAM policy without touching any other member's bindings, and removal subtracts only this exact pair, so grants from different charts, teams, and tools never fight.

## What Gets Created

When you deploy a GcpServiceAccountIamMember resource, Planton provisions:

- **Service Account IAM Member** — one (role, member[, condition]) entry merged into the target service account's IAM policy

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing service account** — referenced via `serviceAccountId`
- **IAM permissions** — `roles/iam.serviceAccountAdmin` on the target account (or its project)

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpServiceAccountIamMember
metadata:
  name: github-deployer-impersonation
spec:
  serviceAccountId:
    value: projects/my-gcp-project-123/serviceAccounts/deployer@my-gcp-project-123.iam.gserviceaccount.com
  role:
    value: roles/iam.workloadIdentityUser
  member:
    value: principalSet://iam.googleapis.com/projects/123456789/locations/global/workloadIdentityPools/github/attribute.repository/my-org/my-repo
```

```shell
planton apply -f impersonation-grant.yaml
```

Compose with first-class nodes instead of literals: reference a GcpServiceAccount's `name` output for the target account and another GcpServiceAccount's `member` output for the grantee.

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `serviceAccountId` | `StringValueOrRef` | — | Required. Fully-qualified account resource name (`projects/<project>/serviceAccounts/<email>`). References GcpServiceAccount by default. Immutable. |
| `role` | `StringValueOrRef` | — | Required. Predefined usage role (`roles/iam.workloadIdentityUser`, `roles/iam.serviceAccountTokenCreator`, `roles/iam.serviceAccountUser`) or custom role name. References GcpIamCustomRole by default. Immutable. |
| `member` | `StringValueOrRef` | — | Required. Identity in IAM member format (`serviceAccount:`, `principalSet://`, `user:`, `group:`, `domain:`). References GcpServiceAccount by default. Immutable. |
| `condition` | `object` | — | Optional IAM Condition (`title`, `expression`, optional `description`). Part of the grant's identity. Immutable. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `service_account_id` | The account whose policy received the grant (after reference resolution) |
| `role` | The granted role (after reference resolution) |
| `member` | The granted member (after reference resolution) |
| `etag` | The account IAM policy etag after the grant — useful for audit correlation |

## Related Components

- [GcpServiceAccount](/docs/catalog/gcp/gcpserviceaccount) — the account being granted on and the most common member
- [GcpWorkloadIdentityPoolProvider](/docs/catalog/gcp/gcpworkloadidentitypoolprovider) — issues the federated principals this grant authorizes
- [GcpGkeWorkloadIdentityBinding](/docs/catalog/gcp/gcpgkeworkloadidentitybinding) — the GKE-specific impersonation convenience
