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
  description = "Specification for the GCP Compute Engine regional address"
  type = object({
    # The GCP project that owns the reservation. The CLI's tfvars converter
    # resolves StringValueOrRef fields to their literal string before the
    # module runs, so this arrives as a plain string.
    # If empty, the provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # Name of the address in GCP (RFC1035). Immutable (ForceNew).
    address_name = string

    # Region for the reservation (e.g. us-central1). Immutable (ForceNew).
    region = string

    # Specific IP (EXTERNAL) or range start (INTERNAL); empty lets GCP
    # assign one. Immutable.
    address = optional(string, "")

    # EXTERNAL (public IP, the default via middleware) or INTERNAL
    # (private IP within a VPC/subnetwork). Immutable.
    address_type = optional(string, "EXTERNAL")

    description = optional(string, "")

    # IPV4 (middleware default) or IPV6. Immutable.
    ip_version = optional(string, "IPV4")

    # VPC network for INTERNAL VPC_PEERING / IPSEC_INTERCONNECT; arrives as
    # a plain self-link string. Immutable.
    network = optional(string, "")

    # Subnetwork for INTERNAL GCE_ENDPOINT / DNS_RESOLVER; arrives as a
    # plain self-link string. Immutable.
    subnetwork = optional(string, "")

    # PREMIUM or STANDARD — EXTERNAL addresses only. Immutable.
    network_tier = optional(string, "")

    # CIDR prefix length for peering/interconnect ranges (8-29). Immutable.
    prefix_length = optional(number, null)

    # GCE_ENDPOINT, SHARED_LOADBALANCER_VIP, VPC_PEERING,
    # IPSEC_INTERCONNECT, or DNS_RESOLVER (INTERNAL only). Immutable.
    purpose = optional(string, "")

    # VM or NETLB — external IPv6 endpoint type. Immutable.
    ipv6_endpoint_type = optional(string, "")

    # User labels merged with the platform labels (platform wins on key
    # conflicts). The one mutable surface on this resource.
    labels = optional(map(string), {})

    # BYOIP source: a PublicDelegatedPrefix URL for EXTERNAL addresses.
    # Immutable.
    ip_collection = optional(string, "")

    # DELETE (default), PREVENT, or ABANDON — what destroy does to the
    # reservation.
    deletion_policy = optional(string, "")
  })
}
