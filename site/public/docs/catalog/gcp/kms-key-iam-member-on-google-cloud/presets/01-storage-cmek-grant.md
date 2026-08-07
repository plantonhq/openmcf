---
title: "Cloud Storage CMEK Grant"
description: "This preset grants `roles/cloudkms.cryptoKeyEncrypterDecrypter` on one key to Cloud Storage's service agent — the permission every CMEK-encrypted bucket requires before its first object write...."
type: "preset"
rank: "01"
presetSlug: "01-storage-cmek-grant"
componentSlug: "kms-key-iam-member-on-google-cloud"
componentTitle: "KMS Key IAM Member on Google Cloud"
provider: "gcp"
icon: "package"
order: 1
---

# Cloud Storage CMEK Grant

This preset grants `roles/cloudkms.cryptoKeyEncrypterDecrypter` on one key to Cloud Storage's service agent — the permission every CMEK-encrypted bucket requires before its first object write. Because the key arrives by reference, the encrypted bucket can depend on this grant and deploy strictly after the permission exists, closing the first-deploy IAM-propagation race.

## When to Use

- A GcpGcsBucket sets a customer-managed encryption key and needs the storage agent authorized on it
- Any first CMEK deploy where the permission and the encrypted resource land in the same run
- Tightening a project-wide encrypter/decrypter grant down to the one key storage actually uses

## Key Configuration Choices

- **`valueFrom` cryptoKeyId** — the reference to GcpKmsKey's `key_id` output is what turns the grant into a DAG edge instead of a raced side effect
- **Project NUMBER in the agent email** — service agent emails embed the numeric project, never the project ID
- **Key-scoped, not ring-scoped** — the agent can use exactly this key; other keys on the ring stay off-limits

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<kms-key-resource-name>` | The Planton resource name of the GcpKmsKey | Your GcpKmsKey manifest's `metadata.name` |
| `<gcp-project-number>` | The numeric project number of the bucket's project | `GcpProject` outputs or GCP Console |

## Related Presets

- **02-workload-key-user** — Grant a workload's own service account use of a key
- **03-conditional-key-access** — Time-bound key access with an IAM condition
