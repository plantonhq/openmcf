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
  description = "Specification for the GCP Compute Engine global address"
  type = object({
    # The GCP project that owns the reservation. The CLI's tfvars converter
    # resolves StringValueOrRef fields to their literal string before the
    # module runs, so this arrives as a plain string.
    # If empty, the provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # Name of the address in GCP (RFC1035). Immutable (ForceNew).
    address_name = string

    # Specific IP (EXTERNAL) or range start (INTERNAL); empty lets GCP
    # assign one. Immutable.
    address = optional(string, "")

    # EXTERNAL (public IP, the default via middleware) or INTERNAL
    # (private range for VPC peering / PSC). Immutable.
    address_type = optional(string, "EXTERNAL")

    description = optional(string, "")

    # IPV4 (middleware default) or IPV6. Immutable.
    ip_version = optional(string, "IPV4")

    # VPC network for INTERNAL reservations; arrives as a plain self-link
    # string. Immutable.
    network = optional(string, "")

    # CIDR prefix length for VPC_PEERING ranges (8-29). Immutable.
    prefix_length = optional(number, null)

    # VPC_PEERING or PRIVATE_SERVICE_CONNECT (INTERNAL only). Immutable.
    purpose = optional(string, "")
  })
}
