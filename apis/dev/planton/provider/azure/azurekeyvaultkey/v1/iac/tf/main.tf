# A Key Vault key is a data-plane object: the provider talks to the vault's
# {name}.vault.azure.net endpoint, not ARM -- which is why creation fails
# with a 403 when the deploying credential lacks data-plane key permissions
# on the vault, even if it owns the subscription.
resource "azurerm_key_vault_key" "main" {
  name         = var.spec.name
  key_vault_id = var.spec.key_vault_id

  # Type, size, and curve are fixed at creation -- Azure key material is
  # immutable by design; changing any of them replaces the key (and every
  # consumer re-encrypts through the new key on its next unwrap).
  key_type = local.key_type
  key_size = var.spec.key_size
  curve    = local.curve

  # The capability boundary: Azure rejects any operation not listed here.
  key_opts = local.key_opts

  not_before_date = var.spec.not_before_date
  expiration_date = var.spec.expiration_date

  # Rotation policy rides the key resource itself (Azure models it as a
  # sub-resource updated in place). expire_after stamps an expiry on every
  # NEW version; the automatic block is what actually rotates.
  dynamic "rotation_policy" {
    for_each = var.spec.rotation_policy != null ? [var.spec.rotation_policy] : []
    content {
      expire_after         = rotation_policy.value.expire_after
      notify_before_expiry = rotation_policy.value.notify_before_expiry

      dynamic "automatic" {
        for_each = rotation_policy.value.automatic != null ? [rotation_policy.value.automatic] : []
        content {
          time_after_creation = automatic.value.time_after_creation
          time_before_expiry  = automatic.value.time_before_expiry
        }
      }
    }
  }

  tags = local.final_tags
}
