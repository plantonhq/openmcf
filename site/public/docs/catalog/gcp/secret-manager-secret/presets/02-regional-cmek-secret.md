---
title: "Regional CMEK Secret"
description: "The data-residency posture: a REGIONAL secret whose payloads never leave the region, encrypted with a customer-managed KMS key, with both destroy guards armed."
type: "preset"
rank: "02"
presetSlug: "02-regional-cmek-secret"
componentSlug: "secret-manager-secret"
componentTitle: "Secret Manager Secret"
provider: "gcp"
icon: "package"
order: 2
---

# Regional CMEK Secret

The data-residency posture: a REGIONAL secret whose payloads never
leave the region, encrypted with a customer-managed KMS key, with both
destroy guards armed.

## What it configures

- `region` — selects the regional Secret Manager API; scope is
  permanent (global <-> regional means destroy-and-recreate).
- `customerManagedEncryption` — CMEK attaches directly on regional
  secrets (replication is a global-only concept). The key must live in
  the SAME region, and the Secret Manager service agent needs
  `roles/cloudkms.cryptoKeyEncrypterDecrypter` on it.
- `deletionProtection: true` + `deletionPolicy: PREVENT` — the
  belt-and-suspenders posture for regulated credentials.

## Adjust before deploying

- **region / kmsKey** — your residency region and key; reference a
  GcpKmsKey resource via valueFrom (its `key_id` output).
- Add `iamMembers` for each consuming workload (see the **App Secret
  with Access** preset).

## When to choose something else

Without a residency requirement, global automatic replication (the
**App Secret with Access** preset) is more available and simpler.
