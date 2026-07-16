---
title: "Symmetric CMEK Key"
description: "The workhorse of customer-managed encryption: a symmetric `ENCRYPT_DECRYPT` key with automatic 90-day rotation, referenced by data services (BigQuery, Spanner, Cloud SQL, GKE, Pub/Sub, ...) through..."
type: "preset"
rank: "01"
presetSlug: "01-symmetric-cmek"
componentSlug: "kms-key"
componentTitle: "KMS Key"
provider: "gcp"
icon: "package"
order: 1
---

# Symmetric CMEK Key

The workhorse of customer-managed encryption: a symmetric
`ENCRYPT_DECRYPT` key with automatic 90-day rotation, referenced by data
services (BigQuery, Spanner, Cloud SQL, GKE, Pub/Sub, ...) through their
CMEK fields.

## What this preset creates

A software-protected symmetric key on a referenced `GcpKmsKeyRing`, using
GCP's default `GOOGLE_SYMMETRIC_ENCRYPTION` algorithm. Rotation mints a
new primary version every 90 days — new data is encrypted under the new
version while old versions keep decrypting existing data.

## Prerequisites

- A `GcpKmsKeyRing` named `prod-encryption` in the location where your
  data lives (see the `GcpKmsKeyRing` presets). Most CMEK integrations
  require the key and the protected resource to be co-located.

## Composing CMEK

Point any consumer's KMS reference at this key's `key_id` output — for
example a BigQuery dataset's `kmsKeyName` or a GKE cluster's database
encryption key. Remember to grant each service agent
`roles/cloudkms.cryptoKeyEncrypterDecrypter` on the key.

## Remix ideas

- Tighten `rotationPeriod` to `2592000s` (30 days) for stricter policies.
- Add `destroyScheduledDuration` to lengthen the recovery window for
  destroyed versions beyond the 30-day default.
- Switch to the HSM preset when compliance requires FIPS 140-2 Level 3.
