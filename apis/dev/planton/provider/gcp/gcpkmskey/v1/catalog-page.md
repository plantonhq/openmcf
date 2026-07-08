# GCP KMS Key

Creates a Cloud KMS cryptographic key within an existing key ring — the resource every CMEK-capable GCP service (BigQuery, Spanner, Cloud SQL, GKE, Pub/Sub, and more) references for encryption with keys you control. Also covers asymmetric signing, raw encryption, MAC, HSM protection, and BYOK/external-key-manager arms.

## What Gets Created

- The Cloud KMS API is enabled on the ring's project (never disabled on destroy)
- A `google_kms_crypto_key` resource in the referenced key ring, carrying your labels merged beneath Planton's attribution labels (`planton-ai_resource`, `planton-ai_name`, `planton-ai_kind`, plus org/env/id when set)

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing key ring** — referenced via `keyRingId` (see `GcpKmsKeyRing`); the key inherits its project and location
- **IAM permissions** — `roles/cloudkms.admin` to create keys; consumers additionally need `roles/cloudkms.cryptoKeyEncrypterDecrypter`

## Quick Start

Create a file `kms-key.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpKmsKey
metadata:
  name: cmek-data-key
spec:
  keyRingId:
    valueFrom:
      kind: GcpKmsKeyRing
      name: prod-encryption
      fieldPath: status.outputs.key_ring_id
  keyName: cmek-data-key
  rotationPeriod: "7776000s"
```

Deploy:

```shell
planton apply -f kms-key.yaml
```

This creates a symmetric encryption key with automatic 90-day rotation — the most common CMEK configuration.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `keyRingId` | `StringValueOrRef` | Fully qualified key ring path (`projects/{p}/locations/{l}/keyRings/{name}`). Reference a `GcpKmsKeyRing` via `valueFrom`. Immutable. | Required |
| `keyName` | `string` | Name of the key in GCP. Immutable — and never reusable within the ring, because keys cannot be deleted. | 1-63 chars; `^[a-zA-Z0-9_-]{1,63}$` |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `purpose` | `string` | `ENCRYPT_DECRYPT` | `ENCRYPT_DECRYPT`, `ASYMMETRIC_SIGN`, `ASYMMETRIC_DECRYPT`, `RAW_ENCRYPT_DECRYPT`, `MAC`, or any newer API purpose (free string — the API validates). Immutable. |
| `rotationPeriod` | `string` | — | Auto-rotation cadence, ≥ `86400s` with ≤ 9 fractional digits (e.g. `"7776000s"` for 90 days). ENCRYPT_DECRYPT keys only (enforced pre-deploy). Mutable. |
| `destroyScheduledDuration` | `string` | 30 days | Recovery window destroyed versions spend in `DESTROY_SCHEDULED` before the material is gone. Immutable. |
| `versionTemplate.algorithm` | `string` | `GOOGLE_SYMMETRIC_ENCRYPTION` | Algorithm for new versions (free string — see the algorithm reference). Mutable, affects future versions only. |
| `versionTemplate.protectionLevel` | `string` | `SOFTWARE` | `SOFTWARE`, `HSM` (FIPS 140-2 L3), `EXTERNAL`, `EXTERNAL_VPC`. Immutable. |
| `skipInitialVersionCreation` | `bool` | `false` | Create the key with no versions (required for import-only keys). Create-time only. |
| `importOnly` | `bool` | `false` | BYOK container: only imported versions, ever. Requires `skipInitialVersionCreation` (enforced pre-deploy). Immutable. |
| `cryptoKeyBackend` | `StringValueOrRef` | — | EKM connection path backing `EXTERNAL_VPC` keys (enforced pre-deploy). Immutable. |
| `labels` | `map` | — | User labels for cost attribution; merged beneath the platform labels. Mutable. |

## Permanence and Destroy Semantics

**Keys cannot be deleted from GCP.** Destroying this resource destroys all key *versions* (data encrypted under them becomes permanently unrecoverable once the recovery window elapses), disables rotation, and removes the key from IaC state — the key object itself remains, inert and free, in the ring forever. Plan key names as permanent; use versioned names (`cmek-data-v2`) where replacement is conceivable.

## Stack Outputs

| Output | Description |
|--------|-------------|
| `key_id` | Fully qualified key path (`projects/{p}/locations/{l}/keyRings/{r}/cryptoKeys/{name}`) — the CMEK reference every downstream consumer takes |
| `key_name` | Short name of the key |
| `primary_version_name` | Resource name of the current primary version (ENCRYPT_DECRYPT keys; empty otherwise) |
| `primary_state` | Lifecycle state of the primary version (e.g. `ENABLED`) |

## Related Components

- [GcpKmsKeyRing](../kms-key-ring/) — the parent container (required)
- [GcpBigQueryDataset](../bigquery-dataset/), [GcpSpannerDatabase](../spanner-database/), [GcpCloudSql](../cloud-sql/), [GcpGkeCluster](../gke-cluster/), [GcpPubSubTopic](../pubsub-topic/) — CMEK consumers referencing `key_id`
- [GcpProjectIamMember](../project-iam-member/) — models encrypter/decrypter grants as first-class nodes
