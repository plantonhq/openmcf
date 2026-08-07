# Enable the Compute Engine API so a fresh project can host the backend
# bucket. disable_on_destroy is false: tearing down one backend bucket must
# never disable the API for everything else in the project.
resource "google_project_service" "compute_api" {
  project = local.project_id
  service = "compute.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A Compute Engine backend bucket — serves a Cloud Storage bucket's objects
# through an external HTTP(S) load balancer, optionally cached at Google's
# edge by Cloud CDN. URL maps route static paths here while dynamic paths go
# to backend services.
#
# name and project are immutable (ForceNew): changing either destroys and
# recreates the backend bucket, briefly breaking every URL map referencing
# the old self_link. bucket_name is deliberately mutable — origin swaps
# (blue/green static releases) are in-place updates.
#
# CDN is a policy ON this resource, not a separate GCP object: enable_cdn
# turns edge caching on and cdn_policy tunes it. TTL/cache-mode coherence
# (e.g. USE_ORIGIN_HEADERS forbids explicit TTLs) is enforced by the spec
# before deploy.
resource "google_compute_backend_bucket" "this" {
  name        = local.backend_bucket_name
  project     = local.project_id
  bucket_name = var.spec.bucket_name
  description = var.spec.description

  enable_cdn            = var.spec.enable_cdn
  compression_mode      = local.compression_mode
  load_balancing_scheme = local.load_balancing_scheme

  # Attaches by reference — the Cloud Armor EDGE policy itself is its own
  # composable node (GcpCloudArmorPolicy), never embedded here.
  edge_security_policy = local.edge_security_policy

  custom_response_headers = length(var.spec.custom_response_headers) > 0 ? var.spec.custom_response_headers : null

  dynamic "cdn_policy" {
    for_each = local.cdn_policy != null ? [local.cdn_policy] : []
    content {
      cache_mode                   = cdn_policy.value.cache_mode
      client_ttl                   = cdn_policy.value.client_ttl
      default_ttl                  = cdn_policy.value.default_ttl
      max_ttl                      = cdn_policy.value.max_ttl
      negative_caching             = cdn_policy.value.negative_caching
      serve_while_stale            = cdn_policy.value.serve_while_stale
      request_coalescing           = cdn_policy.value.request_coalescing
      signed_url_cache_max_age_sec = cdn_policy.value.signed_url_cache_max_age_sec

      dynamic "negative_caching_policy" {
        for_each = cdn_policy.value.negative_caching_policy
        content {
          code = negative_caching_policy.value.code
          # A 0 TTL is meaningful here ("cache but expire immediately" is not
          # useful; GCP treats 0 as don't-cache-this-code), so pass it as-is.
          ttl = negative_caching_policy.value.ttl
        }
      }

      dynamic "cache_key_policy" {
        for_each = cdn_policy.value.cache_key_policy != null ? [cdn_policy.value.cache_key_policy] : []
        content {
          query_string_whitelist = cache_key_policy.value.query_string_whitelist
          include_http_headers   = cache_key_policy.value.include_http_headers
        }
      }

      dynamic "bypass_cache_on_request_headers" {
        for_each = cdn_policy.value.bypass_cache_on_request_headers
        content {
          header_name = bypass_cache_on_request_headers.value.header_name
        }
      }
    }
  }

  depends_on = [google_project_service.compute_api]
}

# Cloud CDN signed-URL keys — folded into this kind rather than modeled as a
# separate node: keys are never referenced by other resources, GCP caps them
# at 3 per bucket, and their lifecycle is the bucket's. Each key is
# immutable in GCP (add/delete only); changing a key's value forces
# replacement of that key resource, which is exactly the rotation semantics
# signed URLs need (add new key -> re-sign -> remove old).
resource "google_compute_backend_bucket_signed_url_key" "this" {
  for_each = { for signed_url_key in var.spec.signed_url_keys : signed_url_key.name => signed_url_key }

  name           = each.value.name
  key_value      = each.value.key_value # secret material; never surfaced in outputs
  backend_bucket = google_compute_backend_bucket.this.name
  project        = local.project_id
}
