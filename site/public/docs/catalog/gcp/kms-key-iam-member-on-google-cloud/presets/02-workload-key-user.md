---
title: "Workload Key User Grant"
description: "This preset grants `roles/cloudkms.cryptoKeyEncrypterDecrypter` on one key to a workload's own service account — for applications that call the KMS API directly (envelope encryption, signing..."
type: "preset"
rank: "02"
presetSlug: "02-workload-key-user"
componentSlug: "kms-key-iam-member-on-google-cloud"
componentTitle: "KMS Key IAM Member on Google Cloud"
provider: "gcp"
icon: "package"
order: 2
---

# Workload Key User Grant

This preset grants `roles/cloudkms.cryptoKeyEncrypterDecrypter` on one key to a workload's own service account — for applications that call the KMS API directly (envelope encryption, signing payloads, decrypting configuration) rather than relying on a Google service's CMEK integration.

## When to Use

- An application encrypts/decrypts data itself via the Cloud KMS API
- Envelope encryption designs where the app wraps data keys with a KMS key
- Both the key and the identity are Planton-managed nodes and the access edge should be visible in the graph

## Key Configuration Choices

- **Both sides referenced** — the key's `key_id` output and the account's `member` output keep the entire access relationship in the resource graph
- **Key-scoped least privilege** — the workload can use exactly this key; a project- or ring-level grant would expose every other key too
- **Encrypter/decrypter only** — the role covers encrypt/decrypt operations, not key administration; the workload cannot rotate, disable, or re-policy the key

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<kms-key-resource-name>` | The Planton resource name of the GcpKmsKey | Your GcpKmsKey manifest's `metadata.name` |
| `<service-account-resource-name>` | The workload's GcpServiceAccount | Your GcpServiceAccount manifest's `metadata.name` |

## Related Presets

- **01-storage-cmek-grant** — Authorize a Google service agent for CMEK
- **03-conditional-key-access** — Time-bound key access with an IAM condition
