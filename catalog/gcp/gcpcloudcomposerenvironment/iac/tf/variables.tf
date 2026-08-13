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
  description = "Specification for the GCP Cloud Composer environment"
  type = object({
    # The GCP project for the environment. The CLI's tfvars converter
    # resolves StringValueOrRef fields to their literal string before the
    # module runs, so this arrives as a plain string. If empty, the
    # provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # Region for the environment (e.g. us-central1). Immutable (ForceNew).
    region = string

    # Environment name. Empty falls back to metadata.name (see locals.tf).
    # Immutable (ForceNew).
    environment_name = optional(string, "")

    # Networking and node settings for the environment's GKE nodes.
    # All node networking is immutable (ForceNew) on Composer 2.
    node_config = optional(object({
      # network XOR composer_network_attachment (resolved refs).
      network         = optional(string, "")
      subnetwork      = optional(string, "")
      service_account = optional(string, "")
      tags            = optional(list(string), [])

      # Composer 3 PSC entry point.
      composer_network_attachment       = optional(string, "")
      composer_internal_ipv4_cidr_block = optional(string, "")

      # SNAT pod traffic to node IPs (ip-masq-agent DaemonSet).
      enable_ip_masq_agent = optional(bool, false)

      # VPC-native range assignment: named secondary range XOR CIDR,
      # per range.
      ip_allocation_policy = optional(object({
        cluster_secondary_range_name  = optional(string, "")
        cluster_ipv4_cidr_block       = optional(string, "")
        services_secondary_range_name = optional(string, "")
        services_ipv4_cidr_block      = optional(string, "")
      }), null)
    }), null)

    # Airflow software configuration. Packages, overrides, and env vars
    # update in place.
    software_config = optional(object({
      image_version            = optional(string, "")
      airflow_config_overrides = optional(map(string), {})
      pypi_packages            = optional(map(string), {})
      env_variables            = optional(map(string), {})
      web_server_plugins_mode  = optional(string, "")
      cloud_data_lineage_integration = optional(object({
        enabled = bool
      }), null)
    }), null)

    # Composer 2.x private networking. Immutable (ForceNew).
    private_environment_config = optional(object({
      enable_private_endpoint                = optional(bool, false)
      connection_type                        = optional(string, "")
      master_ipv4_cidr_block                 = optional(string, "")
      cloud_sql_ipv4_cidr_block              = optional(string, "")
      cloud_composer_network_ipv4_cidr_block = optional(string, "")
      cloud_composer_connection_subnetwork   = optional(string, "")
      enable_privately_used_public_ips       = optional(bool, false)
    }), null)

    # Per-component workload sizing. Updates in place.
    workloads_config = optional(object({
      scheduler = optional(object({
        cpu        = optional(number, 0)
        memory_gb  = optional(number, 0)
        storage_gb = optional(number, 0)
        count      = optional(number, 0)
      }), null)
      web_server = optional(object({
        cpu        = optional(number, 0)
        memory_gb  = optional(number, 0)
        storage_gb = optional(number, 0)
      }), null)
      worker = optional(object({
        cpu        = optional(number, 0)
        memory_gb  = optional(number, 0)
        storage_gb = optional(number, 0)
        min_count  = optional(number, 0)
        max_count  = optional(number, 0)
      }), null)
      triggerer = optional(object({
        cpu       = optional(number, 0)
        memory_gb = optional(number, 0)
        count     = optional(number, 0)
      }), null)
      dag_processor = optional(object({
        cpu        = optional(number, 0)
        memory_gb  = optional(number, 0)
        storage_gb = optional(number, 0)
        count      = optional(number, 0)
      }), null)
    }), null)

    # Managed infrastructure capacity. Updates in place.
    environment_size = optional(string, "")

    # STANDARD_RESILIENCE or HIGH_RESILIENCE. Updates in place.
    resilience_mode = optional(string, "")

    # Resolved CMEK key resource ID. Empty means Google-managed keys.
    # Immutable (ForceNew).
    kms_key_name = optional(string, "")

    # Maintenance window (RFC3339 times + RRULE recurrence). Updates in
    # place.
    maintenance_window = optional(object({
      start_time = string
      end_time   = string
      recurrence = string
    }), null)

    # Scheduled snapshots for disaster recovery.
    recovery_config = optional(object({
      enabled                    = bool
      snapshot_location          = optional(string, "")
      snapshot_creation_schedule = optional(string, "")
      time_zone                  = optional(string, "")
    }), null)

    # IP allowlist for the Airflow web UI. Updates in place.
    web_server_network_access_control = optional(object({
      allowed_ip_ranges = optional(list(object({
        value       = string
        description = optional(string, "")
      })), [])
    }), null)

    # IP-based access control for the environment's GKE cluster master.
    master_authorized_networks_config = optional(object({
      enabled = bool
      cidr_blocks = optional(list(object({
        cidr_block   = string
        display_name = optional(string, "")
      })), [])
    }), null)

    # Retention for task logs and the Airflow metadata database.
    data_retention_config = optional(object({
      task_logs_storage_mode          = optional(string, "")
      airflow_metadata_retention_mode = optional(string, "")
      airflow_metadata_retention_days = optional(number, 0)
    }), null)

    # Existing bucket for DAGs/plugins/data instead of the auto-created
    # one (resolved to the bucket name). Immutable (ForceNew).
    storage_bucket = optional(string, "")

    # Composer 3 private environment flags.
    enable_private_environment = optional(bool, false)
    enable_private_builds_only = optional(bool, false)

    # User labels merged beneath Planton platform labels (platform keys
    # win on conflict).
    labels = optional(map(string), {})

    # Client-side destroy behavior: DELETE (default), PREVENT, ABANDON.
    deletion_policy = optional(string, "")
  })

  validation {
    condition     = var.spec.region != ""
    error_message = "region is required."
  }
}
