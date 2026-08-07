---
title: "KMS Key"
description: "KMS Key deployment documentation"
icon: "package"
order: 100
componentName: "ocikmskey"
---

# KMS Key on OCI

Deploys an Oracle Cloud Infrastructure Key Management Service encryption key inside a KMS vault. Supports AES, RSA, and ECDSA algorithms with configurable key length, HSM/software/external protection modes, and optional automatic key rotation. The key is tied to a specific vault via the `managementEndpoint` field. Downstream OCI services (Block Volume, Object Storage, Database, File Storage) consume the key OCID for customer-managed encryption. The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to both compartments and vaults.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **KMS Key** -- a `kms.Key` in the specified compartment and vault with configurable algorithm, key length, protection mode, and optional auto-rotation schedule. The key is created in ENABLED state and an initial key version is generated automatically.
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the key

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the key in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- A KMS vault with its management endpoint. The `managementEndpoint` ties this key to a specific vault. Provide the endpoint URL directly or reference an OciKmsVault Cloud Resource via ValueFromRef using `status.outputs.managementEndpoint`.
- For external protection mode only: an external key ID on the third-party key manager. The vault itself must also be of type `external`.

## Deploy

### Console

Open the deployment store, find **KMS Key on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **AES-256 HSM Auto-Rotation** preset in the [Presets](#presets) tab to pre-populate a standard data-at-rest encryption key with 90-day rotation.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciKmsKey
metadata:
  name: data-encryption-key
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  managementEndpoint:
    value: "https://xxx-management.kms.us-ashburn-1.oraclecloud.com"
  keyShape:
    algorithm: aes
    length: 32
```

```shell
planton apply -f key.yaml
```

This creates a 256-bit AES key with HSM protection (the default). Auto-rotation is not enabled, and no display name is configured. The key OCID and current key version OCID are exported as stack outputs.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the key to both a compartment and a vault deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: security-compartment
      fieldPath: status.outputs.compartmentId
  managementEndpoint:
    valueFrom:
      kind: OciKmsVault
      name: platform-vault
      fieldPath: status.outputs.managementEndpoint
```

The InfraPipeline resolves the dependency graph, deploys the compartment and vault first, then provisions the key with the resolved values.

## Key Configuration

These are the most important decisions when configuring a KMS key. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Key shape** -- The `keyShape` message defines the algorithm and length. All fields are immutable after creation. AES keys (symmetric) are used for data-at-rest encryption: 16 bytes (128-bit), 24 bytes (192-bit), or 32 bytes (256-bit). RSA keys (asymmetric) are used for signing and verification: 256 bytes (2048-bit), 384 bytes (3072-bit), or 512 bytes (4096-bit). ECDSA keys require a matching `curveId`: `nist_p256` with length 32, `nist_p384` with length 48, or `nist_p521` with length 66.

**Protection mode** -- The `protectionMode` field is immutable after creation. `hsm` (default) stores key material in a FIPS 140-2 Level 3 hardware security module -- the highest isolation level. `software` stores keys in software at lower cost. `external` references key material on a third-party key manager and requires an `externalKeyReference` with the external key ID.

**Auto-rotation** -- Set `isAutoRotationEnabled: true` and configure `autoKeyRotationDetails.rotationIntervalInDays` to enable automatic key version rotation. OCI creates a new key version on schedule; existing data encrypted with older versions remains readable. Auto-rotation is appropriate for symmetric (AES) keys. For asymmetric keys (RSA, ECDSA), rotation requires coordinated consumer updates to distribute the new public key.

**Vault binding** -- The `managementEndpoint` field ties this key to a specific vault. It must match the vault's management endpoint output. Use ValueFromRef to reference an OciKmsVault resource, ensuring the key is always created in the correct vault.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |
| **OciKmsVault** | `managementEndpoint` | `status.outputs.managementEndpoint` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `key_id` | OCID of the KMS key | OciBlockVolume, OciObjectStorageBucket, OciAutonomousDatabase, OciVaultSecret -- any resource supporting customer-managed encryption |
| `current_key_version` | OCID of the currently active key version | Key version auditing, rotation tracking |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**AES-256 HSM with auto-rotation** -- A 256-bit AES symmetric key in an HSM with 90-day automatic rotation. The standard key for encrypting data at rest across OCI services. Meets SOC 2, ISO 27001, HIPAA, and PCI-DSS requirements for HSM-backed encryption with periodic rotation. Start from the **AES-256 HSM Auto-Rotation** preset.

**RSA-4096 HSM signing** -- A 4096-bit RSA asymmetric key in an HSM for digital signing and verification. Auto-rotation is disabled because asymmetric key rotation requires coordinated consumer updates. Used for container image signing, artifact verification, and API request signing. Start from the **RSA-4096 HSM Signing** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this key
- [**KMS Vault on OCI**](/cloud-catalog/oci-kms-vault) -- provides the vault and management endpoint where this key is created