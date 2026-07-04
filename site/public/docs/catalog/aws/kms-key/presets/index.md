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
    excerpt: "SYMMETRIC_DEFAULT key with enable_key_rotation, 365-day rotation period, and alias/ prefix aliases — the shape AWS service integrations expect."
---

# KMS Key Presets

Ready-to-deploy configuration presets for KMS Key. Each preset is a complete manifest you can copy, customize, and deploy.
