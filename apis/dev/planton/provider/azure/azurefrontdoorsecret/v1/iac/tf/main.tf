# The Front Door secret -- the bring-your-own TLS certificate node.
# Fully immutable: Azure exposes no update on secrets, so ANY change
# replaces the resource. That is safe in practice because certificate
# ROTATION happens inside Key Vault: a VERSIONLESS
# key_vault_certificate_id makes Front Door follow the certificate's
# latest version automatically, while a versioned id pins one version.
#
# Operational prerequisite: Front Door's own service principal (the
# Microsoft.AzureFrontDoor-Cdn enterprise application) must hold read
# access on the vault's certificates/secrets before this deploys -- a
# one-time grant per tenant/vault. Without it Azure rejects the create
# with an access-denied error naming the vault.
#
# No Azure tags: ARM does not support tags on Front Door secrets.
resource "azurerm_cdn_frontdoor_secret" "main" {
  name                     = var.spec.secret_name
  cdn_frontdoor_profile_id = var.spec.profile_id

  secret {
    customer_certificate {
      key_vault_certificate_id = var.spec.key_vault_certificate_id
    }
  }
}
