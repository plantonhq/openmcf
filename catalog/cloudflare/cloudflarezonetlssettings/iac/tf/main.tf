# Each TLS setting is its own zone-scoped API object, emitted only when the
# manifest manages it. Universal SSL, Total TLS, auto origin key exchange, and
# CA hostname associations have NO delete at Cloudflare -- destroy drops state
# and abandons the live values. Per-hostname overrides and origin TLS
# compliance modes have real deletes.

# Universal SSL issuance. Disabling can make proxied hostnames unreachable
# over HTTPS unless other certificates cover them -- the spec carries the
# warning. NO delete at Cloudflare.
resource "cloudflare_universal_ssl_setting" "main" {
  count = var.spec.universal_ssl_enabled != null ? 1 : 0

  zone_id = local.zone_id
  enabled = var.spec.universal_ssl_enabled
}

# Total TLS. NO delete at Cloudflare. The certificates' validity period is
# fixed by Cloudflare (90 days) and deliberately not modeled.
resource "cloudflare_total_tls" "main" {
  count = var.spec.total_tls != null ? 1 : 0

  zone_id               = local.zone_id
  enabled               = var.spec.total_tls.enabled
  certificate_authority = var.spec.total_tls.certificate_authority
}

# Automatic origin TLS key exchange. NO delete at Cloudflare.
resource "cloudflare_zone_auto_origin_tls_kex" "main" {
  count = var.spec.auto_origin_tls_kex != null ? 1 : 0

  zone_id = local.zone_id
  enabled = var.spec.auto_origin_tls_kex
}

# Origin TLS compliance modes. Real delete: destroying the resource clears the
# compliance requirement. The module never sends an empty list.
resource "cloudflare_origin_tls_compliance_modes" "main" {
  count = length(var.spec.origin_tls_compliance_modes) > 0 ? 1 : 0

  zone_id = local.zone_id
  value   = var.spec.origin_tls_compliance_modes
}

# Per-hostname TLS overrides: one API object per (setting_id, hostname). Real
# delete: destroy removes the overrides and the hostnames fall back to the
# zone-wide settings.
resource "cloudflare_hostname_tls_setting" "min_tls_version" {
  for_each = local.hostname_min_tls_version

  zone_id    = local.zone_id
  setting_id = "min_tls_version"
  hostname   = each.key
  value      = each.value
}

resource "cloudflare_hostname_tls_setting" "http2" {
  for_each = local.hostname_http2

  zone_id    = local.zone_id
  setting_id = "http2"
  hostname   = each.key
  value      = each.value ? "on" : "off"
}

resource "cloudflare_hostname_tls_setting" "ciphers" {
  for_each = local.hostname_ciphers

  zone_id    = local.zone_id
  setting_id = "ciphers"
  hostname   = each.key
  value      = each.value
}

# Certificate-authority hostname associations: the "managed" key is the zone's
# managed-CA list (the row without an mTLS certificate id); every other key is
# an mTLS certificate's own hostname list. NO delete at Cloudflare.
resource "cloudflare_certificate_authorities_hostname_associations" "main" {
  for_each = local.ca_hostname_associations

  zone_id             = local.zone_id
  hostnames           = each.value.hostnames
  mtls_certificate_id = each.value.mtls_certificate_id == "" ? null : each.value.mtls_certificate_id
}
