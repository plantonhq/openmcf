# main.tf

# Create the Cloudflare DNS Zone
resource "cloudflare_zone" "main" {
  account_id = var.spec.account_id
  zone       = var.spec.zone_name
  plan       = local.plan_id
  paused     = var.spec.paused

  # Jump start automatically scans and imports existing DNS records
  # Set to false for clean zone creation
  jump_start = false

  # Type is always "full" (standard nameserver delegation)
  # "partial" (CNAME setup) is Business/Enterprise only and not in the 80/20
  type = "full"
}

# Configure zone settings including default_proxied
resource "cloudflare_zone_settings_override" "main" {
  zone_id = cloudflare_zone.main.id

  settings {
    # Set default proxied behavior for new DNS records
    # This is configured as a zone setting, not at zone creation
    always_online = "on"
    
    # Note: default_proxied is not a direct zone setting in the Cloudflare API
    # It's a UI/dashboard feature. For Infrastructure-as-Code, you control
    # the proxied state explicitly when creating each DNS record.
  }
}

# Create DNS records within the zone
resource "cloudflare_record" "records" {
  for_each = { for idx, record in var.spec.records : "${record.name}-${record.type}" => record }

  zone_id = cloudflare_zone.main.id
  name    = each.value.name
  type    = each.value.type
  content = each.value.value
  ttl     = each.value.ttl

  # proxied is only applicable to A, AAAA, and CNAME records
  proxied = contains(["A", "AAAA", "CNAME"], each.value.type) ? each.value.proxied : false

  # priority is only used for MX and SRV records
  priority = contains(["MX", "SRV"], each.value.type) ? each.value.priority : null

  # comment for the DNS record
  comment = each.value.comment != "" ? each.value.comment : null

  depends_on = [cloudflare_zone.main]
}

