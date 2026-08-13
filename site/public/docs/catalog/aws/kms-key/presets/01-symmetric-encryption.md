---
title: "Symmetric Encryption Key"
description: "This preset creates a customer-managed symmetric KMS key (`SYMMETRIC_DEFAULT`) with automatic rotation enabled and a 30-day deletion window. Symmetric keys are what AWS service integrations (S3, RDS,..."
type: "preset"
rank: "01"
presetSlug: "01-symmetric-encryption"
componentSlug: "kms-key"
componentTitle: "KMS Key"
provider: "aws"
icon: "package"
order: 1
---

# Symmetric Encryption Key

This preset creates a customer-managed symmetric KMS key (`SYMMETRIC_DEFAULT`) with automatic rotation enabled and a 30-day deletion window. Symmetric keys are what AWS service integrations (S3, RDS, Lambda environment encryption, MSK, ...) expect.

## When to Use

- Encrypting data at rest across AWS services (S3, EBS, RDS, DynamoDB, Secrets Manager)
- Customer-managed key requirements for compliance (HIPAA, PCI-DSS, SOC2)
- Envelope encryption where AWS generates and manages data keys on your behalf

## Key Configuration Choices

- **SYMMETRIC_DEFAULT** (`key_spec: SYMMETRIC_DEFAULT`) — AES-256-GCM; the default and required shape for most AWS integrations
- **Rotation enabled** (`enable_key_rotation: true`) — automatic key material rotation; old material retained for decryption
- **365-day rotation period** (`rotation_period_in_days: 365`) — AWS default interval; only applies when rotation is enabled
- **30-day deletion window** (`deletion_window_days: 30`) — recovery window before permanent destruction
- **Friendly alias** (`aliases`) — human-readable address (`alias/...`); many aliases may point at one key

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aws-region>` | AWS region where the key lives | Your deployment region |
| `<key-description>` | Purpose of the key (e.g. "Production database encryption") | Your team's naming conventions |
| `alias/my-symmetric-key` | Rename to your alias (e.g. `alias/myapp/data-encryption`) | Must match `alias/[0-9A-Za-z_/-]+` |

## Related Presets

- Set `key_spec` to an RSA/ECC/HMAC provider string and matching `key_usage` for signing or MAC workflows (asymmetric keys do not support automatic rotation)
