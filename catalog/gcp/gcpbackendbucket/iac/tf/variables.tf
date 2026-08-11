variable "metadata" {
  description = "Metadata for the resource, including name and labels"
  type = object({
    name    = string,
    id      = optional(string),
    org     = optional(string),
    env     = optional(string),
    labels  = optional(map(string)),
    tags    = optional(list(string)),
    version = optional(object({ id = string, message = string }))
  })
}

variable "spec" {
  description = "Specification for the GCP Compute Engine backend bucket"
  type = object({
    # The GCP project that owns the backend bucket. The CLI's tfvars
    # converter resolves StringValueOrRef fields to their literal string
    # before the module runs, so this arrives as a plain string.
    # If empty, the provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # Name of the backend bucket in GCP (RFC1035). Empty defaults to
    # metadata.name (see locals.tf). Immutable (ForceNew).
    backend_bucket_name = optional(string, "")

    # The Cloud Storage bucket serving as the origin (resolved from a
    # GcpGcsBucket reference or given directly). Mutable — origin swaps are
    # in-place updates.
    bucket_name = string

    description = optional(string)

    # Cache at Google's edge with Cloud CDN. cdn_policy only takes effect
    # while this is true.
    enable_cdn = optional(bool, false)

    # Cloud CDN caching behavior. TTL fields left at 0 are treated as unset
    # so the GCP API applies its own defaults (see locals.tf).
    cdn_policy = optional(object({
      cache_mode                   = optional(string)
      client_ttl                   = optional(number)
      default_ttl                  = optional(number)
      max_ttl                      = optional(number)
      negative_caching             = optional(bool)
      negative_caching_policy      = optional(list(object({ code = number, ttl = optional(number) })), [])
      serve_while_stale            = optional(number)
      request_coalescing           = optional(bool)
      signed_url_cache_max_age_sec = optional(number)
      cache_key_policy = optional(object({
        query_string_whitelist = optional(list(string), [])
        include_http_headers   = optional(list(string), [])
      }))
      bypass_cache_on_request_headers = optional(list(object({ header_name = string })), [])
    }))

    # Load-balancer response compression: AUTOMATIC or DISABLED (empty keeps
    # the GCP default of no compression).
    compression_mode = optional(string, "")

    # Response headers the load balancer adds, "Header-Name: value" form.
    custom_response_headers = optional(list(string), [])

    # Self-link of a Cloud Armor EDGE security policy (resolved from a
    # GcpCloudArmorPolicy reference or given directly).
    edge_security_policy = optional(string, "")

    # INTERNAL_MANAGED for cross-region internal ALBs; empty for external
    # load balancers (the common case). Immutable (ForceNew).
    load_balancing_scheme = optional(string, "")

    # Cloud CDN signed-URL keys (at most 3). Each key_value is secret
    # material — it never appears in outputs.
    signed_url_keys = optional(list(object({
      name      = string
      key_value = string
    })), [])

    # Resource Manager tags bound at create time (tagKeys/{id} ->
    # tagValues/{id}). Immutable.
    resource_manager_tags = optional(map(string), {})

    # DELETE (default), PREVENT, or ABANDON — one switch governing destroy
    # for the backend bucket AND its signed-URL keys.
    deletion_policy = optional(string, "")
  })

  # NOTE: never guard optional strings with coalesce() here — HCL's coalesce
  # skips empty strings as well as nulls, so coalesce("", "") errors and the
  # validation fails on a legitimately-empty value.
  validation {
    condition     = try(var.spec.backend_bucket_name, "") == "" || can(regex("^[a-z]([-a-z0-9]{0,61}[a-z0-9])?$", var.spec.backend_bucket_name))
    error_message = "backend_bucket_name must be RFC1035-compliant: 1-63 lowercase letters, digits, or hyphens."
  }

  validation {
    condition     = length(var.spec.bucket_name) > 0
    error_message = "bucket_name is required — the GCS bucket whose objects are served."
  }

  validation {
    condition     = contains(["", "AUTOMATIC", "DISABLED"], var.spec.compression_mode)
    error_message = "compression_mode must be AUTOMATIC or DISABLED."
  }

  validation {
    condition     = contains(["", "INTERNAL_MANAGED"], var.spec.load_balancing_scheme)
    error_message = "load_balancing_scheme must be INTERNAL_MANAGED, or left unset for external load balancers."
  }

  # HCL's || does not short-circuit, so the nullable bool is guarded with
  # coalesce — Cloud CDN only fronts external load balancers.
  validation {
    condition     = !(coalesce(var.spec.enable_cdn, false) && var.spec.load_balancing_scheme == "INTERNAL_MANAGED")
    error_message = "Cloud CDN cannot be enabled on an INTERNAL_MANAGED backend bucket."
  }

  validation {
    condition     = length(var.spec.signed_url_keys) <= 3
    error_message = "at most 3 signed-URL keys are supported per backend bucket."
  }
}
