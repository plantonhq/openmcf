# Cloudflare Gateway policy: a filter over employee traffic (DNS, HTTP, or
# network) plus the action taken when it matches.
#
# Two provider truths this module encodes:
#   - `enabled` DEFAULTS TO FALSE at Cloudflare; the spec models it explicitly
#     so a policy is never silently inert.
#   - `rule_settings` is ALWAYS sent (empty object when the spec configures
#     nothing) -- the provider's own fixtures do the same to prevent API drift.
#
# KNOWN UPSTREAM DRIFT at v5.23.0: policies carrying add_headers or
# override_ips show computed-field drift even on first apply (the provider's
# own migration tests expect a non-empty plan for them).
resource "cloudflare_zero_trust_gateway_policy" "main" {
  account_id = var.spec.account_id

  name        = var.spec.name
  action      = var.spec.action
  description = var.spec.description != "" ? var.spec.description : null

  filters = local.filters

  # Explicit pass-through: unset inherits Cloudflare's default of FALSE.
  enabled = var.spec.enabled

  # Unset lets Cloudflare assign a precedence; lower runs earlier.
  precedence = var.spec.precedence

  # Wirefilter expressions. The API reformats them before storing; if the plan
  # shows drift, adopt the API-returned form.
  traffic        = var.spec.traffic != "" ? var.spec.traffic : null
  identity       = var.spec.identity != "" ? var.spec.identity : null
  device_posture = var.spec.device_posture != "" ? var.spec.device_posture : null

  expiration = local.expiration
  schedule   = local.schedule

  rule_settings = local.rule_settings
}
