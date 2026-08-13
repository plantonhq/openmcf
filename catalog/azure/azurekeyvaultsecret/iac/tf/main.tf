# A Key Vault secret is a data-plane object: the provider talks to the
# vault's {name}.vault.azure.net endpoint, not ARM -- which is why
# creation fails with a 403 when the deploying credential lacks
# data-plane secret permissions on the vault, even if it owns the
# subscription.
#
# The value arrives already reference-resolved (the spec's sensitive
# StringValueOrRef); the provider's write-only value_wo variant is
# deliberately not wired -- it duplicates `value` for a
# plaintext-in-config problem this module does not have.
resource "azurerm_key_vault_secret" "main" {
  name         = var.spec.name
  key_vault_id = var.spec.key_vault_id

  # Changing the value creates a NEW secret version; the versioned
  # outputs move with it, versionless_id keeps resolving to latest.
  value = var.spec.value

  content_type = var.spec.content_type != "" ? var.spec.content_type : null

  # Advisory attributes -- Key Vault stores and returns them; enforcing
  # them is the consumer's job.
  not_before_date = var.spec.not_before_date
  expiration_date = var.spec.expiration_date

  tags = local.final_tags
}
