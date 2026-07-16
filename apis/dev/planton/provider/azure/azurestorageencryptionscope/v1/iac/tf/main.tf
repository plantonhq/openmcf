# The encryption scope, addressed by the parent account's ARM ID -- a
# pure management-plane resource (no data-plane path exists). Scopes
# carry no Azure tags: ARM does not support tags on encryptionScopes,
# so the platform's identity tags live on the account.
#
# Deletion is a SOFT-DISABLE: ARM has no true delete for scopes, so
# destroy flips the scope's state to Disabled (the provider then treats
# a Disabled scope as gone). The name stays reserved within the account
# -- recreating the same scope name re-enables it.
resource "azurerm_storage_encryption_scope" "main" {
  name               = var.spec.scope_name
  storage_account_id = var.spec.storage_account_id

  # Platform-managed (Microsoft.Storage) or customer-managed
  # (Microsoft.KeyVault) key ownership. With Key Vault, the ACCOUNT must
  # carry an identity with wrap/unwrap access on the key's vault -- the
  # same plumbing as the account-level customer-managed key.
  source           = local.source
  key_vault_key_id = local.key_vault_key_id

  # A second, independent platform-managed encryption layer for just
  # this scope's data -- independent of the account-level infrastructure
  # encryption switch. Fixed at creation.
  infrastructure_encryption_required = var.spec.infrastructure_encryption_required
}
