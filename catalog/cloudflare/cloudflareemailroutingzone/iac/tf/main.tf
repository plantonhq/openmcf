# Enabling Email Routing on the zone. Creating this resource turns Email Routing
# on and provisions the zone's required MX/SPF/DKIM records automatically.
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
