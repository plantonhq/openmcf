# Cloudflare Access identity provider: how users sign in to Access-protected
# applications. The TYPE is immutable -- Terraform replaces the provider (new
# ID) if it changes, invalidating policy rules that reference the old one.
#
# The SCIM secret (scim_config.secret) is minted once when SCIM is first
# enabled and redacted on every later read -- it does not survive an import.
resource "cloudflare_zero_trust_access_identity_provider" "main" {
  account_id = local.account_id != "" ? local.account_id : null
  zone_id    = local.zone_id != "" ? local.zone_id : null

  name = var.spec.name
  type = var.spec.type

  config = local.config

  saml_certificate_set_id = var.spec.saml_certificate_set_id != "" ? var.spec.saml_certificate_set_id : null

  scim_config = local.scim_config

  # A safety latch: while true, Cloudflare refuses API updates and deletes.
  read_only = var.spec.read_only ? true : null
}
