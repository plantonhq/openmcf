---
title: "Key Vault Certificate"
description: "The production posture: the environment pulls the certificate from Azure Key Vault with its system-assigned identity and follows renewals automatically -- no rotation treadmill, no key material in..."
type: "preset"
rank: "01"
presetSlug: "01-key-vault-certificate"
componentSlug: "container-app-environment-certificate"
componentTitle: "Container App Environment Certificate"
provider: "azure"
icon: "package"
order: 1
---

# Key Vault Certificate

The production posture: the environment pulls the certificate from Azure Key Vault with its system-assigned identity and follows renewals automatically -- no rotation treadmill, no key material in configuration.

## When to Use

- Certificates already managed (and auto-renewed) in Key Vault
- Any binding where you never want to touch TLS again after setup

## Key Configuration Choices

- `keyVaultSecretId` references the certificate's VERSIONLESS secret face -- renewals propagate; a versioned URL pins forever
- `identity` is omitted -- the environment's system-assigned identity reads the vault ("System", Azure's default). Set a UAI reference when the environment authenticates with one
- Grant the identity read access to the vault's secrets first (Key Vault Secrets User under RBAC) -- compose an `AzureRoleAssignment`; Azure checks the permission at deploy time

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `app.example.com` | Replace with the certificate's name on the environment (commonly the domain) | The custom domain you are binding |
| `<your-environment>` | The AzureContainerAppEnvironment resource's metadata name | Your Planton resource inventory |
| `<your-certificate>` | The AzureKeyVaultCertificate resource's metadata name | Your Planton resource inventory |

## Related Presets

- `02-inline-pfx-certificate` -- upload a PFX directly, for certificates outside Key Vault
