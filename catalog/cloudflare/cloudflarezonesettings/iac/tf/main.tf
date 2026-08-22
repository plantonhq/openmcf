# The zone-settings fan-out: one cloudflare_zone_setting per MANAGED setting.
# Unset spec fields never enter local.settings, so the module never sends
# defaults -- zone settings have no delete at Cloudflare, and anything sent
# once is owned until reverted explicitly. Destroy drops state and abandons
# the live values (the provider's own destroy is a no-op for this resource).
resource "cloudflare_zone_setting" "main" {
  for_each = local.settings

  zone_id    = local.zone_id
  setting_id = each.key
  value      = each.value
}

# Managed transforms: one zone-wide object carrying both header lists. The API
# treats a transform missing from the list as disabled, so both lists ride
# together whenever either is managed. Real delete: destroy disables all
# managed transforms.
resource "cloudflare_managed_transforms" "main" {
  count = length(var.spec.managed_request_headers) > 0 || length(var.spec.managed_response_headers) > 0 ? 1 : 0

  zone_id = local.zone_id
  managed_request_headers = [
    for header in var.spec.managed_request_headers : {
      id      = header.id
      enabled = header.enabled
    }
  ]
  managed_response_headers = [
    for header in var.spec.managed_response_headers : {
      id      = header.id
      enabled = header.enabled
    }
  ]
}

# URL normalization: a zone singleton with a real delete (destroy resets
# normalization to Cloudflare defaults).
resource "cloudflare_url_normalization_settings" "main" {
  count = var.spec.url_normalization != null ? 1 : 0

  zone_id = local.zone_id
  scope   = var.spec.url_normalization.scope
  type    = var.spec.url_normalization.type
}

# Origin cloud regions: one API object per origin IP (the IP is the row's
# identity and the for_each key, so reordering rows never churns resources).
# Real delete: destroy removes the region hints.
resource "cloudflare_origin_cloud_region" "main" {
  for_each = { for region in var.spec.origin_cloud_regions : region.origin_ip => region }

  zone_id   = local.zone_id
  origin_ip = each.value.origin_ip
  region    = each.value.region
  vendor    = each.value.vendor
}

# Waiting-room crawler bypass: a zone singleton with NO delete at Cloudflare --
# destroy abandons the last-applied value.
resource "cloudflare_waiting_room_settings" "main" {
  count = var.spec.waiting_room_crawler_bypass != null ? 1 : 0

  zone_id                      = local.zone_id
  search_engine_crawler_bypass = var.spec.waiting_room_crawler_bypass
}
