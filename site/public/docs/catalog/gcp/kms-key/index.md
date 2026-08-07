---
title: "KMS Key"
description: "KMS Key deployment documentation"
icon: "package"
order: 100
componentName: "gcpkmskey"
---

# GCP KMS Key

Deploys a Cloud KMS cryptographic key within an existing key ring for symmetric encryption (CMEK), asymmetric signing, asymmetric decryption, or MAC generation. The key supports automatic rotation, configurable protection levels (software or HSM), and scheduled destruction periods. Integrates with Planton's Provider Connections for GCP credential management and supports ValueFromRef wiring to key rings managed as separate Cloud Resources.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **KMS CryptoKey** -- a `kms.CryptoKey` in the specified key ring, configured with the chosen purpose, algorithm, protection level, and rotation schedule
- **Version Template** -- created only when `versionTemplate` is specified; controls the encryption algorithm and protection level (SOFTWARE, HSM, EXTERNAL, or EXTERNAL_VPC) for new key versions
- **Automatic Rotation** -- created only when `rotationPeriod` is set; generates a new primary CryptoKeyVersion at the specified interval (minimum 24 hours)
- **GCP Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking and governance

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### GCP Project

- **A KMS key ring** where the key will be created. Provide the fully qualified key ring path (`projects/{project}/locations/{location}/keyRings/{name}`) directly or reference a GcpKmsKeyRing Cloud Resource via ValueFromRef.
- **Cloud KMS API** (`cloudkms.googleapis.com`) enabled in the project that owns the key ring.

## Deploy

### Console

Open the deployment store, find **GCP KMS Key**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Symmetric Encryption** preset in the [Presets](#presets) tab to pre-populate a standard CMEK key with 90-day rotation.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpKmsKey
metadata:
  name: app-cmek-key
  org: acme-corp
  env: prod
spec:
  keyRingId:
    value: "projects/acme-prod-12345/locations/us-central1/keyRings/prod-encryption"
  keyName: app-data-cmek
  rotationPeriod: "7776000s"
```

```shell
planton apply -f gcp-kms-key.yaml
```

This creates a symmetric encryption key with 90-day automatic rotation using software-level protection. Purpose defaults to ENCRYPT_DECRYPT. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the key to a key ring deployed in the same InfraPipeline:

```yaml
spec:
  keyRingId:
    valueFrom:
      kind: GcpKmsKeyRing
      name: prod-encryption
      fieldPath: status.outputs.key_ring_id
```

The InfraPipeline resolves the dependency graph, deploys the key ring first, then provisions the KMS key with the resolved key ring path.

## Key Configuration

These are the most important decisions when configuring a KMS key. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Purpose** -- Set `purpose` to define what cryptographic operations the key supports. ENCRYPT_DECRYPT (default) is used for CMEK across BigQuery, Spanner, Cloud SQL, GCS, and GKE. ASYMMETRIC_SIGN is used for artifact and container image signing. Purpose is immutable after creation.

**Protection level** -- Set `versionTemplate.protectionLevel` to SOFTWARE (the default when unset), HSM, EXTERNAL, or EXTERNAL_VPC. HSM keys are protected by Cloud HSM hardware modules certified to FIPS 140-2 Level 3, required for PCI DSS, HIPAA, and FedRAMP compliance. The EXTERNAL levels keep the material in your own external key manager (EXTERNAL_VPC pairs with `cryptoKeyBackend`, the EkmConnection reached through your VPC). The protection level is immutable after creation.

**Bring your own key (BYOK)** -- Set `importOnly: true` (with `skipInitialVersionCreation: true` -- GCP cannot generate versions for an import-only key) when regulation requires that Google never generates the key material. Import versions after deploy with `gcloud kms keys versions import`. Start from the **Import-Only BYOK** preset.

**Labels** -- `labels` attach key-value metadata for inventory filtering and cost attribution across a key fleet; freely mutable.

**Rotation period** -- Set `rotationPeriod` for automatic key version rotation (e.g., `"7776000s"` for 90 days). Only applies to ENCRYPT_DECRYPT keys. Asymmetric keys require manual version management. Omit to disable automatic rotation.

**Destroy scheduled duration** -- `destroyScheduledDuration` controls how long destroyed key versions remain recoverable (default 30 days, minimum 24 hours). Shorter durations reduce the recovery window but limit accidental-destruction protection.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpKmsKeyRing** | `keyRingId` | `status.outputs.key_ring_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `key_id` | Fully qualified crypto key path (`projects/{p}/locations/{l}/keyRings/{kr}/cryptoKeys/{name}`) | CMEK references in BigQuery, Spanner, Cloud SQL, GCS, GKE; GcpKmsKeyIamMember `cryptoKeyId` grants |
| `key_name` | Short name of the key | Display, logging, human-readable references |
| `primary_version_name` | The current primary CryptoKeyVersion | Populated only for symmetric encryption keys; empty for asymmetric/MAC/import-only keys |
| `primary_state` | Lifecycle state of the primary version (e.g. `ENABLED`) | Health checks on the encryption path; same population rules as the version name |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Symmetric encryption (CMEK)** -- Standard customer-managed encryption key with 90-day rotation and software protection. The most common pattern for encrypting data at rest across GCP services. Start from the **Symmetric Encryption** preset.

**HSM-protected symmetric encryption** -- Same CMEK pattern but with HSM protection level for compliance scenarios requiring FIPS 140-2 Level 3 certification (PCI DSS, HIPAA, FedRAMP). Start from the **HSM Symmetric Encryption** preset.

**Asymmetric signing** -- ECDSA P-256 key for signing build artifacts, container images, JWTs, or code. No automatic rotation -- key versions are managed manually. Start from the **Asymmetric Signing** preset.

## Works With

- [**GCP KMS Key Ring**](/cloud-catalog/gcp-kms-key-ring) -- provides the key ring that contains this cryptographic key