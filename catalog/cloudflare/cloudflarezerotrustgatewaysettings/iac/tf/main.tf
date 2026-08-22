# Secure Web Gateway configuration: three folded Cloudflare surfaces with
# different lifecycles.
#
#   - cloudflare_zero_trust_gateway_settings: the account configuration
#     SINGLETON. Create and update are the same PUT; DESTROY IS A NO-OP that
#     abandons the live configuration exactly as last applied. An unset spec
#     sub-object is never sent, leaving the live value (dashboard-set or
#     default) untouched.
#   - cloudflare_zero_trust_gateway_logging: the logging SINGLETON, same
#     lifecycle. The module sends the COMPLETE logging tree when the spec
#     declares it -- omitted logging fields drift at Cloudflare (the
#     provider's own tests accept non-empty plans on partial sends), so
#     explicit emission is what keeps apply idempotent.
#   - cloudflare_zero_trust_gateway_pacfile: one resource per spec row, each
#     with a real create/update/delete lifecycle.
resource "cloudflare_zero_trust_gateway_settings" "main" {
  count = try(var.spec.settings, null) != null ? 1 : 0

  account_id = var.spec.account_id

  settings = {
    activity_log = try(var.spec.settings.activity_log, null) != null ? {
      enabled = try(var.spec.settings.activity_log.enabled, null)
    } : null

    antivirus = try(var.spec.settings.antivirus, null) != null ? {
      enabled_download_phase = try(var.spec.settings.antivirus.enabled_download_phase, null)
      enabled_upload_phase   = try(var.spec.settings.antivirus.enabled_upload_phase, null)
      fail_closed            = try(var.spec.settings.antivirus.fail_closed, null)
      notification_settings = try(var.spec.settings.antivirus.notification_settings, null) != null ? {
        enabled         = try(var.spec.settings.antivirus.notification_settings.enabled, null)
        msg             = try(var.spec.settings.antivirus.notification_settings.msg, "") != "" ? var.spec.settings.antivirus.notification_settings.msg : null
        support_url     = try(var.spec.settings.antivirus.notification_settings.support_url, "") != "" ? var.spec.settings.antivirus.notification_settings.support_url : null
        include_context = try(var.spec.settings.antivirus.notification_settings.include_context, null)
      } : null
    } : null

    block_page = try(var.spec.settings.block_page, null) != null ? {
      enabled          = try(var.spec.settings.block_page.enabled, null)
      mode             = try(var.spec.settings.block_page.mode, "") != "" ? var.spec.settings.block_page.mode : null
      target_uri       = try(var.spec.settings.block_page.target_uri, "") != "" ? var.spec.settings.block_page.target_uri : null
      include_context  = try(var.spec.settings.block_page.include_context, null)
      name             = try(var.spec.settings.block_page.name, "") != "" ? var.spec.settings.block_page.name : null
      header_text      = try(var.spec.settings.block_page.header_text, "") != "" ? var.spec.settings.block_page.header_text : null
      footer_text      = try(var.spec.settings.block_page.footer_text, "") != "" ? var.spec.settings.block_page.footer_text : null
      suppress_footer  = try(var.spec.settings.block_page.suppress_footer, null)
      background_color = try(var.spec.settings.block_page.background_color, "") != "" ? var.spec.settings.block_page.background_color : null
      logo_path        = try(var.spec.settings.block_page.logo_path, "") != "" ? var.spec.settings.block_page.logo_path : null
      mailto_address   = try(var.spec.settings.block_page.mailto_address, "") != "" ? var.spec.settings.block_page.mailto_address : null
      mailto_subject   = try(var.spec.settings.block_page.mailto_subject, "") != "" ? var.spec.settings.block_page.mailto_subject : null
    } : null

    body_scanning = try(var.spec.settings.body_scanning, null) != null ? {
      inspection_mode = try(var.spec.settings.body_scanning.inspection_mode, "") != "" ? var.spec.settings.body_scanning.inspection_mode : null
    } : null

    browser_isolation = try(var.spec.settings.browser_isolation, null) != null ? {
      non_identity_enabled          = try(var.spec.settings.browser_isolation.non_identity_enabled, null)
      url_browser_isolation_enabled = try(var.spec.settings.browser_isolation.url_browser_isolation_enabled, null)
    } : null

    certificate = try(var.spec.settings.certificate, null) != null ? {
      id = var.spec.settings.certificate.id
    } : null

    extended_email_matching = try(var.spec.settings.extended_email_matching, null) != null ? {
      enabled = try(var.spec.settings.extended_email_matching.enabled, null)
    } : null

    fips = try(var.spec.settings.fips, null) != null ? {
      tls = try(var.spec.settings.fips.tls, null)
    } : null

    host_selector = try(var.spec.settings.host_selector, null) != null ? {
      enabled = try(var.spec.settings.host_selector.enabled, null)
    } : null

    inspection = try(var.spec.settings.inspection, null) != null ? {
      mode = try(var.spec.settings.inspection.mode, "") != "" ? var.spec.settings.inspection.mode : null
    } : null

    max_ttl_secs = try(var.spec.settings.max_ttl_secs, null)

    protocol_detection = try(var.spec.settings.protocol_detection, null) != null ? {
      enabled = try(var.spec.settings.protocol_detection.enabled, null)
    } : null

    sandbox = try(var.spec.settings.sandbox, null) != null ? {
      enabled         = try(var.spec.settings.sandbox.enabled, null)
      fallback_action = try(var.spec.settings.sandbox.fallback_action, "") != "" ? var.spec.settings.sandbox.fallback_action : null
    } : null

    tls_decrypt = try(var.spec.settings.tls_decrypt, null) != null ? {
      enabled = try(var.spec.settings.tls_decrypt.enabled, null)
    } : null
  }
}

# The logging singleton: the COMPLETE tree is always sent (unset spec fields
# become false, Cloudflare's own default) because partial sends drift.
resource "cloudflare_zero_trust_gateway_logging" "main" {
  count = try(var.spec.logging, null) != null ? 1 : 0

  account_id = var.spec.account_id
  redact_pii = try(var.spec.logging.redact_pii, null) != null ? var.spec.logging.redact_pii : false

  settings_by_rule_type = {
    dns = {
      log_all    = try(var.spec.logging.settings_by_rule_type.dns.log_all, null) != null ? var.spec.logging.settings_by_rule_type.dns.log_all : false
      log_blocks = try(var.spec.logging.settings_by_rule_type.dns.log_blocks, null) != null ? var.spec.logging.settings_by_rule_type.dns.log_blocks : false
    }
    http = {
      log_all    = try(var.spec.logging.settings_by_rule_type.http.log_all, null) != null ? var.spec.logging.settings_by_rule_type.http.log_all : false
      log_blocks = try(var.spec.logging.settings_by_rule_type.http.log_blocks, null) != null ? var.spec.logging.settings_by_rule_type.http.log_blocks : false
    }
    l4 = {
      log_all    = try(var.spec.logging.settings_by_rule_type.l4.log_all, null) != null ? var.spec.logging.settings_by_rule_type.l4.log_all : false
      log_blocks = try(var.spec.logging.settings_by_rule_type.l4.log_blocks, null) != null ? var.spec.logging.settings_by_rule_type.l4.log_blocks : false
    }
  }
}

# PAC files: one resource per row, keyed by name so a row edit replaces only
# its own file. slug forces replacement at the provider (it is baked into
# the file's public URL).
resource "cloudflare_zero_trust_gateway_pacfile" "main" {
  for_each = { for pac_file in try(var.spec.pac_files, []) : pac_file.name => pac_file }

  account_id  = var.spec.account_id
  name        = each.value.name
  contents    = each.value.contents
  slug        = try(each.value.slug, "") != "" ? each.value.slug : null
  description = try(each.value.description, "") != "" ? each.value.description : null
}
