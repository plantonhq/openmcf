# hold.tf

# Zone hold, provisioned only when the spec enables it. While active, Cloudflare
# rejects attempts to create this zone's hostname (and optionally any subdomain)
# as a zone in any other account. An empty hold_after is sent as null: the
# provider treats an empty string as "hold from the current time", which would
# produce a perpetual diff against an unset spec field.
resource "cloudflare_zone_hold" "main" {
  count   = var.spec.hold != null ? (var.spec.hold.enabled ? 1 : 0) : 0
  zone_id = cloudflare_zone.main.id

  include_subdomains = var.spec.hold.include_subdomains
  hold_after         = var.spec.hold.hold_after != "" ? var.spec.hold.hold_after : null
}
