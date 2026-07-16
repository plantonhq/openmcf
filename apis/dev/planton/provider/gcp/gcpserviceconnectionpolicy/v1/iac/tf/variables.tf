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
  description = "Specification for the GCP service connection policy"
  type = object({
    # The GCP project owning the network and policy. The CLI's tfvars
    # converter resolves StringValueOrRef fields to their literal string
    # before the module runs, so this arrives as a plain string.
    # If empty, the provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # Name for the policy resource. Empty falls back to metadata.name.
    # Immutable (ForceNew).
    policy_name = optional(string, "")

    # GCP region the policy applies to. Immutable (ForceNew).
    location = string

    # The consumer VPC network — a resource path or self-link URL (a
    # GcpVpcNetwork reference resolves to the resource path). Normalized
    # to projects/{project}/global/networks/{name} in locals.tf.
    # Immutable (ForceNew).
    network = string

    # The producer's published service class (e.g. gcp-memorystore).
    # Immutable (ForceNew).
    service_class = string

    # Free-text description of the policy's purpose.
    description = optional(string, "")

    # User labels merged beneath Planton platform labels (platform keys
    # win on conflict).
    labels = optional(map(string), {})

    # Private Service Connect configuration: subnets for endpoint IPs,
    # optional connection limit, optional producer-location allowlist.
    psc_config = optional(object({
      # Subnet resource paths or self-links (GcpSubnetwork references
      # resolve to self-links). Normalized to relative resource paths in
      # locals.tf.
      subnetworks = list(string)

      # Max PSC connections under this policy; 0 leaves GCP's default.
      limit = optional(number, 0)

      # PRODUCER_INSTANCE_LOCATION_UNSPECIFIED or
      # CUSTOM_RESOURCE_HIERARCHY_LEVELS.
      producer_instance_location = optional(string, "")

      # projects/{id}, folders/{num}, organizations/{num} entries; only
      # consulted under CUSTOM_RESOURCE_HIERARCHY_LEVELS.
      allowed_google_producers_resource_hierarchy_levels = optional(list(string), [])
    }), null)
  })

  validation {
    condition     = var.spec.location != ""
    error_message = "location is required."
  }

  validation {
    condition     = var.spec.network != ""
    error_message = "network is required."
  }

  validation {
    condition     = var.spec.service_class != ""
    error_message = "service_class is required."
  }

  validation {
    condition     = try(length(var.spec.psc_config.subnetworks) > 0, var.spec.psc_config == null)
    error_message = "psc_config.subnetworks must contain at least one subnet when psc_config is set."
  }
}
