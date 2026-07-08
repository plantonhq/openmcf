---
title: "Presets"
description: "Ready-to-deploy configuration presets for KMS Key"
type: "preset-list"
componentSlug: "kms-key"
componentTitle: "KMS Key"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-symmetric-cmek"
    rank: "01"
    title: "Symmetric CMEK Key"
    excerpt: "The workhorse of customer-managed encryption: a symmetric `ENCRYPT_DECRYPT` key with automatic 90-day rotation, referenced by data services (BigQuery, Spanner, Cloud SQL, GKE, Pub/Sub, ...) through..."
  - slug: "02-hsm-symmetric-encryption"
    rank: "02"
    title: "HSM-Protected Symmetric Key"
    excerpt: "The compliance-grade CMEK key: identical to the symmetric workhorse but with every key version generated and held inside Cloud HSM (FIPS 140-2 Level 3 validated hardware)."
  - slug: "03-asymmetric-signing"
    rank: "03"
    title: "Asymmetric Signing Key"
    excerpt: "The code- and artifact-signing primitive: an `ASYMMETRIC_SIGN` key whose private material never leaves Cloud KMS, with public keys fetched by verifiers per version."
  - slug: "04-import-only-byok"
    rank: "04"
    title: "Import-Only BYOK Key"
    excerpt: "The bring-your-own-key container: a key that can only ever hold imported material, for organizations whose policy requires key generation outside Google's infrastructure."
---

# KMS Key Presets

Ready-to-deploy configuration presets for KMS Key. Each preset is a complete manifest you can copy, customize, and deploy.
