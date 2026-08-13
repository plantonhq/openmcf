---
title: "CMK-Hardened AI Hub"
description: "The regulated-estate hub: customer-managed-key encryption, disabled public access, approved-outbound network isolation, and reduced diagnostics (high business impact). Every project created in this..."
type: "preset"
rank: "02"
presetSlug: "02-cmk-hardened-hub"
componentSlug: "ai-foundry-hub"
componentTitle: "AI Foundry Hub"
provider: "azure"
icon: "package"
order: 2
---

# CMK-Hardened AI Hub

The regulated-estate hub: customer-managed-key encryption, disabled
public access, approved-outbound network isolation, and reduced
diagnostics (high business impact). Every project created in this hub
inherits the posture.

## When to Use

- Compliance regimes requiring customer-owned encryption keys
- Estates where AI workloads must not reach the public internet
- Data-sensitivity postures that minimize Microsoft-collected
  diagnostics

## Key Configuration Choices

- `identity.type: USER_ASSIGNED` with `primaryUserAssignedIdentity` --
  the identity (and its wrap/unwrap grant on the key) must exist
  BEFORE the hub, which a system-assigned identity cannot satisfy.
- `encryption.keyId` is a VERSIONED key URL -- the provider's hub
  contract; key rotation does NOT auto-propagate (re-point the field
  to rotate). This deliberately differs from the classic ML
  workspace's versionless guidance.
- `highBusinessImpactEnabled: true` -- explicit here; the service
  flips it true anyway when encryption is on, and setting it keeps the
  intent visible in the manifest.
- Encryption and the identity flavor are fixed at creation -- moving
  an unencrypted hub to CMK replaces it.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-resource-group-name>` | The resource group to create the hub in | Portal -> Resource groups |
| `<your-key-vault-id>` | ARM ID of the key vault (holds hub secrets AND the CMK) | Portal -> Key Vault -> Properties -> Resource ID |
| `<your-storage-account-id>` | ARM ID of the hub's storage account (HNS off) | Portal -> Storage account -> Properties -> Resource ID |
| `<your-user-assigned-identity-id>` | ARM ID of the identity holding key access | Portal -> Managed Identities -> Properties -> Resource ID |
| `keyId` example URL | Your key's VERSIONED data-plane URL | Portal -> Key Vault -> Keys -> the key -> the version -> Key Identifier |

## Related Presets

- `01-team-hub` -- the plain shared hub without CMK or isolation.
