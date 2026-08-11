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
  description = "Specification for the GCP DNS managed zone"
  type = object({
    # StringValueOrRef fields arrive as PLAIN STRINGS: the tfvars converter
    # flattens refs before the module ever sees them.
    project_id = optional(string, "")

    dns_name    = optional(string, "")
    description = optional(string, "")
    visibility  = optional(string, "public")

    private_visibility_config = optional(object({
      networks = optional(list(object({
        network_url = string
      })), [])
      gke_clusters = optional(list(object({
        gke_cluster_name = string
      })), [])
    }), null)

    dnssec_config = optional(object({
      state         = optional(string, "off")
      non_existence = optional(string, "")
      default_key_specs = optional(list(object({
        algorithm  = optional(string, "")
        key_length = optional(number, 0)
        key_type   = optional(string, "")
      })), [])
    }), null)

    forwarding_config = optional(object({
      target_name_servers = optional(list(object({
        ipv4_address    = optional(string, "")
        domain_name     = optional(string, "")
        forwarding_path = optional(string, "")
        # One address family per target — the spec CEL rejects a target
        # carrying both IPv4 and IPv6.
        ipv6_address = optional(string, "")
      })), [])
    }), null)

    peering_config = optional(object({
      target_network = string
    }), null)

    cloud_logging_config = optional(object({
      enable_logging = bool
    }), null)

    force_destroy = optional(bool, false)
    labels        = optional(map(string), {})

    # DELETE (default), PREVENT, or ABANDON — what destroy does to the zone
    # shell (force_destroy governs the records inside it).
    deletion_policy = optional(string, "")
  })
}
