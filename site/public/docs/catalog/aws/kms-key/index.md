---
title: "KMS Key"
description: "KMS Key deployment documentation"
icon: "package"
order: 100
componentName: "awskmskey"
---

# AWS KMS Key

Deploys a customer-managed KMS key with configurable cryptographic shape (key spec and usage), a custom key policy, automatic key rotation with a configurable period, multi-Region designation, a scheduled deletion window, and a list of aliases. The key integrates with Planton's Provider Connections for AWS credential management and provides `key_arn` and `key_id` outputs that downstream Cloud Resources consume via ValueFromRef for envelope encryption.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **KMS Key** -- a customer-managed key with the specified key spec and usage, description, key policy, rotation configuration, multi-Region designation, and deletion window
- **KMS Aliases** -- one alias resource per entry in `aliases`; friendly names (e.g., `alias/data-encryption`) that reference the key, added and removed in place as the list changes
- **AWS Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the key

## Before You Deploy

### Planton Setup

- **AWS Provider Connection** -- an active connection in the Connect module with credentials for the target AWS account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or cross-account trust authentication modes.

### AWS Account

- **KMS permissions** -- the credentials used by the Provider Connection must have `kms:CreateKey`, `kms:CreateAlias`, `kms:EnableKeyRotation`, `kms:PutKeyPolicy`, and related KMS permissions.
- **Key policy** -- leaving `policy` empty keeps the AWS default policy, which grants the account root full access and enables IAM-policy delegation. Provide a custom policy document for cross-account grants or restricted administration -- and always keep an administrator principal in it, or the key can become unmanageable.
- **Region selection** -- KMS keys are regional resources. Create keys in the same region as the resources that will use them (S3 buckets, RDS instances, EKS clusters). For cross-region decryption with the same material, enable `multiRegion` and create replica keys from this primary.

## Deploy

### Console

Open the deployment store, find **AWS KMS Key**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Symmetric Encryption** preset in the [Presets](#presets) tab to pre-populate a standard encryption key configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsKmsKey
metadata:
  name: data-encryption-key
  org: acme-corp
  env: prod
spec:
  region: us-west-2
  description: "Customer-managed key for data encryption"
  keySpec: SYMMETRIC_DEFAULT
  enableKeyRotation: true
  deletionWindowDays: 30
  aliases:
    - "alias/acme-data-encryption"
```

```shell
planton apply -f kms-key.yaml
```

This creates a symmetric encryption key with automatic rotation enabled, a 30-day deletion window, and one alias. SYMMETRIC_DEFAULT is suitable for use with S3 SSE-KMS, EBS volume encryption, RDS storage encryption, and EKS secrets encryption. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a KMS key. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Key spec and usage** -- The `keySpec` field defaults to `SYMMETRIC_DEFAULT`, the shape every AWS service integration requires. Choose an RSA (`RSA_2048`/`RSA_3072`/`RSA_4096`) or ECC (`ECC_NIST_P256`/`ECC_NIST_P384`/`ECC_NIST_P521`/`ECC_SECG_P256K1`) spec for signing or public-key workflows, an HMAC spec (`HMAC_224`/`HMAC_256`/`HMAC_384`/`HMAC_512`) for token authentication, or `SM2` in China regions. `keyUsage` is coupled to the spec's family: symmetric keys encrypt/decrypt, HMAC keys generate/verify MACs, ECC keys sign/verify, and RSA/SM2 keys choose between encryption and signing. Both fields are create-time immutable -- changing them replaces the key.

**Key rotation** -- Set `enableKeyRotation: true` to rotate the material automatically (symmetric keys only). Rotation is transparent to callers: the key ID, ARN, and aliases never change, and old ciphertext keeps decrypting. `rotationPeriodInDays` (90-2560) tunes the cadence; leaving it unset keeps AWS's 365-day default.

**Key policy** -- The `policy` field takes a JSON policy document that is the root of access control on the key. Empty keeps the AWS default policy (account root full access, IAM delegation enabled) -- the right choice for most keys. `bypassPolicyLockoutSafetyCheck` skips AWS's protection against lockout policies; leave it false unless deliberately constructing a lockout.

**Multi-Region** -- `multiRegion: true` creates a multi-Region PRIMARY key whose replicas in other regions share its material, so ciphertext encrypted in one region decrypts in another. Create-time immutable.

**Deletion window** -- The `deletionWindowDays` field (7-30 days, default 30) defines the waiting period before a scheduled key deletion becomes permanent. During this window, the deletion can be cancelled. Use 30 days for production keys to allow time for discovery of dependent resources.

**Aliases** -- Each entry in `aliases` must start with `alias/` (the `alias/aws/` prefix is reserved for AWS-managed keys). Aliases are how humans and SDK callers address the key without its generated ID; re-pointing an alias to a new key is the manual-rotation idiom for specs AWS cannot rotate automatically.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `key_id` | Generated ID of the KMS key (`mrk-…` for multi-Region keys) | Direct key references in AWS service configurations |
| `key_arn` | Amazon Resource Name of the KMS key | EKS secrets encryption, S3 SSE-KMS, RDS storage encryption, EBS volume encryption |
| `alias_names` | Alias names attached to the key, in spec order | Human-readable key references in application configuration |

The `key_arn` output is the primary value consumed by downstream resources. S3 buckets, EKS clusters, RDS instances, and EBS volumes reference it to enable customer-managed encryption instead of AWS-managed default keys.

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Symmetric encryption key** -- A SYMMETRIC_DEFAULT key with rotation enabled and a 30-day deletion window. The standard pattern for encrypting data at rest across AWS services (S3, EBS, RDS, EKS secrets, DynamoDB). Start from the **Symmetric Encryption** preset.

## Works With

This component operates independently and does not reference other components.
