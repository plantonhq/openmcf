---
title: "Team AI Hub"
description: "The 30-second hub: a shared AI Foundry foundation with a system-assigned identity on your existing key vault and storage account. Create it once; every team then gets its own AzureAiFoundryProject..."
type: "preset"
rank: "01"
presetSlug: "01-team-hub"
componentSlug: "ai-foundry-hub"
componentTitle: "AI Foundry Hub"
provider: "azure"
icon: "package"
order: 1
---

# Team AI Hub

The 30-second hub: a shared AI Foundry foundation with a
system-assigned identity on your existing key vault and storage
account. Create it once; every team then gets its own
AzureAiFoundryProject inside it and inherits this hub's security and
connectivity posture.

## When to Use

- Setting up AI Foundry for the first time
- One shared foundation, several team projects on top
- No customer-managed-key or network-isolation requirements yet

## Key Configuration Choices

- `identity.type: SYSTEM_ASSIGNED` -- Azure creates and rotates the
  hub's identity; grant it access on the key vault and storage after
  creation (or reference Planton-managed ones and compose the grants).
- Public network access stays enabled (the default) -- tighten later
  with `publicNetworkAccessEnabled: false` plus private endpoints.
- No `managedNetwork` block -- isolation stays off; the mode can be
  turned on in place later.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-resource-group-name>` | The resource group to create the hub in | Portal -> Resource groups, or reference an AzureResourceGroup |
| `<your-key-vault-id>` | ARM ID of the key vault for hub secrets | Portal -> Key Vault -> Properties -> Resource ID, or reference an AzureKeyVault |
| `<your-storage-account-id>` | ARM ID of the hub's storage account (HNS off) | Portal -> Storage account -> Properties -> Resource ID, or reference an AzureStorageAccount |

## Related Presets

- `02-cmk-hardened-hub` -- customer-managed-key encryption, private
  access, and managed-network isolation for regulated estates.
