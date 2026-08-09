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
  description = "Specification for the GCP VPC"
  type = object({
    # The GCP project that owns the network. The CLI's tfvars converter
    # resolves StringValueOrRef fields to their literal string before the
    # module runs, so this arrives as a plain string.
    # If empty, the provider's default project is used.
    project_id = optional(string, "")

    # Auto (true) vs custom (false) subnet mode. Custom is the default and
    # the production recommendation.
    auto_create_subnetworks = optional(bool, false)

    # Dynamic routing mode: REGIONAL (default) or GLOBAL. The converter
    # emits the proto enum's NAME as a string.
    routing_mode = optional(string, "REGIONAL")

    # Name of the VPC network (RFC1035). Immutable.
    network_name = string

    # Human-readable description. Immutable on this resource.
    description = optional(string, "")

    # MTU in bytes (1300-8896); null falls through to the API default (1460).
    mtu = optional(number)

    # ULA internal IPv6 enablement and optional explicit /48 from fd20::/20.
    enable_ula_internal_ipv6 = optional(bool, false)
    internal_ipv6_range      = optional(string, "")

    # Firewall policy vs classic rule evaluation order; empty falls through
    # to the API default (AFTER_CLASSIC_FIREWALL).
    network_firewall_policy_enforcement_order = optional(string, "")

    # Full or partial network profile URL. Immutable.
    network_profile = optional(string, "")

    # BGP best-path selection block; all fields mutable. No object default:
    # absence must stay distinguishable from an all-default block, because
    # always_compare_med is only sent when the block is present.
    bgp_best_path_selection = optional(object({
      mode               = optional(string, "")
      always_compare_med = optional(bool, false)
      inter_region_cost  = optional(string, "")
    }))

    # Suppress automatic 0.0.0.0/0 routes at creation. Immutable.
    delete_default_routes_on_create = optional(bool, false)

    # Create-time Resource Manager tag bindings (tagKeys/{id} => tagValues/{id}).
    resource_manager_tags = optional(map(string), {})

    # DELETE (default) / PREVENT / ABANDON; empty falls through to the
    # provider default (DELETE).
    deletion_policy = optional(string, "")
  })

  validation {
    condition     = can(regex("^[a-z]([a-z0-9-]{0,61}[a-z0-9])?$", var.spec.network_name))
    error_message = "Network name must be 1-63 characters, lowercase letters, numbers, or hyphens, starting with a letter and ending with a letter or number."
  }

  validation {
    condition     = contains(["REGIONAL", "GLOBAL"], var.spec.routing_mode)
    error_message = "routing_mode must be REGIONAL or GLOBAL."
  }
}
