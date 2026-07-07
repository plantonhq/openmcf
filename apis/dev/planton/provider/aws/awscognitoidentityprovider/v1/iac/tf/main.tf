# Cognito identity provider -- a SAML/OIDC/social IdP federated into a user
# pool. The referenced pool ID arrives pre-resolved as a plain string.

resource "aws_cognito_identity_provider" "this" {
  user_pool_id = var.spec.user_pool_id

  # Both are the IdP's AWS identity within its pool (ForceNew): app clients
  # list provider_name in supported_identity_providers.
  provider_name = var.spec.provider_name
  provider_type = var.spec.provider_type

  provider_details = local.provider_details

  # Empty means AWS applies its per-provider-type default mappings.
  attribute_mapping = length(var.spec.attribute_mapping) > 0 ? var.spec.attribute_mapping : null

  idp_identifiers = length(var.spec.idp_identifiers) > 0 ? var.spec.idp_identifiers : null

  lifecycle {
    # AWS auto-populates ActiveEncryptionCertificate for SAML providers.
    # Ignore it to prevent perpetual drift.
    ignore_changes = [
      provider_details["ActiveEncryptionCertificate"],
    ]
  }
}
