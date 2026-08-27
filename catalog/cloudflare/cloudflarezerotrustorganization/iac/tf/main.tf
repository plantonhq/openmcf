# Zero Trust organization: the Access login experience (team domain, login
# design, session defaults, MFA policy) plus the folded Access service-key
# rotation cadence.
#
# Cloudflare has no create call for an organization -- both create and update
# are the same PUT (an upsert of the singleton the account or zone already
# carries), and DESTROY IS A NO-OP: removing this resource abandons the live
# configuration exactly as last applied. Unset spec fields are never sent,
# leaving the live value untouched.
resource "cloudflare_zero_trust_organization" "main" {
  # Exactly one scope is set (spec validation enforces it).
  account_id = var.spec.account_id != "" ? var.spec.account_id : null
  zone_id    = try(var.spec.zone_id, "") != "" ? var.spec.zone_id : null

  # Always sent: Cloudflare rejects any organization write without
  # auth_domain (API error 11004 invalid_auth_domain, live-measured) --
  # the spec requires it.
  auth_domain                        = var.spec.auth_domain
  name                               = try(var.spec.name, "") != "" ? var.spec.name : null
  session_duration                   = try(var.spec.session_duration, "") != "" ? var.spec.session_duration : null
  warp_auth_session_duration         = try(var.spec.warp_auth_session_duration, "") != "" ? var.spec.warp_auth_session_duration : null
  user_seat_expiration_inactive_time = try(var.spec.user_seat_expiration_inactive_time, "") != "" ? var.spec.user_seat_expiration_inactive_time : null

  deny_unmatched_requests                     = try(var.spec.deny_unmatched_requests, null)
  deny_unmatched_requests_exempted_zone_names = length(try(var.spec.deny_unmatched_requests_exempted_zone_names, [])) > 0 ? var.spec.deny_unmatched_requests_exempted_zone_names : null

  allow_authenticate_via_warp = try(var.spec.allow_authenticate_via_warp, null)
  auto_redirect_to_identity   = try(var.spec.auto_redirect_to_identity, null)
  mfa_required_for_all_apps   = try(var.spec.mfa_required_for_all_apps, null)
  is_ui_read_only             = try(var.spec.is_ui_read_only, null)
  ui_read_only_toggle_reason  = try(var.spec.ui_read_only_toggle_reason, "") != "" ? var.spec.ui_read_only_toggle_reason : null

  custom_pages = try(var.spec.custom_pages, null) != null ? {
    forbidden       = try(var.spec.custom_pages.forbidden, "") != "" ? var.spec.custom_pages.forbidden : null
    identity_denied = try(var.spec.custom_pages.identity_denied, "") != "" ? var.spec.custom_pages.identity_denied : null
  } : null

  login_design = try(var.spec.login_design, null) != null ? {
    background_color = try(var.spec.login_design.background_color, "") != "" ? var.spec.login_design.background_color : null
    text_color       = try(var.spec.login_design.text_color, "") != "" ? var.spec.login_design.text_color : null
    logo_path        = try(var.spec.login_design.logo_path, "") != "" ? var.spec.login_design.logo_path : null
    header_text      = try(var.spec.login_design.header_text, "") != "" ? var.spec.login_design.header_text : null
    footer_text      = try(var.spec.login_design.footer_text, "") != "" ? var.spec.login_design.footer_text : null
  } : null

  mfa_config = try(var.spec.mfa_config, null) != null ? {
    allowed_authenticators        = length(try(var.spec.mfa_config.allowed_authenticators, [])) > 0 ? var.spec.mfa_config.allowed_authenticators : null
    session_duration              = try(var.spec.mfa_config.session_duration, "") != "" ? var.spec.mfa_config.session_duration : null
    amr_matching_session_duration = try(var.spec.mfa_config.amr_matching_session_duration, "") != "" ? var.spec.mfa_config.amr_matching_session_duration : null
    required_aaguids              = try(var.spec.mfa_config.required_aaguids, "") != "" ? var.spec.mfa_config.required_aaguids : null
  } : null

  mfa_ssh_piv_key_requirements = try(var.spec.mfa_ssh_piv_key_requirements, null) != null ? {
    pin_policy          = try(var.spec.mfa_ssh_piv_key_requirements.pin_policy, "") != "" ? var.spec.mfa_ssh_piv_key_requirements.pin_policy : null
    touch_policy        = try(var.spec.mfa_ssh_piv_key_requirements.touch_policy, "") != "" ? var.spec.mfa_ssh_piv_key_requirements.touch_policy : null
    ssh_key_types       = length(try(var.spec.mfa_ssh_piv_key_requirements.ssh_key_type, [])) > 0 ? var.spec.mfa_ssh_piv_key_requirements.ssh_key_type : null
    ssh_key_sizes       = length(try(var.spec.mfa_ssh_piv_key_requirements.ssh_key_size, [])) > 0 ? var.spec.mfa_ssh_piv_key_requirements.ssh_key_size : null
    require_fips_device = try(var.spec.mfa_ssh_piv_key_requirements.require_fips_device, null)
  } : null
}

# The folded Access service-key rotation cadence: its own Cloudflare surface
# with the same singleton/upsert/no-op-destroy lifecycle, deployed only when
# the spec declares a cadence (account scope only -- spec validation enforces
# that). Ordered after the organization so a fresh Zero Trust account applies
# the org first.
resource "cloudflare_zero_trust_access_key_configuration" "main" {
  count = try(var.spec.key_rotation_interval_days, null) != null ? 1 : 0

  account_id                 = var.spec.account_id
  key_rotation_interval_days = var.spec.key_rotation_interval_days

  depends_on = [cloudflare_zero_trust_organization.main]
}
