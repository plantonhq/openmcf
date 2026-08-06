# Store a bring-your-own TLS certificate on a Container App Environment --
# either an inline PFX upload or a Key Vault reference the environment
# keeps current across renewals. The certificate is shared by every app in
# the environment; AzureContainerAppCustomDomain bindings reference it by
# the certificate_id output.
#
# Lifecycle notes worth knowing before operating this resource:
# - Only tags update in place; every other change replaces the certificate
#   (and briefly re-binds any custom domain using it).
# - Azure never returns the PFX blob on reads, so blob drift is invisible
#   -- rotating an inline certificate means updating the spec.
# - The Key Vault path requires the environment's managed identity (the
#   one named in certificate_key_vault.identity) to already hold read
#   access to the vault's secrets; Azure checks it at deploy time.
resource "azurerm_container_app_environment_certificate" "main" {
  name                         = var.spec.certificate_name
  container_app_environment_id = var.spec.container_app_environment_id

  # Spec validation guarantees exactly one source: the inline PFX (with
  # its possibly-empty password -- passwordless PFX bundles are legal,
  # and Azure expects the password argument alongside the blob either
  # way) or the Key Vault reference. The unused arguments stay null so
  # the provider's own exactly-one-of contract is satisfied.
  certificate_blob_base64 = var.spec.certificate_blob_base64 != "" ? var.spec.certificate_blob_base64 : null
  certificate_password    = var.spec.certificate_blob_base64 != "" ? var.spec.certificate_password : null

  dynamic "certificate_key_vault" {
    for_each = var.spec.certificate_key_vault != null ? [var.spec.certificate_key_vault] : []
    content {
      key_vault_secret_id = certificate_key_vault.value.key_vault_secret_id
      # Unset deploys "System" -- Azure's own default identity for the
      # vault read; the explicit fallback keeps both engines identical.
      identity = certificate_key_vault.value.identity != "" ? certificate_key_vault.value.identity : "System"
    }
  }

  tags = local.final_tags
}
