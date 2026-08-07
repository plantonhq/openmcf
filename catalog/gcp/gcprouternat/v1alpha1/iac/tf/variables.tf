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
  description = "Specification for the GCP Cloud Router + NAT gateway"
  type = object({
    # The GCP project that owns the router and NAT. The CLI's tfvars
    # converter resolves StringValueOrRef fields to their literal string
    # before the module runs, so this arrives as a plain string.
    # If empty, the provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # Name of the Cloud Router. Immutable.
    router_name = string

    # Name of the NAT configuration on the router. Immutable.
    nat_name = string

    # Region hosting the router and NAT (e.g. us-central1). Immutable.
    region = string

    # VPC self link; arrives as a plain string after ref resolution.
    # Immutable.
    vpc_self_link = string

    # BGP ASN for the router (private ASN). 0 means GCP assigns one.
    router_asn = optional(number, 0)

    # BGP keepalive interval in seconds (20-60). 0 means the API default.
    router_keepalive_interval = optional(number, 0)

    # PUBLIC (internet egress, default) or PRIVATE (NCC spoke-to-spoke NAT).
    # Immutable.
    type = optional(string, "")

    # ALL_SUBNETWORKS_ALL_IP_RANGES, ALL_SUBNETWORKS_ALL_PRIMARY_IP_RANGES,
    # or LIST_OF_SUBNETWORKS. Empty derives: LIST_OF_SUBNETWORKS when
    # subnetworks are listed, else ALL_SUBNETWORKS_ALL_IP_RANGES.
    source_subnetwork_ip_ranges_to_nat = optional(string, "")

    # Per-subnetwork NAT scoping (LIST_OF_SUBNETWORKS mode). Subnetwork
    # references arrive as plain self-link strings.
    subnetworks = optional(list(object({
      subnetwork               = string
      source_ip_ranges_to_nat  = optional(list(string), [])
      secondary_ip_range_names = optional(list(string), [])
    })), [])

    # Self links of GcpAddress reservations for manual NAT IP allocation.
    # Non-empty selects MANUAL_ONLY; empty selects AUTO_ONLY.
    nat_ips = optional(list(string), [])

    # Self links of addresses being drained (must already be in nat_ips).
    drain_nat_ips = optional(list(string), [])

    # PREMIUM or STANDARD tier for auto-allocated NAT IPs.
    auto_network_tier = optional(string, "")

    # Port allocation tuning. 0 means the API default.
    min_ports_per_vm               = optional(number, 0)
    max_ports_per_vm               = optional(number, 0)
    enable_dynamic_port_allocation = optional(bool, false)

    # Endpoint-independent mapping (RFC 5128); mutually exclusive with
    # dynamic port allocation.
    enable_endpoint_independent_mapping = optional(bool, false)

    # ENDPOINT_TYPE_VM (default), ENDPOINT_TYPE_SWG, or
    # ENDPOINT_TYPE_MANAGED_PROXY_LB. At most one. Immutable.
    endpoint_types = optional(list(string), [])

    # Connection idle/wait timeouts in seconds. 0 means the API default.
    udp_idle_timeout_sec             = optional(number, 0)
    icmp_idle_timeout_sec            = optional(number, 0)
    tcp_established_idle_timeout_sec = optional(number, 0)
    tcp_transitory_idle_timeout_sec  = optional(number, 0)
    tcp_time_wait_timeout_sec        = optional(number, 0)

    # NAT rules: dedicated NAT IPs/ranges for matching egress traffic.
    # Address and subnetwork references arrive as plain self-link strings.
    rules = optional(list(object({
      rule_number = number
      match       = string
      description = optional(string, "")
      action = optional(object({
        source_nat_active_ips    = optional(list(string), [])
        source_nat_drain_ips     = optional(list(string), [])
        source_nat_active_ranges = optional(list(string), [])
        source_nat_drain_ranges  = optional(list(string), [])
      }), null)
    })), [])

    # NAT translation logging: DISABLED, ERRORS_ONLY, ALL, or
    # TRANSLATIONS_ONLY. The tfvars converter emits the proto enum NAME as
    # a string; empty means the ERRORS_ONLY default.
    log_filter = optional(string, "ERRORS_ONLY")
  })
}
