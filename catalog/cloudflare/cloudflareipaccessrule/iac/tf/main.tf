# Cloudflare IP Access rule: an allow/block/challenge decision on an IP, IP
# range, ASN, or country, applied account-wide or to a single zone. The spec's
# CEL guarantees exactly one scope is set, so exactly one of account_id/zone_id
# renders here (the provider would silently prefer account_id if both arrived).
#
# API limitation taught by the provider's own tests: only mode and notes update
# in place. A configuration (target/value) change plans as an in-place update
# but does NOT stick at the API -- change what a rule MATCHES by recreating it.
resource "cloudflare_access_rule" "main" {
  account_id = var.spec.account_id != "" ? var.spec.account_id : null
  zone_id    = var.spec.zone_id != "" ? var.spec.zone_id : null

  mode = var.spec.mode

  configuration = {
    target = var.spec.configuration.target
    value  = var.spec.configuration.value
  }

  notes = var.spec.notes != "" ? var.spec.notes : null
}
