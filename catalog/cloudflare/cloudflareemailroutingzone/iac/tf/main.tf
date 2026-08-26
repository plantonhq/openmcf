# Enabling Email Routing on the zone. Creating this resource turns Email Routing
# on and provisions the zone's required MX/SPF/DKIM records automatically.
#
# UPSTREAM DEFECT (open at provider v5.23.0 AND v5.24.0, measured live
# 2026-08-26): every create/refresh of this resource fails with "Value
# Conversion Error ... Struct defines fields not found in object:
# support_subaddress" -- the provider's model gained support_subaddress in
# v5.23.0 but its schema never did (terraform-provider-cloudflare issue
# #7301, fix PR #7302). No configuration works around a schema/model
# mismatch; `plan` stays green because plan never runs the conversion. When
# a release carries the fix, this module works unchanged.
resource "cloudflare_email_routing_settings" "main" {
  zone_id = var.spec.zone_id
}

# The single per-zone catch-all rule (folded), created only when configured.
# NOTE: the provider's Delete for this resource is a genuine no-op (no API
# call) -- destroying it drops it from state and the zone keeps whatever
# catch-all configuration it last had. The zone-level destroy (disabling Email
# Routing) is what actually retires the behavior.
resource "cloudflare_email_routing_catch_all" "main" {
  count   = local.catch_all != null ? 1 : 0
  zone_id = var.spec.zone_id
  enabled = try(local.catch_all.enabled, false)
  name    = try(local.catch_all.name, "") != "" ? local.catch_all.name : null

  # "all" is the only matcher type Cloudflare permits on the catch-all, so the
  # module supplies it rather than modeling it in the spec.
  matchers = [{
    type = "all"
  }]

  actions = [
    for a in local.catch_all_actions : {
      type  = a.type
      value = length(a.value) > 0 ? a.value : null
    }
  ]

  depends_on = [cloudflare_email_routing_settings.main]
}

# Optionally lock the Email Routing DNS records so they cannot be modified
# out-of-band. dns_name targets a subdomain of the zone; empty routes the apex.
resource "cloudflare_email_routing_dns" "main" {
  count   = try(var.spec.lock_dns_records, false) ? 1 : 0
  zone_id = var.spec.zone_id
  name    = try(var.spec.dns_name, "") != "" ? var.spec.dns_name : null

  depends_on = [cloudflare_email_routing_settings.main]
}
