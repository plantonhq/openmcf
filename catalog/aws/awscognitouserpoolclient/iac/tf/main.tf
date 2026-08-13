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

# ---------------------------------------------------------------------------
# Client-scoped risk configuration (threat protection's automated responses
# for THIS client only). Overrides the pool-wide configuration set on the
# AwsCognitoUserPool spec. Requires the pool's threat protection to be active
# (advanced_security_mode AUDIT or ENFORCED) -- a cross-resource requirement
# AWS enforces at apply time.
# ---------------------------------------------------------------------------

resource "aws_cognito_risk_configuration" "this" {
  count = var.spec.risk_configuration != null ? 1 : 0

  user_pool_id = var.spec.user_pool_id
  client_id    = aws_cognito_user_pool_client.this.id

  dynamic "account_takeover_risk_configuration" {
    for_each = var.spec.risk_configuration.account_takeover != null ? [var.spec.risk_configuration.account_takeover] : []
    content {
      # The provider requires the actions block; the spec's CEL requires at
      # least one action inside it.
      actions {
        dynamic "low_action" {
          for_each = account_takeover_risk_configuration.value.low_action != null ? [account_takeover_risk_configuration.value.low_action] : []
          content {
            event_action = low_action.value.event_action
            notify       = low_action.value.notify
          }
        }

        dynamic "medium_action" {
          for_each = account_takeover_risk_configuration.value.medium_action != null ? [account_takeover_risk_configuration.value.medium_action] : []
          content {
            event_action = medium_action.value.event_action
            notify       = medium_action.value.notify
          }
        }

        dynamic "high_action" {
          for_each = account_takeover_risk_configuration.value.high_action != null ? [account_takeover_risk_configuration.value.high_action] : []
          content {
            event_action = high_action.value.event_action
            notify       = high_action.value.notify
          }
        }
      }

      dynamic "notify_configuration" {
        for_each = account_takeover_risk_configuration.value.notify_configuration != null ? [account_takeover_risk_configuration.value.notify_configuration] : []
        content {
          source_arn = notify_configuration.value.source_arn
          from       = notify_configuration.value.from != "" ? notify_configuration.value.from : null
          reply_to   = notify_configuration.value.reply_to != "" ? notify_configuration.value.reply_to : null

          dynamic "block_email" {
            for_each = notify_configuration.value.block_email != null ? [notify_configuration.value.block_email] : []
            content {
              subject   = block_email.value.subject
              html_body = block_email.value.html_body
              text_body = block_email.value.text_body
            }
          }

          dynamic "mfa_email" {
            for_each = notify_configuration.value.mfa_email != null ? [notify_configuration.value.mfa_email] : []
            content {
              subject   = mfa_email.value.subject
              html_body = mfa_email.value.html_body
              text_body = mfa_email.value.text_body
            }
          }

          dynamic "no_action_email" {
            for_each = notify_configuration.value.no_action_email != null ? [notify_configuration.value.no_action_email] : []
            content {
              subject   = no_action_email.value.subject
              html_body = no_action_email.value.html_body
              text_body = no_action_email.value.text_body
            }
          }
        }
      }
    }
  }

  dynamic "compromised_credentials_risk_configuration" {
    for_each = var.spec.risk_configuration.compromised_credentials != null ? [var.spec.risk_configuration.compromised_credentials] : []
    content {
      actions {
        event_action = compromised_credentials_risk_configuration.value.event_action
      }

      # Empty means AWS's default (all supported events) -- send absence.
      event_filter = length(compromised_credentials_risk_configuration.value.event_filter) > 0 ? compromised_credentials_risk_configuration.value.event_filter : null
    }
  }

  dynamic "risk_exception_configuration" {
    for_each = var.spec.risk_configuration.risk_exception != null ? [var.spec.risk_configuration.risk_exception] : []
    content {
      blocked_ip_range_list = length(risk_exception_configuration.value.blocked_ip_ranges) > 0 ? risk_exception_configuration.value.blocked_ip_ranges : null
      skipped_ip_range_list = length(risk_exception_configuration.value.skipped_ip_ranges) > 0 ? risk_exception_configuration.value.skipped_ip_ranges : null
    }
  }
}
