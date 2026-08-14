---
title: "MACsec Port"
description: "This preset creates a MACsec-encrypted ExpressRoute Direct port: layer-2 encryption on both physical links, keyed from your own Key Vault secrets through a user-assigned managed identity. Both spec..."
type: "preset"
rank: "02"
presetSlug: "02-macsec-port"
componentSlug: "expressroute-port"
componentTitle: "ExpressRoute Port"
provider: "azure"
icon: "package"
order: 2
---

# MACsec Port

This preset creates a MACsec-encrypted ExpressRoute Direct port: layer-2 encryption on both physical links, keyed from your own Key Vault secrets through a user-assigned managed identity. Both spec contracts -- the CKN/CAK pair travelling together, and the user-assigned identity being present -- are enforced at validation time.

## When to Use

- Compliance regimes that require link-layer encryption between your edge and Microsoft's
- Any estate where the colocation cross-connects traverse infrastructure you do not control

## Key Configuration Choices

- **Three parties must align** -- the Key Vault secrets (CKN names the key, CAK is the key material), the port's identity (Key Vault secret GET on both), and the facility side's matching key/SCI configuration; a mismatch reads as a dead link
- **GCM_AES_256 is the strong default choice** -- use the XPN variants for sustained rates above ~40 Gbps where 32-bit packet numbering wraps too quickly
- **Grant Key Vault access BEFORE deploying** -- reference an AzureUserAssignedIdentity so the grant composes ahead of the port

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | The ARM metadata region | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-peering-location>` | The ExpressRoute Direct facility | `az network express-route port location list` |
| `<your-user-assigned-identity>` | The AzureUserAssignedIdentity resource name | Your Planton environment's resources |
| `<your-ckn-secret-id>` / `<your-cak-secret-id>` | Versioned Key Vault secret URLs holding the MACsec keys | Key Vault portal -> the secret -> Secret Identifier |
