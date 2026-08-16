resource "digitalocean_spaces_bucket" "main" {
  name = var.spec.bucket_name

  # When unset, the provider applies its own default region (nyc3);
  # changing the region replaces the bucket (provider ForceNew).
  region = var.spec.region != "" ? var.spec.region : null

  acl = local.acl

  # Versioning: once enabled it can never be removed — flipping
  # versioning_enabled back to false suspends it (the provider suppresses
  # the block-removal diff), keeping existing versions.
  dynamic "versioning" {
    for_each = var.spec.versioning_enabled ? [1] : []
    content {
      enabled = true
    }
  }

  dynamic "lifecycle_rule" {
    for_each = var.spec.lifecycle_rules
    content {
      # When omitted, the provider generates a rule id.
      id      = lifecycle_rule.value.id != "" ? lifecycle_rule.value.id : null
      prefix  = lifecycle_rule.value.prefix != "" ? lifecycle_rule.value.prefix : null
      enabled = lifecycle_rule.value.enabled

      abort_incomplete_multipart_upload_days = lifecycle_rule.value.abort_incomplete_multipart_upload_days > 0 ? lifecycle_rule.value.abort_incomplete_multipart_upload_days : null

      dynamic "expiration" {
        for_each = lifecycle_rule.value.expiration != null ? [lifecycle_rule.value.expiration] : []
        content {
          # The spec requires exactly one trigger; only the set one is sent.
          date                         = expiration.value.date != "" ? expiration.value.date : null
          days                         = expiration.value.days > 0 ? expiration.value.days : null
          expired_object_delete_marker = expiration.value.expired_object_delete_marker ? true : null
        }
      }

      dynamic "noncurrent_version_expiration" {
        for_each = lifecycle_rule.value.noncurrent_version_expiration != null ? [lifecycle_rule.value.noncurrent_version_expiration] : []
        content {
          days = noncurrent_version_expiration.value.days
        }
      }
    }
  }

  # DANGER: when true, destroy empties the bucket — every object AND every
  # object version — before deleting it.
  force_destroy = var.spec.force_destroy
}

# --- Per-bucket settings satellites ---------------------------------------
# Separate provider resources whose lifecycle is identical to the bucket's,
# managed as part of this kind. Their region argument is required by the
# provider, so the spec requires an explicit region whenever any satellite
# is configured (CEL rule satellites_require_region).

# CORS through the standalone resource: the bucket's inline cors_rule is
# deprecated at the pinned provider and performs no drift detection.
resource "digitalocean_spaces_bucket_cors_configuration" "main" {
  count = length(var.spec.cors_rules) > 0 ? 1 : 0

  bucket = digitalocean_spaces_bucket.main.id
  region = var.spec.region

  dynamic "cors_rule" {
    for_each = var.spec.cors_rules
    content {
      allowed_methods = cors_rule.value.allowed_methods
      allowed_origins = cors_rule.value.allowed_origins
      allowed_headers = length(cors_rule.value.allowed_headers) > 0 ? cors_rule.value.allowed_headers : null
      expose_headers  = length(cors_rule.value.expose_headers) > 0 ? cors_rule.value.expose_headers : null
      id              = cors_rule.value.id != "" ? cors_rule.value.id : null
      max_age_seconds = cors_rule.value.max_age_seconds > 0 ? cors_rule.value.max_age_seconds : null
    }
  }
}

resource "digitalocean_spaces_bucket_policy" "main" {
  count = var.spec.policy != "" ? 1 : 0

  bucket = digitalocean_spaces_bucket.main.id
  region = var.spec.region
  policy = var.spec.policy
}

resource "digitalocean_spaces_bucket_logging" "main" {
  count = var.spec.logging != null ? 1 : 0

  bucket = digitalocean_spaces_bucket.main.id
  region = var.spec.region

  # The FK arrives flattened: the log-receiving bucket's name.
  target_bucket = var.spec.logging.target_bucket
  target_prefix = var.spec.logging.target_prefix
}
