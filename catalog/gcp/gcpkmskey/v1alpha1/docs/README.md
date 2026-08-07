# GcpKmsKey — Deep Dive

## The problem this resource solves

Customer-managed encryption is the control plane of data trust on GCP: whoever controls the key controls access to everything encrypted under it, regardless of storage IAM. This kind models that key as a first-class node so every CMEK relationship in an environment is an explicit, reviewable reference — a BigQuery dataset, a Spanner database, a GKE cluster's etcd, a Pub/Sub topic all pointing at the same auditable key resource — instead of a string pasted into twenty specs.

## Where it sits in the composition

- **GcpKmsKeyRing** — the parent container; this key references its `status.outputs.key_ring_id` and inherits project + location.
- **GcpKmsKey** — this resource: purpose, rotation, protection level, version lifecycle.
- **CMEK consumers** — 20+ catalog kinds reference this key's `status.outputs.key_id` (the fully qualified `projects/{p}/locations/{l}/keyRings/{r}/cryptoKeys/{k}` path) from their CMEK fields.
- **Bare-name consumers** — components that take a key name plus separately supplied project/location (for example OpenBao's GCP KMS seal) compose from `key_name`.
- **IAM** — each consumer's service agent needs `roles/cloudkms.cryptoKeyEncrypterDecrypter` on the key; model those grants as first-class IAM nodes next to the consumers.

## Lifecycle contract

| Property | Behavior |
|---|---|
| `keyRingId`, `keyName`, `purpose`, `destroyScheduledDuration`, `importOnly`, `cryptoKeyBackend`, `versionTemplate.protectionLevel` | Immutable (ForceNew) — for an undeletable resource, "change" means abandon and create under a new name |
| `rotationPeriod`, `versionTemplate.algorithm`, `labels` | Mutable in place |
| Deletion | **The key object can never be deleted.** Destroy schedules every version for destruction, disables rotation, and removes the key from state; the key remains (inert, free) in the ring forever |
| Version destruction | Destroyed versions sit in `DESTROY_SCHEDULED` for the recovery window (default 30 days), then the material is gone — and with it, everything encrypted under those versions |

## The version model (what rotation actually does)

A crypto key is a container of versions. Exactly one version is *primary* (ENCRYPT_DECRYPT keys only): it encrypts all new data. Decryption automatically selects the version that encrypted the ciphertext, so **rotation never re-encrypts anything** — it caps how much data any single version protects going forward. Consequences worth designing for:

- Old versions must stay ENABLED as long as any data encrypted under them lives.
- Cost scales with live versions (each HSM version bills individually) — rotation cadence is a cost lever, not just a security one.
- Asymmetric/MAC keys have no primary and never rotate automatically: verifiers pin exact versions, so new versions are minted deliberately.

## Protection levels and the external arms

| Level | Material lives | Use |
|---|---|---|
| `SOFTWARE` (default) | Google's software fleet | Standard CMEK |
| `HSM` | Cloud HSM (FIPS 140-2 L3) | Compliance-grade CMEK/signing |
| `EXTERNAL` | An external key manager, linked per version via external key URIs | Cloud EKM over the internet |
| `EXTERNAL_VPC` | An external key manager reached over VPC; the key-level `cryptoKeyBackend` names the EKM connection | Cloud EKM, private path |

`importOnly` is orthogonal: a BYOK container GCP never generates material for (requires `skipInitialVersionCreation`, enforced pre-deploy). The import ceremony itself (key ring import jobs, 3-day wrapping keys) is operational, not durable infrastructure — deliberately outside this kind.

## Pre-deploy coherence rules

The spec enforces before any cloud call what the API would reject or silently mishandle at create time:

- `rotation_period` only on ENCRYPT_DECRYPT keys (and ≥ 86400s with ≤ 9 fractional digits — the provider's own validator, applied earlier).
- `import_only` ⇒ `skip_initial_version_creation` (an auto-generated version would violate the import-only guarantee).
- `crypto_key_backend` ⇔ `EXTERNAL_VPC` protection (EXTERNAL keys link material per version instead).
- `purpose` and `versionTemplate.algorithm` are deliberately free strings: GCP adds purposes and algorithms over time (post-quantum schemes being the current example), and a hardcoded allowlist would reject valid new values while the API accepts them. The API remains the validator of record for those two.

## 90/10 coverage vs the provider resource

| Provider field (`google_kms_crypto_key`) | Modeled | Notes |
|---|---|---|
| `key_ring` | ✅ `keyRingId` (ref → GcpKmsKeyRing) | ForceNew |
| `name` | ✅ `keyName` | ForceNew |
| `purpose` | ✅ `purpose` | free string, ForceNew |
| `rotation_period` | ✅ `rotationPeriod` | provider-validator CEL |
| `destroy_scheduled_duration` | ✅ `destroyScheduledDuration` | ForceNew |
| `version_template.algorithm` | ✅ | mutable |
| `version_template.protection_level` | ✅ (all four levels) | ForceNew |
| `skip_initial_version_creation` | ✅ | create-time only |
| `import_only` | ✅ | ForceNew |
| `crypto_key_backend` | ✅ `cryptoKeyBackend` (un-defaulted ref) | ForceNew |
| `labels` | ✅ `labels` (merged beneath attribution labels) | mutable |
| `primary` (computed) | ✅ outputs `primary_version_name` / `primary_state` | ENCRYPT_DECRYPT only |

## Deliberately not modeled (recorded reasons)

- **`deletion_policy`** — absent from the released `google ~> 6.x` GA resource (it exists on the provider's unreleased line); modeling it would create a one-engine field. The bridged Pulumi provider *does* carry it client-side, so the Pulumi module pins it explicitly to `DELETE` — the released Terraform destroy behavior — keeping destroy semantics byte-identical across engines.
- **`key_access_justifications_policy`** — beta-only (google-beta / unreleased GA); revisit when it reaches the released GA line.
- **`google_kms_crypto_key_version` as a kind** — versions are the rotation *product*, managed by the service; IaC-pinned versions matter only for EXTERNAL key URIs (Tier-2 with the EKM family).
- **`google_kms_ekm_connection`** — the external-key-manager connection is a Tier-2 kind candidate; `cryptoKeyBackend` is already ref-shaped (un-defaulted) so the future kind attaches with a one-line `default_kind`.
- **Autokey (`google_kms_autokey_config`, `google_kms_key_handle`)** — folder-scoped singletons on a different provisioning philosophy (GCP creates keys on demand); revisit on concrete pull.
- **`google_kms_secret_ciphertext`** — a one-shot encryption transform, not durable infrastructure.
- **Per-key IAM (`google_kms_crypto_key_iam_*`)** — resource-scoped IAM stays out of the catalog pending concrete pull.

## Provider mapping

Maps to `google_kms_crypto_key` (`google/services/kms/resource_kms_crypto_key.go`). The provider's Delete destroys all versions and disables rotation, then removes the key from state — the key object persists in GCP. Only `rotation_period`, `version_template.algorithm`, and `labels` are in the resource's update mask; everything else is ForceNew.
