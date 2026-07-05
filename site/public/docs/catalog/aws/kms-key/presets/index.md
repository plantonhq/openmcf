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
---

# KMS Key Presets

Ready-to-deploy configuration presets for KMS Key. Each preset is a complete manifest you can copy, customize, and deploy.
