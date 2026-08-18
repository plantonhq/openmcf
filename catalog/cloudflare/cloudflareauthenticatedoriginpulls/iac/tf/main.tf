# Cloudflare Authenticated Origin Pulls enablement: the zone-wide toggle plus
# per-hostname certificate associations. Destroy semantics differ per surface
# and neither is a delete: the zone toggle has NO delete at Cloudflare
# (destroy abandons the live value), and an association is removed by a
# revert write (enabled: null) the provider issues from state -- the
# association's cert_id must still be in state for that revert to land.

# The zone-wide toggle. Managed only when the spec sets zone_enabled --
# unset means "leave the zone's toggle alone".
resource "cloudflare_authenticated_origin_pulls_settings" "zone" {
  count = var.spec.zone_enabled != null ? 1 : 0

  zone_id = var.spec.zone_id
  enabled = var.spec.zone_enabled
}

# One association resource per hostname row. The provider requires the config
# list to hold exactly one element per resource (it hard-fails otherwise), so
# each row fans out to its own resource. An omitted enabled is sent as true:
# Cloudflare treats null as "void the association", and a declared row is
# meant to exist.
resource "cloudflare_authenticated_origin_pulls" "association" {
  for_each = local.hostname_associations

  zone_id = var.spec.zone_id

  config = [{
    hostname = each.value.hostname
    cert_id  = each.value.certificate_id != "" ? each.value.certificate_id : null
    enabled  = coalesce(each.value.enabled, true)
  }]
}
