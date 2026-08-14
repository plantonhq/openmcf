---
title: "Point-to-Site with Entra ID"
description: "This preset creates a VPN gateway whose point-to-site clients authenticate with Entra ID (Azure AD) over OpenVPN -- remote-workforce VPN access with your existing identity provider, no certificate..."
type: "preset"
rank: "03"
presetSlug: "03-point-to-site-entra-id"
componentSlug: "virtual-network-gateway"
componentTitle: "Virtual Network Gateway"
provider: "azure"
icon: "package"
order: 3
---

# Point-to-Site with Entra ID

This preset creates a VPN gateway whose point-to-site clients authenticate with Entra ID (Azure AD) over OpenVPN -- remote-workforce VPN access with your existing identity provider, no certificate distribution, and conditional-access policies applied at sign-in.

## When to Use

- Individual users (developers, operators) reaching private VNet resources from anywhere
- Replacing certificate-based P2S (and its enrollment/revocation burden) with identity-based access
- Organizations that want conditional access and MFA on VPN sign-in

## Key Configuration Choices

- **The Entra ID trio travels together** -- `aadTenant`, `aadAudience`, and `aadIssuer` (spec-enforced). The audience is the Azure VPN client application's id; the issuer ends with a trailing slash
- **`vpnClientProtocols: [OpenVPN]`** -- Entra ID authentication requires OpenVPN; IKEv2/SSTP pair with certificate or RADIUS auth instead
- **Client pool sizing** -- each connected client consumes one address from `addressSpaces`; carve a range that never overlaps the VNet or any tunneled network
- **VPN_GW_1_AZ floor** -- IKEv2/OpenVPN point-to-site needs VpnGw1AZ or higher (BASIC cannot; the non-AZ VpnGw SKUs are retired for new creates). The same gateway can carry site-to-site tunnels simultaneously
- **Per-group routing** -- add `policyGroups` and `vpnClientConfiguration.clientConnections` when different user groups need distinct address pools

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-gateway-subnet-resource-name>` | Planton metadata name of the `AzureSubnet` named "GatewaySubnet" | Your subnet resource |
| `<your-public-ip-resource-name>` | Planton metadata name of the `AzurePublicIp` | Your public IP resource |
| `<your-tenant-id>` | Entra ID tenant id | Entra admin center -> Overview |
| `<azure-vpn-application-client-id>` | The Azure VPN enterprise application's client id | Entra admin center -> Enterprise applications -> Azure VPN |
