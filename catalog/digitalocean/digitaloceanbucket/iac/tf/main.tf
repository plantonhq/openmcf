# DigitalOcean Spaces Bucket
# https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs/resources/spaces_bucket

resource "digitalocean_spaces_bucket" "main" {
  name   = var.spec.bucket_name
  region = var.spec.region

  # Access control (private or public-read)
  acl = local.acl

  # Versioning configuration. Once enabled, Spaces versioning cannot be
  # disabled -- only suspended -- so flipping versioning_enabled back to
  # false suspends versioning rather than removing it.
  dynamic "versioning" {
    for_each = var.spec.versioning_enabled ? [1] : []

    content {
      enabled = true
    }
  }

  # Force destroy for development/testing
  # WARNING: If true, bucket will be deleted even if it contains objects
  force_destroy = false
}

# Optional: Create Spaces bucket CORS configuration
# Uncomment if CORS is needed for web applications
# resource "digitalocean_spaces_bucket_cors_configuration" "main" {
#   bucket = digitalocean_spaces_bucket.main.id
#   region = var.spec.region
#
#   cors_rule {
#     allowed_headers = ["*"]
#     allowed_methods = ["GET", "HEAD"]
#     allowed_origins = ["*"]
#     max_age_seconds = 3000
#   }
# }

