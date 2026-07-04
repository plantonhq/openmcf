---
title: "KMS Key"
description: "KMS Key deployment documentation"
icon: "package"
order: 100
componentName: "awskmskey"
---

# AWS KMS Key

Deploys a customer-managed AWS KMS key: cryptographic shape (`key_spec`, `key_usage`), optional key policy, rotation, multi-region designation, and friendly aliases. KMS has no name in AWS — identity is the generated key ID/ARN; consumers compose through `status.outputs.key_arn` or an alias.

## What Gets Created

When you deploy an AwsKmsKey resource, Planton provisions:

- **KMS key** — an `aws_kms_key` with the specified `key_spec`, `key_usage`, policy, rotation settings, and `disabled` state
- **KMS aliases** — one `aws_kms_alias` per entry in `aliases` (each `alias/...` string), materialized in spec order

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **IAM permissions** to create and manage KMS keys (`kms:CreateKey`, `kms:CreateAlias`, `kms:EnableKeyRotation`, `kms:TagResource`)

## Quick Start

Create a file `kms-key.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsKmsKey
metadata:
  name: my-key
spec:
  region: us-west-2
  aliases:
    - alias/my-app
```

Deploy:

```shell
planton apply -f kms-key.yaml
```

This creates a symmetric encrypt/decrypt key with AWS default settings (no automatic rotation until you opt in).

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region where the key lives. | Required; non-empty |

All other fields are optional.

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `description` | `string` | — | Console description (max 8192 characters). |
| `key_spec` | `string` | `SYMMETRIC_DEFAULT` | Provider string: `SYMMETRIC_DEFAULT`, `RSA_2048`, `RSA_3072`, `RSA_4096`, `ECC_NIST_P256`, `HMAC_256`, etc. Create-time immutable. |
| `key_usage` | `string` | `ENCRYPT_DECRYPT` | `ENCRYPT_DECRYPT`, `SIGN_VERIFY`, or `GENERATE_VERIFY_MAC`. Create-time immutable. |
| `policy` | `string` | AWS default | Key policy as a JSON document. Empty grants account root full access and enables IAM delegation. |
| `bypass_policy_lockout_safety_check` | `bool` | `false` | Skip lockout safety check when setting a custom policy — use with extreme care. |
| `disabled` | `bool` | `false` | Create or pause the key in disabled state (all crypto operations fail). |
| `enable_key_rotation` | `bool` | `false` | Automatic rotation (SYMMETRIC_DEFAULT only). |
| `rotation_period_in_days` | `int32` | 365 | Rotation interval 90–2560; requires `enable_key_rotation`. |
| `multi_region` | `bool` | `false` | Create as multi-Region primary key (create-time immutable). |
| `deletion_window_days` | `int32` | 30 | Pending deletion window 7–30 days. |
| `aliases` | `string[]` | — | Friendly names each starting with `alias/` (e.g. `alias/orders-db`). |

## Examples

### Symmetric Key with Rotation

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsKmsKey
metadata:
  name: app-encryption
spec:
  region: us-west-2
  description: Application data encryption
  key_spec: SYMMETRIC_DEFAULT
  enable_key_rotation: true
  aliases:
    - alias/app-encryption
```

### RSA Signing Key

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsKmsKey
metadata:
  name: jwt-signing
spec:
  region: us-west-2
  key_spec: RSA_4096
  key_usage: SIGN_VERIFY
  description: JWT signing key
  aliases:
    - alias/jwt-signing
  deletion_window_days: 14
```

### Custom Key Policy

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsKmsKey
metadata:
  name: cross-account
spec:
  region: us-west-2
  policy: |
    {
      "Version": "2012-10-17",
      "Statement": [{
        "Sid": "EnableRoot",
        "Effect": "Allow",
        "Principal": {"AWS": "arn:aws:iam::123456789012:root"},
        "Action": "kms:*",
        "Resource": "*"
      }]
    }
  aliases:
    - alias/shared-encryption
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `key_id` | `string` | Generated key ID (UUID; `mrk-...` for multi-Region keys) |
| `key_arn` | `string` | Key ARN — join key for encryption-at-rest fields across the catalog |
| `alias_names` | `string[]` | Attached alias names (`alias/...`) in spec order |

## Related Components

- [AwsLambda](/docs/catalog/aws/lambda) — `kms_key_arn` encrypts environment variables
- [AwsS3Bucket](/docs/catalog/aws/s3-bucket) — server-side encryption with customer-managed keys
- [AwsRdsCluster](/docs/catalog/aws/rds-cluster) — storage encryption
- [AwsLambdaEventSourceMapping](/docs/catalog/aws/lambda-event-source-mapping) — encrypts filter criteria at rest
