locals {
  # Resource naming
  resource_name = coalesce(try(var.metadata.name, null), "cloudflare-r2-bucket")
  
  # Labels/tags
  labels = merge({
    "name" = local.resource_name
  }, try(var.metadata.labels, {}))

  # Location enum to string mapping
  location_map = {
    0 = "auto"         # auto (Cloudflare chooses optimal region)
    1 = "WNAM"         # Western North America
    2 = "ENAM"         # Eastern North America
    3 = "WEUR"         # Western Europe
    4 = "EEUR"         # Eastern Europe
    5 = "APAC"         # Asia-Pacific
    6 = "OC"           # Oceania
  }
  
  location = lookup(local.location_map, try(var.spec.location, 0), "auto")
  
  # Bucket configuration
  bucket_name = var.spec.bucket_name
  account_id  = var.spec.account_id
  
  # Public access and versioning flags
  public_access       = coalesce(try(var.spec.public_access, null), false)
  versioning_enabled  = coalesce(try(var.spec.versioning_enabled, null), false)

  # Custom domain configuration
  custom_domain_enabled = try(var.spec.custom_domain.enabled, false)
  custom_domain_zone_id = try(var.spec.custom_domain.zone_id.value, "")
  custom_domain_name    = try(var.spec.custom_domain.domain, "")
}

