# AliCloud KMS Key

Deploys an Alibaba Cloud Key Management Service (KMS) customer-managed key (CMK) for encrypting data across Alibaba Cloud services. KMS keys are used by RDS (Transparent Data Encryption), OSS (Server-Side Encryption), ECS (disk encryption), PolarDB, and NAS. Keys can also be configured for asymmetric signing when using RSA or EC key specs.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **KMS Key** -- an `alicloud_kms_key` resource with configurable key spec, usage, protection level, rotation policy, and deletion protection

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Alibaba Cloud Account

- **Immutable decisions** -- `keySpec`, `keyUsage`, and `protectionLevel` cannot be changed after creation. Choose carefully based on your encryption and compliance requirements.
- **Deletion risk** -- losing a KMS key means permanent, irrecoverable data loss for any data encrypted with it. Enable `deletionProtection` for production keys.

## Deploy

### Console

Open the deployment store, find **AliCloud KMS Key**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields including key spec, rotation, and deletion protection.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudKmsKey
metadata:
  name: platform-encryption-key
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  description: Platform encryption key for RDS and OSS
  automaticRotation: true
  rotationInterval: "365d"
  deletionProtection: true
  deletionProtectionDescription: Protects production data encryption key
  tags:
    team: platform
    compliance: soc2
```

```shell
planton apply -f alicloud-kms-key.yaml
```

This creates an AES-256 key with annual rotation and deletion protection. A Stack Job tracks the provisioning in real time.

### InfraChart

KMS keys are standalone resources with no upstream dependencies. Downstream components reference the key via ValueFromRef:

```yaml
spec:
  encryption:
    kmsMasterKeyId:
      valueFrom:
        kind: AliCloudKmsKey
        name: platform-encryption-key
        fieldPath: status.outputs.key_id
```

The InfraPipeline resolves the dependency graph and creates the KMS key before any dependent resources.

## Key Configuration

These are the most important decisions when configuring a KMS key. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Key spec** -- The `keySpec` field selects the cryptographic algorithm. "Aliyun_AES_256" (default) for symmetric encryption. "RSA_2048" or "EC_P256" for asymmetric signing. Immutable after creation.

**Key usage** -- The `keyUsage` field must match the key spec. "ENCRYPT/DECRYPT" for symmetric keys (AES/SM4). "SIGN/VERIFY" for asymmetric keys (RSA/EC). Immutable after creation.

**Automatic rotation** -- The `automaticRotation` field enables periodic key material rotation. Previous versions remain available for decryption. Only supported for symmetric keys. Set `rotationInterval` (e.g., "365d") when enabled.

**Deletion protection** -- The `deletionProtection` field prevents accidental deletion. Strongly recommended for production keys.

**Pending deletion window** -- The `pendingWindowInDays` field (7-366 days, default 30) controls how long a deleted key can be recovered before permanent destruction.

## Outputs and Dependencies

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `key_id` | KMS key ID assigned by Alibaba Cloud | AliCloudRdsInstance TDE, AliCloudStorageBucket SSE, AliCloudNasFileSystem encryption |
| `arn` | Key ARN (acs:kms:{region}:{account}:key/{id}) | RAM policies for encrypt/decrypt permissions |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard** -- A default AES-256 encryption key with no rotation. Suitable for development and non-critical workloads. Start from the **Standard** preset.

**Production with rotation** -- A key with annual automatic rotation, deletion protection, and compliance tags. Start from the **Production With Rotation** preset.

**Asymmetric signing** -- An RSA-2048 key configured for SIGN/VERIFY operations. Start from the **Asymmetric Signing** preset.

## Works With

- [**AliCloud RDS Instance**](/cloud-catalog/ali-cloud-rds-instance) -- Transparent Data Encryption with customer-managed key
- [**AliCloud Storage Bucket**](/cloud-catalog/ali-cloud-storage-bucket) -- KMS-based server-side encryption
- [**AliCloud NAS File System**](/cloud-catalog/ali-cloud-nas-file-system) -- customer-managed encryption for NAS
- [**AliCloud PolarDB Cluster**](/cloud-catalog/ali-cloud-polardb-cluster) -- TDE encryption
