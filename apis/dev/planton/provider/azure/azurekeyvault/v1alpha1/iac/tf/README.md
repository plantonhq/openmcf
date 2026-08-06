# AzureKeyVault Terraform Module

## Overview

This Terraform module provisions an Azure Key Vault using the `azurerm`
provider. It creates a single `azurerm_key_vault` -- the tenant-scoped
container whose vault-wide controls (authorization mode, network rules,
deletion safety, SKU) everything inside inherits.

The vault authenticates against the deploying credential's Azure AD tenant,
read from the ambient client configuration -- a vault cannot be managed
cross-tenant, so the tenant is never modeled as an input.

What lives inside the vault is composed, never created here: keys are
`AzureKeyVaultKey` resources, certificates are `AzureKeyVaultCertificate`
resources, and secret values belong to a secrets-management workflow rather
than IaC state.

## Resources Created

- `azurerm_key_vault.main` -- the vault, with inline legacy access policies
  when the spec carries them (inline so the vault owns its complete grant
  list; azurerm cannot mix inline and standalone policy resources on one
  vault without perpetual drift)

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Key Vault specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` | yes | Azure region; changing it replaces the vault |
| `resource_group` | yes | Resource group name |
| `vault_name` | yes | 3-24 chars, globally unique (becomes `{name}.vault.azure.net`) |
| `sku` | no | `STANDARD` (default) or `PREMIUM` (HSM-backed key types); updatable |
| `rbac_authorization_enabled` | no | Azure RBAC (default true, the recommended posture) vs legacy access policies |
| `access_policies` | no | Legacy grants; permission lists arrive as the spec enums' full names and are mapped to ARM's values through exhaustive vocabularies in locals |
| `enabled_for_deployment` / `_disk_encryption` / `_template_deployment` | no | Resource-manager integration switches (Azure defaults: false) |
| `public_network_access_enabled` | no | false = private-endpoints-only (Azure default: true) |
| `purge_protection_enabled` | no | Irreversible once on; destroy then schedules deletion instead of purging (Azure default: false) |
| `soft_delete_retention_days` | no | 7-90 (Azure default 90); fixed at creation |
| `network_acls` | no | default_action + bypass + IP/subnet allowlists; an unspecified bypass materializes Azure's default (AzureServices) |
| `tags` | no | User tags, merged over metadata-derived tags (user wins) |

## Outputs

| Output | Description |
|--------|-------------|
| `key_vault_id` | The vault's ARM resource ID -- what keys, certificates, and vault-scoped role assignments reference |
| `key_vault_name` | The vault's name |
| `vault_uri` | The data-plane URI applications call |
| `tenant_id` | The Azure AD tenant the vault authenticates against |
| `resource_group_name` | The resource group the vault was created in |

## Operational Notes

- **Soft delete is always on.** With the provider's default features, a
  destroy purges the soft-deleted vault so the globally unique name frees
  up immediately -- unless purge protection is on, in which case destroy
  becomes a scheduled deletion at the end of the retention window.
- **A name-colliding create against a soft-deleted vault auto-recovers it**
  (provider default) rather than failing.
- Switching authorization modes on a live vault requires
  Microsoft.Authorization write permission (Owner / User Access
  Administrator) on the vault.
