# AliCloudKmsKey

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1alpha1`

AliCloudKmsKeySpec defines the configuration for an Alibaba Cloud Key
Management Service (KMS) customer-managed key (CMK).

A KMS key is used to encrypt and decrypt data across Alibaba Cloud services
such as RDS (Transparent Data Encryption), OSS (Server-Side Encryption),
ECS (disk encryption), and PolarDB. Keys can also be used for digital
signing when configured with an asymmetric key spec and SIGN/VERIFY usage.

Keys are immutable in several dimensions: key_spec, key_usage,
protection_level cannot be changed after creation. Rotation and deletion
protection can be toggled at any time.

Deletion is never immediate -- Alibaba Cloud enforces a configurable
pending window (7-366 days) during which the key can be recovered.

Provider resources:
  Terraform: alicloud_kms_key
  Pulumi:    kms.Key

## Example

```yaml
apiVersion: alicloud.planton.dev/v1alpha1
kind: AliCloudKmsKey
metadata:
  name: alicloudkmskey-demo
spec:
  region: cn-hangzhou
  description: Demo KMS key for local testing
  tags:
    purpose: demo
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.keySpec` | `string` |  | `Aliyun_AES_256` |  |
| `spec.keyUsage` | `string` |  | `ENCRYPT/DECRYPT` |  |
| `spec.protectionLevel` | `string` |  | `SOFTWARE` |  |
| `spec.automaticRotation` | `bool` |  | `false` |  |
| `spec.rotationInterval` | `string` |  |  |  |
| `spec.pendingWindowInDays` | `int32` |  | `30` |  |
| `spec.deletionProtection` | `bool` |  | `false` |  |
| `spec.deletionProtectionDescription` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region where the KMS key will be created.
The key can only be used by resources in the same region.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.description

`string`

Human-readable description of the key's intended use.

### spec.keySpec

`string` · optional (explicit presence)

Cryptographic algorithm for the key material.

Symmetric keys (ENCRYPT/DECRYPT):
  "Aliyun_AES_256" -- AES-256 (default, recommended for most use cases)
  "Aliyun_AES_128" -- AES-128 (Dedicated KMS only)
  "Aliyun_AES_192" -- AES-192 (Dedicated KMS only)
  "Aliyun_SM4"     -- Chinese national standard SM4

Asymmetric keys (SIGN/VERIFY):
  "RSA_2048"  -- RSA 2048-bit
  "RSA_3072"  -- RSA 3072-bit
  "EC_P256"   -- NIST P-256 elliptic curve
  "EC_P256K"  -- secp256k1 elliptic curve
  "EC_SM2"    -- Chinese national standard SM2

This field is immutable after creation (ForceNew in the provider).
Default: "Aliyun_AES_256"

- default: `Aliyun_AES_256`
- rule: key_spec must be one of: Aliyun_AES_256, Aliyun_AES_128, Aliyun_AES_192, Aliyun_SM4, RSA_2048, RSA_3072, EC_P256, EC_P256K, EC_SM2

### spec.keyUsage

`string` · optional (explicit presence)

Intended cryptographic operation for the key.

"ENCRYPT/DECRYPT" -- symmetric encryption/decryption (use with AES/SM4 key specs).
"SIGN/VERIFY"     -- asymmetric signing/verification (use with RSA/EC key specs).

This field is immutable after creation (ForceNew in the provider).
Default: "ENCRYPT/DECRYPT"

- default: `ENCRYPT/DECRYPT`
- rule: key_usage must be one of: ENCRYPT/DECRYPT, SIGN/VERIFY

### spec.protectionLevel

`string` · optional (explicit presence)

Protection level for the key material.

"SOFTWARE" -- key material is protected by software-based cryptographic
  modules. Suitable for most workloads.
"HSM" -- key material is stored and used exclusively within a Hardware
  Security Module. Required for regulatory compliance in some industries.

This field is immutable after creation (ForceNew in the provider).
Default: "SOFTWARE"

- default: `SOFTWARE`
- rule: protection_level must be one of: SOFTWARE, HSM

### spec.automaticRotation

`bool` · optional (explicit presence)

Enable automatic key rotation.

When enabled, Alibaba Cloud periodically generates new key material for
the CMK according to rotation_interval. Previous key versions remain
available for decryption, but new encrypt operations use the latest version.

Automatic rotation is only supported for symmetric keys (Aliyun_AES_*,
Aliyun_SM4). Asymmetric keys (RSA_*, EC_*) do not support rotation.

When set to true, rotation_interval must also be provided.
Default: false

- default: `false`

### spec.rotationInterval

`string`

Period between automatic rotations.

Accepted formats: "Nd" for days (e.g., "365d"), "Ns" for seconds
(e.g., "604800s"), "Nh" for hours (e.g., "8760h").

Required when automatic_rotation is true. Ignored when false.
Common production value: "365d" (annual rotation).

### spec.pendingWindowInDays

`int32` · optional (explicit presence)

Number of days to wait before permanently deleting the key after a
deletion request. During this window the key enters PendingDeletion
state and can be recovered by cancelling the deletion.

Range: 7-366 days.
Default: 30

- default: `30`
- rule: {"int32":{"lte":366,"gte":7}}

### spec.deletionProtection

`bool` · optional (explicit presence)

Enable deletion protection to prevent accidental key deletion.

When enabled, any attempt to delete the key will be rejected until
deletion protection is explicitly disabled. Strongly recommended for
production encryption keys -- losing a KMS key means permanent,
irrecoverable data loss for any data encrypted with it.
Default: false

- default: `false`

### spec.deletionProtectionDescription

`string`

Reason for enabling deletion protection. Provides an audit trail for
why the key is protected. Only meaningful when deletion_protection
is true.

### spec.tags

`map<string, string>`

Tags to apply to the KMS key.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudKmsKey, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.key_id` | `string` | The KMS key ID assigned by Alibaba Cloud. This is the primary identifier used by downstream components (RDS, OSS, ECS) when configuring encryption with a customer-managed key. |
| `status.outputs.arn` | `string` | The key ARN (Alibaba Cloud Resource Name). Format: acs:kms:{region}:{account-id}:key/{key-id} Used in RAM policies to grant encrypt/decrypt permissions on this key. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AliCloudKubernetesCluster | `spec.encryptionProviderKey` | `status.outputs.key_id` |

## See Also

- [Overview](../README.md)
