# Enable the Identity Toolkit API so a fresh project can be initialized.
# disable_on_destroy is false: tearing down this resource must never disable
# authentication for everything else in the project.
resource "google_project_service" "identitytoolkit_api" {
  project = local.project_id
  service = "identitytoolkit.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# Resolves the ambient project ONLY for the adoption import ID below, and
# only when the spec omits project_id — plans that name their project stay
# credential-free (the count-gated client-config pattern).
data "google_client_config" "this" {
  count = var.spec.adopt_existing && var.spec.project_id == "" ? 1 : 0
}

# Adoption: initialization is one-way and ONCE-ONLY — GCP rejects a second
# initializeAuth with 400 "Identity Platform has already been enabled for
# this project" (live-verified). When the spec arms adopt_existing, the
# deterministic singleton (projects/{project}/config) is imported and the
# spec applies as an update; the block is inert once the resource is in
# state, so re-applies stay clean.
import {
  for_each = var.spec.adopt_existing ? toset(["adopt"]) : toset([])
  to       = google_identity_platform_config.this
  id       = "projects/${var.spec.project_id != "" ? var.spec.project_id : data.google_client_config.this[0].project}/config"
}

# The project's Identity Platform configuration — a ONE-WAY project
# singleton. The provider's create is a bare initializeAuth call
# (permanently enabling Identity Platform on the project, billing required)
# followed by an update that applies every setting, and its delete abandons
# the configuration in place — GCP has no de-initialize; a create on an
# already-initialized project hard-fails (arm spec.adopt_existing there).
# That is why this resource carries no deletion_policy: the spec's
# deletion_policy governs only the composed IdP configs below.
#
# Every sign-in arm's `enabled` is sent EXPLICITLY whenever its object is
# present: the fields drive live authentication surfaces, and a spec
# transition true -> false must reach the API rather than being omitted.
resource "google_identity_platform_config" "this" {
  project = local.project_id

  authorized_domains         = length(var.spec.authorized_domains) > 0 ? var.spec.authorized_domains : null
  autodelete_anonymous_users = var.spec.autodelete_anonymous_users ? true : null

  # The API always materializes the email and phone_number policies in its
  # read-back (enabled=false when never configured), so BOTH blocks render
  # explicitly whenever sign_in is present — omitting one leaves a perpetual
  # block-removal diff on every re-plan (idempotency-gate caught,
  # live-verified). The anonymous arm is NOT echoed when unset and renders
  # only when the spec sets it.
  dynamic "sign_in" {
    for_each = var.spec.sign_in != null ? [var.spec.sign_in] : []
    content {
      allow_duplicate_emails = sign_in.value.allow_duplicate_emails

      email {
        enabled           = sign_in.value.email != null ? sign_in.value.email.enabled : false
        password_required = sign_in.value.email != null ? sign_in.value.email.password_required : null
      }
      phone_number {
        enabled            = sign_in.value.phone_number != null ? sign_in.value.phone_number.enabled : false
        test_phone_numbers = sign_in.value.phone_number != null && length(coalesce(sign_in.value.phone_number.test_phone_numbers, {})) > 0 ? sign_in.value.phone_number.test_phone_numbers : null
      }
      dynamic "anonymous" {
        for_each = sign_in.value.anonymous != null ? [sign_in.value.anonymous] : []
        content {
          enabled = anonymous.value.enabled
        }
      }
    }
  }

  dynamic "mfa" {
    for_each = var.spec.mfa != null ? [var.spec.mfa] : []
    content {
      state             = mfa.value.state != "" ? mfa.value.state : null
      enabled_providers = length(mfa.value.enabled_providers) > 0 ? mfa.value.enabled_providers : null

      dynamic "provider_configs" {
        for_each = mfa.value.provider_configs
        content {
          state = provider_configs.value.state != "" ? provider_configs.value.state : null

          dynamic "totp_provider_config" {
            for_each = provider_configs.value.totp_provider_config != null ? [provider_configs.value.totp_provider_config] : []
            content {
              adjacent_intervals = totp_provider_config.value.adjacent_intervals
            }
          }
        }
      }
    }
  }

  dynamic "blocking_functions" {
    for_each = var.spec.blocking_functions != null ? [var.spec.blocking_functions] : []
    content {
      dynamic "triggers" {
        for_each = blocking_functions.value.triggers
        content {
          event_type   = triggers.value.event_type
          function_uri = triggers.value.function_uri
        }
      }
      dynamic "forward_inbound_credentials" {
        for_each = blocking_functions.value.forward_inbound_credentials != null ? [blocking_functions.value.forward_inbound_credentials] : []
        content {
          access_token  = forward_inbound_credentials.value.access_token
          id_token      = forward_inbound_credentials.value.id_token
          refresh_token = forward_inbound_credentials.value.refresh_token
        }
      }
    }
  }

  dynamic "quota" {
    for_each = var.spec.sign_up_quota != null ? [var.spec.sign_up_quota] : []
    content {
      sign_up_quota_config {
        quota          = quota.value.quota
        quota_duration = quota.value.quota_duration
        start_time     = quota.value.start_time
      }
    }
  }

  dynamic "sms_region_config" {
    for_each = var.spec.sms_region_config != null ? [var.spec.sms_region_config] : []
    content {
      dynamic "allow_by_default" {
        for_each = sms_region_config.value.allow_by_default != null ? [sms_region_config.value.allow_by_default] : []
        content {
          disallowed_regions = length(allow_by_default.value.disallowed_regions) > 0 ? allow_by_default.value.disallowed_regions : null
        }
      }
      dynamic "allowlist_only" {
        for_each = sms_region_config.value.allowlist_only != null ? [sms_region_config.value.allowlist_only] : []
        content {
          allowed_regions = length(allowlist_only.value.allowed_regions) > 0 ? allowlist_only.value.allowed_regions : null
        }
      }
    }
  }

  dynamic "client" {
    for_each = var.spec.client_permissions != null ? [var.spec.client_permissions] : []
    content {
      permissions {
        disabled_user_signup   = client.value.disabled_user_signup
        disabled_user_deletion = client.value.disabled_user_deletion
      }
    }
  }

  dynamic "monitoring" {
    for_each = var.spec.request_logging_enabled != null ? [var.spec.request_logging_enabled] : []
    content {
      request_logging {
        # Explicit send whenever set — true AND false must reach the API.
        enabled = monitoring.value
      }
    }
  }

  dynamic "multi_tenant" {
    for_each = var.spec.multi_tenant != null ? [var.spec.multi_tenant] : []
    content {
      allow_tenants           = multi_tenant.value.allow_tenants
      default_tenant_location = multi_tenant.value.default_tenant_location != "" ? multi_tenant.value.default_tenant_location : null
    }
  }

  depends_on = [google_project_service.identitytoolkit_api]
}

# The composed project-level IdP configs. Each lands only after the project
# is initialized, and each carries the spec's deletion_policy (the config
# itself cannot be deleted). A disabled IdP is sent explicitly — coalesce
# defaults an unset `enabled` to true, the posture users expect when adding
# a provider.
resource "google_identity_platform_default_supported_idp_config" "this" {
  for_each = local.default_supported_idps

  project = local.project_id

  idp_id        = each.value.idp_id
  client_id     = each.value.client_id
  client_secret = each.value.client_secret
  enabled       = coalesce(each.value.enabled, true)

  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  depends_on = [google_identity_platform_config.this]
}

resource "google_identity_platform_oauth_idp_config" "this" {
  for_each = local.oauth_idp_configs

  project = local.project_id

  name          = each.value.name
  display_name  = each.value.display_name != "" ? each.value.display_name : null
  issuer        = each.value.issuer
  client_id     = each.value.client_id
  client_secret = each.value.client_secret != "" ? each.value.client_secret : null
  enabled       = coalesce(each.value.enabled, true)

  dynamic "response_type" {
    for_each = each.value.response_type != null ? [each.value.response_type] : []
    content {
      code     = response_type.value.code
      id_token = response_type.value.id_token
    }
  }

  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  depends_on = [google_identity_platform_config.this]
}

resource "google_identity_platform_inbound_saml_config" "this" {
  for_each = local.inbound_saml_configs

  project = local.project_id

  name         = each.value.name
  display_name = each.value.display_name
  enabled      = coalesce(each.value.enabled, true)

  idp_config {
    idp_entity_id = each.value.idp_config.idp_entity_id
    sso_url       = each.value.idp_config.sso_url
    sign_request  = each.value.idp_config.sign_request

    dynamic "idp_certificates" {
      for_each = each.value.idp_config.idp_certificates
      content {
        x509_certificate = idp_certificates.value.x509_certificate != "" ? idp_certificates.value.x509_certificate : null
      }
    }
  }

  dynamic "sp_config" {
    for_each = each.value.sp_config != null ? [each.value.sp_config] : []
    content {
      callback_uri = sp_config.value.callback_uri != "" ? sp_config.value.callback_uri : null
      sp_entity_id = sp_config.value.sp_entity_id != "" ? sp_config.value.sp_entity_id : null
    }
  }

  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  depends_on = [google_identity_platform_config.this]
}
