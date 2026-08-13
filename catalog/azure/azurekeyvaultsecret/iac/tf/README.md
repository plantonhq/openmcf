# AzureKeyVaultSecret Terraform Module

## Overview

Stores a secret in an Azure Key Vault -- a versioned data-plane object the vault guards. The value arrives already reference-resolved (the spec's sensitive field); it is never written to outputs.

## Resources Created

- `azurerm_key_vault_secret` -- the secret

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureKeyVaultSecretSpec fields; the vault reference and the sensitive value arrive as resolved literals

## Outputs

- `secret_id` -- the versioned data-plane ID (pins consumers to this version)
- `versionless_id` -- the versionless data-plane ID (the reference consumers should use; value updates propagate automatically)
- `secret_name` -- the secret's name within the vault
- `version` -- the current version identifier
- `resource_id` / `resource_versionless_id` -- the versioned/versionless ARM resource IDs (control-plane identity)

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **A value change creates a NEW version** -- the versioned outputs move with it; `versionless_id` keeps resolving to latest.
- **The provider's write-only `value_wo`/`value_wo_version` pair is deliberately not wired** -- it duplicates `value` for a plaintext-in-config problem this module does not have (values are reference-resolved at deploy, never stored in manifests).
- **Data-plane permissions are separate from ARM**: creation fails with a 403 when the deploying credential lacks secret permissions on the vault, even if it owns the subscription.
- **Destroy purges the soft-deleted secret** (the provider's default feature behavior) so the name frees immediately -- skipped automatically when the vault has purge protection, where the name stays reserved for the retention window by design.
- **Key Vault strips raw newlines** -- multi-line values should be base64-encoded, with the encoding recorded in `content_type`.

## Required Permissions

The deploying principal needs the "Key Vault Administrator" or "Key Vault Secrets Officer" RBAC role on the vault (or secret permissions in a legacy access policy), plus ARM read on the vault.
