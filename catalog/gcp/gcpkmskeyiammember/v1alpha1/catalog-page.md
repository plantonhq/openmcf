# GCP KMS Key IAM Member

Grants one role to one identity ON a KMS crypto key — the least-privilege unit of customer-managed encryption (CMEK) access control. The grant merges into the key's IAM policy without touching any other member's bindings, and removal subtracts only this exact pair, so grants from different charts, teams, and tools never fight.

## What Gets Created

When you deploy a GcpKmsKeyIamMember resource, Planton provisions:

- **Crypto Key IAM Member** — one (role, member[, condition]) entry merged into the target key's IAM policy

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing crypto key** — referenced via `cryptoKeyId`
- **IAM permissions** — `roles/cloudkms.admin` on the target key (or its ring/project)

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpKmsKeyIamMember
metadata:
  name: gcs-state-key-grant
spec:
  cryptoKeyId:
    value: projects/my-gcp-project-123/locations/us-central1/keyRings/app-ring/cryptoKeys/state-key
  role:
    value: roles/cloudkms.cryptoKeyEncrypterDecrypter
  member:
    value: serviceAccount:service-123456789@gs-project-accounts.iam.gserviceaccount.com
```

```shell
planton apply -f cmek-grant.yaml
```

Compose with first-class nodes instead of literals: reference a GcpKmsKey's `key_id` output for the key, so the encrypted resource can be ordered after the grant it needs.

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `cryptoKeyId` | `StringValueOrRef` | — | Required. Fully-qualified key path (`projects/<project>/locations/<location>/keyRings/<ring>/cryptoKeys/<key>`). References GcpKmsKey by default. Immutable. |
| `role` | `StringValueOrRef` | — | Required. Predefined role (typically `roles/cloudkms.cryptoKeyEncrypterDecrypter`) or custom role name. References GcpIamCustomRole by default. Immutable. |
| `member` | `StringValueOrRef` | — | Required. Identity in IAM member format (`serviceAccount:`, `user:`, `group:`, `domain:`, `principal://`). References GcpServiceAccount by default. Immutable. |
| `condition` | `object` | — | Optional IAM Condition (`title`, `expression`, optional `description`). Part of the grant's identity. Immutable. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `crypto_key_id` | The key whose policy received the grant (after reference resolution) |
| `role` | The granted role (after reference resolution) |
| `member` | The granted member (after reference resolution) |
| `etag` | The key IAM policy etag after the grant — useful for audit correlation |

## Related Components

- [GcpKmsKey](/docs/catalog/gcp/gcpkmskey) — the key being granted on
- [GcpKmsKeyRing](/docs/catalog/gcp/gcpkmskeyring) — the key's parent ring
- [GcpServiceAccount](/docs/catalog/gcp/gcpserviceaccount) — a grantable identity
