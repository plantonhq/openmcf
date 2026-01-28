# Auth0Connection Main Resources
# This file creates the Auth0 connection based on the strategy type

# Auth0 Connection Resource
resource "auth0_connection" "this" {
  name                 = local.connection_name
  strategy             = local.strategy
  display_name         = local.display_name
  is_domain_connection = local.is_domain_connection
  show_as_button       = local.show_as_button
  realms               = length(local.realms) > 0 ? local.realms : null
  metadata             = length(local.connection_metadata) > 0 ? local.connection_metadata : null

  # Options block - configuration varies by strategy
  options {
    # Database connection options (auth0 strategy)
    password_policy        = local.is_database_strategy && local.database_options != null ? local.database_options.password_policy : null
    requires_username      = local.is_database_strategy && local.database_options != null ? local.database_options.requires_username : null
    disable_signup         = local.is_database_strategy && local.database_options != null ? local.database_options.disable_signup : null
    brute_force_protection = local.is_database_strategy && local.database_options != null ? local.database_options.brute_force_protection : null

    # Password history for database connections
    dynamic "password_history" {
      for_each = local.is_database_strategy && local.database_options != null && local.database_options.password_history_size > 0 ? [1] : []
      content {
        enable = true
        size   = local.database_options.password_history_size
      }
    }

    # Password no personal info for database connections
    dynamic "password_no_personal_info" {
      for_each = local.is_database_strategy && local.database_options != null && local.database_options.password_no_personal_info ? [1] : []
      content {
        enable = true
      }
    }

    # Password dictionary for database connections
    dynamic "password_dictionary" {
      for_each = local.is_database_strategy && local.database_options != null && local.database_options.password_dictionary ? [1] : []
      content {
        enable = true
      }
    }

    # MFA for database connections
    dynamic "mfa" {
      for_each = local.is_database_strategy && local.database_options != null && local.database_options.mfa_enabled ? [1] : []
      content {
        active                 = true
        return_enroll_settings = true
      }
    }

    # Social connection options (google-oauth2, facebook, github, etc.)
    client_id     = local.is_social_strategy && local.social_options != null ? local.social_options.client_id : (local.is_oidc_strategy && local.oidc_options != null ? local.oidc_options.client_id : (local.is_waad_strategy && local.azure_ad_options != null ? local.azure_ad_options.client_id : null))
    client_secret = local.is_social_strategy && local.social_options != null ? local.social_options.client_secret : (local.is_oidc_strategy && local.oidc_options != null ? local.oidc_options.client_secret : (local.is_waad_strategy && local.azure_ad_options != null ? local.azure_ad_options.client_secret : null))
    scopes        = local.is_social_strategy && local.social_options != null && length(local.social_options.scopes) > 0 ? local.social_options.scopes : (local.is_oidc_strategy && local.oidc_options != null ? local.oidc_options.scopes : null)

    # SAML connection options (samlp strategy)
    sign_in_endpoint    = local.is_saml_strategy && local.saml_options != null ? local.saml_options.sign_in_endpoint : null
    signing_cert        = local.is_saml_strategy && local.saml_options != null ? local.saml_options.signing_cert : null
    sign_out_endpoint   = local.is_saml_strategy && local.saml_options != null ? local.saml_options.sign_out_endpoint : null
    entity_id           = local.is_saml_strategy && local.saml_options != null ? local.saml_options.entity_id : null
    protocol_binding    = local.is_saml_strategy && local.saml_options != null ? local.saml_options.protocol_binding : null
    sign_saml_request   = local.is_saml_strategy && local.saml_options != null ? local.saml_options.sign_request : null
    signature_algorithm = local.is_saml_strategy && local.saml_options != null ? local.saml_options.signature_algorithm : null
    digest_algorithm    = local.is_saml_strategy && local.saml_options != null ? local.saml_options.digest_algorithm : null
    fields_map          = local.is_saml_strategy && local.saml_options != null && length(local.saml_options.attribute_mappings) > 0 ? jsonencode(local.saml_options.attribute_mappings) : null

    # OIDC connection options (oidc strategy)
    issuer                 = local.is_oidc_strategy && local.oidc_options != null ? local.oidc_options.issuer : null
    type                   = local.is_oidc_strategy && local.oidc_options != null ? local.oidc_options.type : null
    authorization_endpoint = local.is_oidc_strategy && local.oidc_options != null ? local.oidc_options.authorization_endpoint : null
    token_endpoint         = local.is_oidc_strategy && local.oidc_options != null ? local.oidc_options.token_endpoint : null
    jwks_uri               = local.is_oidc_strategy && local.oidc_options != null ? local.oidc_options.jwks_uri : null

    # Azure AD connection options (waad strategy)
    domain                = local.is_waad_strategy && local.azure_ad_options != null ? local.azure_ad_options.domain : null
    tenant_domain         = local.is_waad_strategy && local.azure_ad_options != null ? local.azure_ad_options.tenant_id : null
    max_groups_to_retrieve = local.is_waad_strategy && local.azure_ad_options != null ? tostring(local.azure_ad_options.max_groups_to_retrieve) : null
    api_enable_users      = local.is_waad_strategy && local.azure_ad_options != null ? local.azure_ad_options.api_enable_users : null
  }
}

# Enable clients for the connection
resource "auth0_connection_clients" "this" {
  count = length(local.enabled_clients) > 0 ? 1 : 0

  connection_id   = auth0_connection.this.id
  enabled_clients = local.enabled_clients
}
