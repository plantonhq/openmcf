# Cognito User Pool app client -- the OAuth 2.0 / OIDC contract between ONE
# application and a user pool. The referenced pool ID arrives pre-resolved as
# a plain string; supported identity providers resolve the same way (an
# AwsCognitoIdentityProvider's provider_name output, or literals like
# "COGNITO"/"Google").

resource "aws_cognito_user_pool_client" "this" {
  name         = local.resource_name
  user_pool_id = var.spec.user_pool_id

  # ForceNew: confidential (secret-holding) vs public (PKCE) is decided at
  # creation.
  generate_secret = var.spec.generate_secret

  # ---------------------------------------------------------------------------
  # OAuth 2.0 / OIDC contract. Empty collections must reach AWS as absence --
  # the provider treats [] and unset differently on this Framework resource.
  # ---------------------------------------------------------------------------

  allowed_oauth_flows_user_pool_client = var.spec.allowed_oauth_flows_user_pool_client
  allowed_oauth_flows                  = length(var.spec.allowed_oauth_flows) > 0 ? var.spec.allowed_oauth_flows : null
  allowed_oauth_scopes                 = length(var.spec.allowed_oauth_scopes) > 0 ? var.spec.allowed_oauth_scopes : null
  callback_urls                        = length(var.spec.callback_urls) > 0 ? var.spec.callback_urls : null
  logout_urls                          = length(var.spec.logout_urls) > 0 ? var.spec.logout_urls : null
  default_redirect_uri                 = var.spec.default_redirect_uri != "" ? var.spec.default_redirect_uri : null
  supported_identity_providers         = length(var.spec.supported_identity_providers) > 0 ? var.spec.supported_identity_providers : null

  # ---------------------------------------------------------------------------
  # Authentication flows
  # ---------------------------------------------------------------------------

  explicit_auth_flows   = length(var.spec.explicit_auth_flows) > 0 ? var.spec.explicit_auth_flows : null
  auth_session_validity = var.spec.auth_session_validity

  # ---------------------------------------------------------------------------
  # Token lifetimes. Values pair with token_validity_units; omitted values
  # keep AWS defaults (1h access/ID, 30d refresh).
  # ---------------------------------------------------------------------------

  access_token_validity  = var.spec.access_token_validity
  id_token_validity      = var.spec.id_token_validity
  refresh_token_validity = var.spec.refresh_token_validity

  dynamic "token_validity_units" {
    for_each = var.spec.token_validity_units != null ? [var.spec.token_validity_units] : []
    content {
      access_token  = token_validity_units.value.access_token != "" ? token_validity_units.value.access_token : null
      id_token      = token_validity_units.value.id_token != "" ? token_validity_units.value.id_token : null
      refresh_token = token_validity_units.value.refresh_token != "" ? token_validity_units.value.refresh_token : null
    }
  }

  # Rotating refresh tokens shrink the blast radius of a stolen token; the
  # grace period absorbs clients that lose the response carrying the new one.
  dynamic "refresh_token_rotation" {
    for_each = var.spec.refresh_token_rotation != null ? [var.spec.refresh_token_rotation] : []
    content {
      feature                    = refresh_token_rotation.value.feature
      retry_grace_period_seconds = refresh_token_rotation.value.retry_grace_period_seconds
    }
  }

  # ---------------------------------------------------------------------------
  # Security posture
  # ---------------------------------------------------------------------------

  # Presence-aware: absent means AWS's default (true). Only an explicit choice
  # is forwarded so the module never silently flips revocability.
  enable_token_revocation = var.spec.enable_token_revocation

  enable_propagate_additional_user_context_data = var.spec.enable_propagate_additional_user_context_data ? true : null
  prevent_user_existence_errors                 = var.spec.prevent_user_existence_errors != "" ? var.spec.prevent_user_existence_errors : null

  # ---------------------------------------------------------------------------
  # Attribute access. Omitted means AWS grants access to all (mutable)
  # attributes.
  # ---------------------------------------------------------------------------

  read_attributes  = length(var.spec.read_attributes) > 0 ? var.spec.read_attributes : null
  write_attributes = length(var.spec.write_attributes) > 0 ? var.spec.write_attributes : null

  # ---------------------------------------------------------------------------
  # Pinpoint analytics. Exactly one identity arm (spec CEL): the ARN arm
  # derives the publish role; the ID arm wires it explicitly.
  # ---------------------------------------------------------------------------

  dynamic "analytics_configuration" {
    for_each = var.spec.analytics_configuration != null ? [var.spec.analytics_configuration] : []
    content {
      application_arn  = analytics_configuration.value.application_arn != "" ? analytics_configuration.value.application_arn : null
      application_id   = analytics_configuration.value.application_id != "" ? analytics_configuration.value.application_id : null
      external_id      = analytics_configuration.value.application_id != "" ? analytics_configuration.value.external_id : null
      role_arn         = analytics_configuration.value.application_id != "" ? analytics_configuration.value.role_arn : null
      user_data_shared = analytics_configuration.value.user_data_shared
    }
  }
}
