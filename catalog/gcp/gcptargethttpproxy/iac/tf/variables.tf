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
  description = "Specification for the GCP Compute Engine global target HTTP proxy"
  type = object({
    # The GCP project that owns the proxy. The CLI's tfvars converter
    # resolves StringValueOrRef fields to their literal string before the
    # module runs, so this arrives as a plain string.
    # If empty, the provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # Name of the proxy in GCP (RFC1035). Empty defaults to metadata.name
    # (see locals.tf). Immutable (ForceNew).
    proxy_name = optional(string, "")

    description = optional(string, "")

    # The URL map the proxy routes through — required; arrives as a plain
    # string (self-link or name). The only mutable field on the resource.
    url_map = string

    # Idle client keep-alive in seconds (5-1200); 0 lets GCP apply its
    # default. Only honored by EXTERNAL_MANAGED load balancers. Immutable.
    http_keep_alive_timeout_sec = optional(number, 0)

    # Traffic Director binding; false for internet-facing frontends.
    # Immutable.
    proxy_bind = optional(bool, false)

    # DELETE (default) / PREVENT / ABANDON; empty falls through to the
    # provider default (DELETE).
    deletion_policy = optional(string, "")
  })
}
