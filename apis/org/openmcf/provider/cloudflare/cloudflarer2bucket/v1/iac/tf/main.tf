# Cloudflare R2 Bucket
# R2 is Cloudflare's S3-compatible object storage with zero egress fees
resource "cloudflare_r2_bucket" "main" {
  account_id = local.account_id
  name       = local.bucket_name
  location   = local.location
}

# Note: Public access (r2.dev subdomain) is not yet supported in the Terraform provider
# It must be enabled manually via the Cloudflare Dashboard or API
# See: https://developers.cloudflare.com/r2/buckets/public-buckets/

# Note: R2 does not support object versioning
# The versioning_enabled field is ignored

# Custom domain for the R2 bucket
# Created only when custom_domain.enabled is true
resource "cloudflare_r2_custom_domain" "main" {
  count = local.custom_domain_enabled ? 1 : 0

  account_id  = local.account_id
  bucket_name = cloudflare_r2_bucket.main.name
  zone_id     = local.custom_domain_zone_id
  hostname    = local.custom_domain_name
}
