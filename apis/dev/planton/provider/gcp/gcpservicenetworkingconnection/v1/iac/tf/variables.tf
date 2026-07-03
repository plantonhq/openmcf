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
  description = "Specification for the GCP private services access connection"
  type = object({
    # The GCP project used for in-module API enablement. The CLI's tfvars
    # converter resolves StringValueOrRef fields to their literal string
    # before the module runs, so this arrives as a plain string.
    # If empty, the provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # The VPC network to peer with the service producer — a network name or
    # full self-link URL (a GcpVpc reference resolves to the self-link).
    # Immutable (ForceNew).
    network = string

    # The service producer to peer with. Empty falls through to
    # servicenetworking.googleapis.com — the producer behind Cloud SQL,
    # AlloyDB, Memorystore, and Filestore private IP. Immutable.
    service = optional(string, "")

    # Names of INTERNAL VPC_PEERING global address ranges reserved for the
    # producer (GcpGlobalAddress references resolve to their names). At least
    # one is required. Mutable: appending ranges grows producer capacity
    # without disturbing already-provisioned service subnets.
    reserved_peering_ranges = list(string)

    # Adopt a pre-existing connection for the same (network, service) pair
    # by converting the create-time "Cannot modify allocated ranges" failure
    # into an in-place update of its reserved ranges.
    update_on_creation_fail = optional(bool, false)
  })

  validation {
    condition     = length(var.spec.reserved_peering_ranges) > 0
    error_message = "at least one reserved peering range is required."
  }

  validation {
    condition     = var.spec.network != ""
    error_message = "network is required."
  }
}
