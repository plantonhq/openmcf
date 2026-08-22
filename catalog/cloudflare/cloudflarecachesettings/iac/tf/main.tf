# Each cache setting is its own zone-scoped API object, emitted only when the
# manifest manages it. Most have NO delete at Cloudflare -- destroy drops state
# and abandons the live value (smart tiered cache and cache variants are the
# real-delete exceptions). Argo Smart Routing is paid and KEEPS BILLING after
# destroy; apply false first when retiring it.

# Smart Tiered Cache (the dashboard's Tiered Cache toggle). Real delete:
# destroy disables it.
resource "cloudflare_tiered_cache" "main" {
  count = var.spec.smart_tiered_cache != null ? 1 : 0

  zone_id = local.zone_id
  value   = var.spec.smart_tiered_cache ? "on" : "off"
}

# Generic tiered caching. NO delete at Cloudflare.
resource "cloudflare_argo_tiered_caching" "main" {
  count = var.spec.tiered_caching != null ? 1 : 0

  zone_id = local.zone_id
  value   = var.spec.tiered_caching ? "on" : "off"
}

# Cache Reserve (paid; storage keeps billing while on). NO delete at
# Cloudflare -- turn it off explicitly before retiring.
resource "cloudflare_zone_cache_reserve" "main" {
  count = var.spec.cache_reserve != null ? 1 : 0

  zone_id = local.zone_id
  value   = var.spec.cache_reserve ? "on" : "off"
}

# Regional Tiered Cache. NO delete at Cloudflare.
resource "cloudflare_regional_tiered_cache" "main" {
  count = var.spec.regional_tiered_cache != null ? 1 : 0

  zone_id = local.zone_id
  value   = var.spec.regional_tiered_cache ? "on" : "off"
}

# Argo Smart Routing. PAID, NO delete at Cloudflare -- destroying with the
# value still on KEEPS BILLING.
resource "cloudflare_argo_smart_routing" "main" {
  count = var.spec.argo_smart_routing != null ? 1 : 0

  zone_id = local.zone_id
  value   = var.spec.argo_smart_routing ? "on" : "off"
}

# Cache variants. Real delete: destroy resets variants to defaults.
resource "cloudflare_zone_cache_variants" "main" {
  count = var.spec.cache_variants != null ? 1 : 0

  zone_id = local.zone_id
  value   = local.cache_variants_value
}
