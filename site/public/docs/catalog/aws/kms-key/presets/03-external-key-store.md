---
title: "External Key Store Key"
description: "This preset creates a KMS key whose cryptographic material lives OUTSIDE standard AWS KMS: in an external key manager reached through an external key store (XKS), or -- by dropping `xks_key_id` -- in..."
type: "preset"
rank: "03"
presetSlug: "03-external-key-store"
componentSlug: "kms-key"
componentTitle: "KMS Key"
provider: "aws"
icon: "package"
order: 3
---

# External Key Store Key

This preset creates a KMS key whose cryptographic material lives OUTSIDE standard AWS KMS: in an external key manager reached through an external key store (XKS), or -- by dropping `xks_key_id` -- in your own CloudHSM cluster behind a CloudHSM key store. AWS services and SDKs use the key exactly like any other KMS key; KMS forwards every operation to the backing store.

## When to Use

- Regulatory or sovereignty requirements that key material never resides in AWS-managed HSMs
- Hold-your-own-key (HYOK) postures where an on-premises or third-party key manager is the root of trust
- CloudHSM-backed workloads that need single-tenant FIPS 140-2 Level 3 hardware under your control

## Key Configuration Choices

- **External key store** (`custom_key_store_id` + `xks_key_id`) — the store must exist and be CONNECTED before the key is created; `xks_key_id` names the existing key in your external manager
- **CloudHSM variant** — set `custom_key_store_id` alone (pointing at a CloudHSM key store): KMS generates the material inside your cluster
- **Symmetric-only, by AWS contract** — custom key store keys must be `SYMMETRIC_DEFAULT` / `ENCRYPT_DECRYPT` (leave `key_spec`/`key_usage` empty), never rotate automatically, and cannot be multi-Region; the spec enforces all three at validate time
- **Store lifecycle is yours** — key stores are account-level infrastructure the catalog does not provision; supply the literal store id

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aws-region>` | AWS region where the key lives | Your deployment region |
| `<key-description>` | Purpose of the key | Your team's naming conventions |
| `alias/xks-backed-key` | Rename to your alias | Must match `alias/[0-9A-Za-z_/-]+` |
| `<custom-key-store-id>` | The connected key store (cks-...) | KMS console → Custom key stores, or DescribeCustomKeyStores |
| `<xks-key-id>` | The existing key's id in your external manager | Your external key manager (omit for CloudHSM stores) |

## Related Presets

- `01-symmetric-encryption` — standard KMS-held material with rotation
- `02-cross-account-grants` — delegated access on any key, including this one
