---
title: "Entra ID Remote Workforce"
description: "This preset authenticates remote users with Entra ID (Azure AD): sign-in rides your tenant's conditional access and MFA, revoking a person is an identity operation, and no certificates are..."
type: "preset"
rank: "01"
presetSlug: "01-entra-remote-workforce"
componentSlug: "vpn-server-configuration"
componentTitle: "VPN Server Configuration"
provider: "azure"
icon: "package"
order: 1
---

# Entra ID Remote Workforce

This preset authenticates remote users with Entra ID (Azure AD): sign-in rides your tenant's conditional access and MFA, revoking a person is an identity operation, and no certificates are distributed. The standard choice for a managed workforce.

## When to Use

- Remote employees on managed devices with Entra ID accounts
- Organizations that want conditional access and MFA on VPN sign-in
- Anywhere certificate distribution is the thing you are trying to avoid

## Key Configuration Choices

- **Entra ID only** -- `vpnAuthenticationTypes: [AAD]`; the `aadAuthentication` trio carries the tenant plumbing
- **The audience is fixed** -- `41b23e61-6c1e-4545-b367-cd054e0ed4b4` is the Microsoft-managed Azure VPN Client's application ID; only the tenant/issuer URLs embed your directory
- **OpenVPN** -- Entra ID authentication requires the OpenVPN tunnel protocol

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group-name>` | Resource group for the configuration object | `AzureResourceGroup` status outputs (`resource_group_name`), or reference it with valueFrom |
| `<your-tenant-id>` | Your Entra ID directory (tenant) ID | Azure portal → Microsoft Entra ID → Overview → Tenant ID |

The issuer URL keeps its trailing slash -- a typo here deploys fine and then fails every sign-in.
