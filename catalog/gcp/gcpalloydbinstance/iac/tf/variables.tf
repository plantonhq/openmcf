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
  description = "Specification for the GCP AlloyDB instance"
  type = object({
    project_id    = optional(string, "")
    cluster       = string
    instance_id   = string
    instance_type = optional(string, "READ_POOL")
    cpu_count     = optional(number, 0)
    machine_type  = optional(string, "")

    read_pool_config = optional(object({
      node_count = number
    }), null)

    availability_type = optional(string, "")
    database_flags    = optional(map(string), {})
    display_name      = optional(string, "")

    query_insights_config = optional(object({
      query_plans_per_minute  = optional(number, 0)
      query_string_length     = optional(number, 0)
      record_application_tags = optional(bool, false)
      record_client_address   = optional(bool, false)
    }), null)

    require_connectors = optional(bool, false)
    ssl_mode           = optional(string, "")

    activation_policy = optional(string, "")

    enable_public_ip          = optional(bool, false)
    enable_outbound_public_ip = optional(bool, false)
    authorized_external_networks = optional(list(object({
      cidr_range = string
    })), [])

    psc_instance_config = optional(object({
      allowed_consumer_projects = optional(list(string), [])
      psc_auto_connections = optional(list(object({
        consumer_network = optional(string, "")
        consumer_project = optional(string, "")
      })), [])
      psc_interface_configs = optional(list(object({
        network_attachment_resource = string
      })), [])
    }), null)

    # User labels; merged beneath the platform attribution labels.
    labels = optional(map(string), {})

    # Unstructured client metadata (annotations, not labels).
    annotations = optional(map(string), {})

    # Pin a ZONAL instance to a specific zone (ZONAL instances only).
    gce_zone = optional(string, "")

    # AlloyDB managed connection pooling (built-in pooler).
    connection_pool_config = optional(object({
      enabled = optional(bool, false)
      flags   = optional(map(string), {})
    }), null)

    # Draw private IPs from a specific PSA allocated range instead of the
    # parent cluster's. Immutable.
    allocated_ip_range_override = optional(string, "")

    # What happens to the instance in GCP on destroy:
    # DELETE (provider default), PREVENT, or ABANDON.
    deletion_policy = optional(string, "")
  })
}
