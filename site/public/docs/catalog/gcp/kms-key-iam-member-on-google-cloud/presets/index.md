---
title: "Presets"
description: "Ready-to-deploy configuration presets for KMS Key IAM Member on Google Cloud"
type: "preset-list"
componentSlug: "kms-key-iam-member-on-google-cloud"
componentTitle: "KMS Key IAM Member on Google Cloud"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-storage-cmek-grant"
    rank: "01"
    title: "Cloud Storage CMEK Grant"
    excerpt: "This preset grants `roles/cloudkms.cryptoKeyEncrypterDecrypter` on one key to Cloud Storage's service agent — the permission every CMEK-encrypted bucket requires before its first object write...."
  - slug: "02-workload-key-user"
    rank: "02"
    title: "Workload Key User Grant"
    excerpt: "This preset grants `roles/cloudkms.cryptoKeyEncrypterDecrypter` on one key to a workload's own service account — for applications that call the KMS API directly (envelope encryption, signing..."
  - slug: "03-conditional-key-access"
    rank: "03"
    title: "Conditional Key Access (Time-Bound)"
    excerpt: "This preset grants key access that expires on its own: an IAM condition gates the grant with a CEL expression, so a human operator gets encrypt/decrypt on one key until a fixed timestamp — no cleanup..."
---

# KMS Key IAM Member on Google Cloud Presets

Ready-to-deploy configuration presets for KMS Key IAM Member on Google Cloud. Each preset is a complete manifest you can copy, customize, and deploy.
