# Cloudflare Gateway policy: a filter over employee traffic (DNS, HTTP, or
# network) plus the action taken when it matches.
#
# Two provider truths this module encodes:
#   - `enabled` DEFAULTS TO FALSE at Cloudflare; the spec models it explicitly
#     so a policy is never silently inert.
#   - `rule_settings` is ALWAYS sent (empty object when the spec configures
#     nothing) -- the provider's own fixtures do the same to prevent API drift.
#
# KNOWN UPSTREAM DEFECT at v5.23.0 (unfixed through v5.24.0 and provider
# main; upstream issue #7106): the resource's Computed attributes ship no
# UseStateForUnknown plan modifiers, so refresh-inclusive plans NEVER
# converge -- measured live 2026-08-26 on every configuration shape
# (rule_settings absent, empty object, and populated all re-plan an
# in-place update forever; the visible driver is
# rule_settings.ignore_cname_category_matches, which Cloudflare's API never
# echoes back). Module-side remedies were measured non-fixing: an explicit
# send is accepted-and-dropped by the API, and lifecycle ignore_changes
# cannot suppress a provider-planned unknown. Applies succeed and the
# repeated update is a no-op write; the perpetual pending plan is upstream's
# to fix. The earlier-known add_headers/override_ips first-apply drift is
# the same family.
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
