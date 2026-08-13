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
  description = "Specification for the GCP Memorystore for Redis instance"
  type = object({
    # The GCP project that owns the instance. The CLI's tfvars converter
    # resolves StringValueOrRef fields to their literal string before the
    # module runs, so this arrives as a plain string.
    # If empty, the provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # Name of the Redis instance in GCP. Immutable.
    instance_name = string

    # Region hosting the instance (e.g. us-central1). Immutable.
    region = string

    # BASIC (standalone, no SLA) or STANDARD_HA (primary + failover replica,
    # 99.9% SLA). Immutable.
    tier = string

    # Memory in GiB. Mutable (in-place resize); STANDARD_HA and read
    # replicas require at least 5.
    memory_size_gb = number

    # Engine version, e.g. REDIS_7_2. Upgrades apply in place; a downgrade
    # replaces the instance.
    redis_version = optional(string, "")

    # Human-readable display name.
    display_name = optional(string, "")

    # Primary zone within the region. Immutable.
    location_id = optional(string, "")

    # Replica zone (STANDARD_HA only); must differ from location_id.
    # Immutable.
    alternative_location_id = optional(string, "")

    # VPC self link; arrives as a plain string after ref resolution.
    # Empty means the project's default network. Immutable.
    authorized_network = optional(string, "")

    # DIRECT_PEERING (default) or PRIVATE_SERVICE_ACCESS (requires the VPC
    # to already carry a service networking connection). Immutable.
    connect_mode = optional(string, "")

    # DIRECT_PEERING: /29 CIDR (or empty for auto). PRIVATE_SERVICE_ACCESS:
    # the NAME of an allocated address range on the PSA connection.
    # Immutable.
    reserved_ip_range = optional(string, "")

    # Additional range for node placement — required when enabling read
    # replicas on an existing instance. /28 CIDR, range name, or "auto".
    # Mutable.
    secondary_ip_range = optional(string, "")

    # Redis AUTH: when true GCP generates and rotates the AUTH string
    # (exported as a sensitive output).
    auth_enabled = optional(bool, false)

    # DISABLED or SERVER_AUTHENTICATION (TLS; pair with the server_ca_certs
    # output). Immutable.
    transit_encryption_mode = optional(string, "")

    # Redis configuration parameters (e.g. maxmemory-policy).
    redis_configs = optional(map(string), {})

    # Weekly maintenance window start (UTC). Fixed 1-hour duration.
    maintenance_window = optional(object({
      day    = string
      hour   = optional(number, 0)
      minute = optional(number, 0)
      # Human-readable description of the policy (max 512 characters).
      description = optional(string, "")
    }), null)

    # Self-service maintenance version — set to a newer available version
    # to apply maintenance on your schedule instead of GCP's rollout.
    maintenance_version = optional(string, "")

    # READ_REPLICAS_DISABLED or READ_REPLICAS_ENABLED (STANDARD_HA only).
    # Set at creation time.
    read_replicas_mode = optional(string, "")

    # Read replica count (1-5) when read replicas are enabled.
    replica_count = optional(number, 0)

    # RDB snapshot persistence.
    persistence_config = optional(object({
      persistence_mode        = string
      rdb_snapshot_period     = optional(string, "")
      rdb_snapshot_start_time = optional(string, "")
    }), null)

    # CMEK key resource id; arrives as a plain string after ref resolution.
    # Immutable.
    customer_managed_key = optional(string, "")

    # User labels, merged beneath the platform attribution labels
    # (see locals.tf).
    labels = optional(map(string), {})

    # Destroy guard. The spec defaults this to true (Planton middleware
    # materializes the default), and the module sends it explicitly so
    # destroy behavior is identical on both engines.
    deletion_protection = optional(bool, true)

    # Deletion policy: "", "DELETE" (default), "PREVENT" (destroy fails),
    # or "ABANDON" (remove from management, leave running in GCP).
    deletion_policy = optional(string, "")
  })

  validation {
    condition     = contains(["", "DELETE", "PREVENT", "ABANDON"], var.spec.deletion_policy)
    error_message = "deletion_policy must be one of: DELETE, PREVENT, ABANDON."
  }
}
