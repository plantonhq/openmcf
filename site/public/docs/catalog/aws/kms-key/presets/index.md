---
title: "Presets"
description: "Ready-to-deploy configuration presets for KMS Key"
type: "preset-list"
componentSlug: "kms-key"
componentTitle: "KMS Key"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-symmetric-encryption"
    rank: "01"
    title: "Symmetric Encryption Key"
    excerpt: "This preset creates a customer-managed symmetric KMS key (`SYMMETRIC_DEFAULT`) with automatic rotation enabled and a 30-day deletion window. Symmetric keys are what AWS service integrations (S3, RDS,..."
  - slug: "02-cross-account-grants"
    rank: "02"
    title: "Shared Key with Grants"
    excerpt: "This preset creates a rotating symmetric key whose access is delegated through KMS grants instead of key-policy edits: one grant wires a workload role (referenced from an `AwsIamRole` in the same..."
  - slug: "03-external-key-store"
    rank: "03"
    title: "External Key Store Key"
    excerpt: "This preset creates a KMS key whose cryptographic material lives OUTSIDE standard AWS KMS: in an external key manager reached through an external key store (XKS), or -- by dropping `xks_key_id` -- in..."
---

# KMS Key Presets

Ready-to-deploy configuration presets for KMS Key. Each preset is a complete manifest you can copy, customize, and deploy.
