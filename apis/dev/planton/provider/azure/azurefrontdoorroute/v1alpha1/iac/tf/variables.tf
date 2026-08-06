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
  description = "Azure Front Door route specification"
  type = object({
    # The Front Door endpoint the route attaches to, by ARM ID.
    # References are resolved to literal IDs by the platform before the
    # module runs. ForceNew.
    endpoint_id = string

    # The route's name -- unique within the endpoint. ForceNew.
    route_name = string

    # The origin group that answers matched requests, by ARM ID.
    # Updatable in place.
    origin_group_id = string

    # Origin ARM IDs -- never sent to Azure; they exist purely to
    # sequence provisioning after the group's origins exist.
    origin_ids = optional(list(string))

    # Rule set ARM IDs whose delivery policies apply on this route.
    rule_set_ids = optional(list(string))

    # URL path patterns this route matches (each starts with "/").
    patterns_to_match = list(string)

    # Client-facing protocols as spec enum value names (HTTP / HTTPS).
    supported_protocols = list(string)

    # Origin-leg protocol as the spec enum's value name (MATCH_REQUEST /
    # HTTP_ONLY / HTTPS_ONLY). Absent means MATCH_REQUEST (tfvars drops
    # zero-valued proto fields).
    forwarding_protocol = optional(string)

    # Edge HTTP->HTTPS redirect. Azure defaults to true; requires both
    # protocols (enforced by the spec CEL).
    https_redirect_enabled = optional(bool)

    # Custom domain ARM IDs this route serves. The route side owns the
    # domain attachment.
    custom_domain_ids = optional(list(string))

    # Serve on the endpoint's generated *.azurefd.net hostname. Azure
    # defaults to true; false requires at least one custom domain
    # (enforced by the spec CEL).
    link_to_default_domain = optional(bool)

    # Whether the route matches traffic. Azure defaults to true.
    enabled = optional(bool)

    # Path prepended on the origin side before forwarding.
    origin_path = optional(string)

    # Edge caching. Absent means caching disabled (a real behavior
    # switch -- the provider transmits an explicit null).
    cache = optional(object({
      # Cache-key behavior as the spec enum's value name. Absent means
      # IGNORE_QUERY_STRING.
      query_string_caching_behavior = optional(string)
      # Parameter names for the *_SPECIFIED behaviors.
      query_strings = optional(list(string))
      # Edge compression for eligible content types.
      compression_enabled       = optional(bool)
      content_types_to_compress = optional(list(string))
    }))
  })
}
